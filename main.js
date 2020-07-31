const {app, BrowserWindow} = require('electron')
let mainWindow = null

function initialize () {
  app.setName('Electron Simple App')

  app.on('ready', () => {
    createWindow()
  })

  app.on('window-all-closed', () => {
      app.quit()
  })

  app.on('activate', () => {
    if (mainWindow === null) {
      createWindow()
    }
  })
}

function createWindow () {
  const windowOptions = {
    width: 600,
    minWidth: 600,
    height: 300,
    title: app.getName()
  }

  mainWindow = new BrowserWindow(windowOptions)
  mainWindow.loadURL('file://' + __dirname + '/index.html')

  mainWindow.webContents.openDevTools()

  mainWindow.on('closed', () => {
    mainWindow = null
  })
}

initialize()
