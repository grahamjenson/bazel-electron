package main

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const PACKAGE_JSON string = `{ "main": "main.js" }`

func main() {
	outputFile := os.Args[1]
	name := os.Args[2]
	mainJS := os.Args[3]
	indexHTML := os.Args[4]
	electronZIP := os.Args[5]

	// 1. Extract electronZIP file
	// 2. add Contents/Contents/Resources/app/ main.js index.html package.json
	// 3. tar contents back into <name>.app

	// from https://golangcode.com/unzip-files-in-go/
	rawFiles, err := Unzip(electronZIP, "electronZIP")
	if err != nil {
		log.Fatal("\nError UNZIP: %s\n\n", err.Error())
	}

	appFiles := map[string]string{}
	for _, f := range rawFiles {
		if strings.HasPrefix(f, "electronZIP/Electron.app/") {
			appFiles[f] = name + ".app/" + strings.TrimPrefix(f, "electronZIP/Electron.app/")
		}
	}

	// Write out our application files
	appFiles[mainJS] = name + ".app/Contents/Resources/app/main.js"
	appFiles[indexHTML] = name + ".app/Contents/Resources/app/index.html"
	ioutil.WriteFile("package.json", []byte(PACKAGE_JSON), 0644)
	appFiles["package.json"] = name + ".app/Contents/Resources/app/package.json"

	err = tarFiles(outputFile, appFiles)
	if err != nil {
		log.Fatal("\nError ZIP: %s\n\n", err.Error())
	}
}

// tarrer walks paths to create tar file tarName
func tarFiles(tarName string, files map[string]string) error {
	tarFile, err := os.Create(tarName)
	if err != nil {
		return err
	}

	defer func() {
		err = tarFile.Close()
	}()

	tw := tar.NewWriter(tarFile)
	defer func() {
		err = tw.Close()
	}()

	// walk each specified path and add encountered file to tar
	for path, name := range files {
		// validate path
		path = filepath.Clean(path)

		// make sure to handle symlinks properly
		finfo, link, err := getFileInfo(path)
		if err != nil {
			return err
		}

		// add file to tar
		srcFile, err := os.Open(filepath.Clean(path))
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// fill in header info using func FileInfoHeader
		hdr, err := tar.FileInfoHeader(finfo, finfo.Name())
		if err != nil {
			return err
		}

		// ensure header has relative file path
		hdr.Name = name

		// strip headers that can change the SHA
		hdr.ModTime = time.Unix(0, 0)
		hdr.AccessTime = time.Unix(0, 0)
		hdr.ChangeTime = time.Unix(0, 0)
		hdr.Uid = 0
		hdr.Gid = 0
		hdr.Uname = ""
		hdr.Gname = ""
		hdr.Mode = 0755

		if finfo.Mode()&os.ModeSymlink != 0 {
			hdr.Linkname = link
		}

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		// Copy the contents of the file only if it is the regular type
		if hdr.Typeflag != tar.TypeReg {
			continue
		}

		_, err = io.Copy(tw, srcFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func getFileInfo(path string) (os.FileInfo, string, error) {
	finfo, err := os.Lstat(path)
	if err != nil {
		return finfo, "", err
	}

	// Return unless is it a Symlink
	if finfo.Mode()&os.ModeSymlink == 0 {
		return finfo, "", nil
	}

	// Get Link
	link, err := os.Readlink(path)
	if err != nil {
		return finfo, "", err
	}

	// If a non-relative symlink then return non-symlink FileInfo
	if strings.HasPrefix(link, "/") {
		finfo, err := os.Stat(path)
		return finfo, "", err
	}
	return finfo, link, err
}

// from https://golangcode.com/unzip-files-in-go/
func Unzip(src string, dest string) ([]string, error) {
	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		filenames = append(filenames, fpath)

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		defer rc.Close()
		if err != nil {
			return filenames, err
		}

		if f.Mode()&os.ModeSymlink != 0 {
			buff, err := ioutil.ReadAll(rc)
			if err != nil {
				return filenames, err
			}
			oldname := string(buff)
			err = os.Symlink(oldname, fpath)
			if err != nil {
				return filenames, err
			}
			continue
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		defer outFile.Close()
		_, err = io.Copy(outFile, rc)

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
