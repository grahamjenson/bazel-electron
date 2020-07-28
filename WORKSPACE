workspace(name = "com_github_grahamjenson_bazel_electron")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file")


http_archive(
    name = "io_bazel_rules_go",
    sha256 = "a8d6b1b354d371a646d2f7927319974e0f9e52f73a2452d2b3877118169eb6bb",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.23.3/rules_go-v0.23.3.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.23.3/rules_go-v0.23.3.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
go_rules_dependencies()
go_register_toolchains()


###
# Electron Binaries
###

# We do not use http_archive because Bazel does not deal very well with symlinks in the file
http_file(
    name = "electron_release",
    sha256 = "594326256e5dc6ddf6c5b1ecc35563416b9d98b824c73cba287aa92cca1f41ec",
    urls = ["https://github.com/electron/electron/releases/download/v8.4.1/electron-v8.4.1-darwin-x64.zip"],
)