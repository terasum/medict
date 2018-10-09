'use strict'

import { app, BrowserWindow, ipcMain } from 'electron'
import mt from '../common/msgType'

/**
 * Set `__static` path to static files in production
 * https://simulatedgreg.gitbooks.io/electron-vue/content/en/using-static-assets.html
 */
if (process.env.NODE_ENV !== 'development') {
  global.__static = require('path').join(__dirname, '/static').replace(/\\/g, '\\\\')
}

let mainWindow
const winURL = process.env.NODE_ENV === 'development'
  ? `http://localhost:9080`
  : `file://${__dirname}/index.html`

let bgWindow
const bgURL = process.env.NODE_ENV === 'development'
  ? `http://localhost:9080/background.html`
  : `file://${__dirname}/background.html`

/**
 *
 */
function createBackgroundWin () {
  bgWindow = new BrowserWindow({
    show: false,
    webPreferences: {
      webSecurity: false,
      nodeIntegrationInWorker: true
    }
  })

  bgWindow.loadURL(bgURL)
  // if bgwin Closed close main window too
  bgWindow.on('ready-to-show', () => {
    bgWindow.webContents.openDevTools({ mode: 'detach' })
  })

  bgWindow.on('closed', () => {
    unsetMainBgBridge()
    // close the main window
    if (mainWindow !== null) {
      mainWindow.close()
    }

    bgWindow = null
  })
}

function createWindow () {
  /**
   * Initial window options
   */
  mainWindow = new BrowserWindow({
    height: 473,
    useContentSize: true,
    show: false,
    width: 743,
    titleBarStyle: 'hidden',
    'auto-hide-menu-bar': true,
    // enable multi-worker
    webPreferences: {
      nodeIntegrationInWorker: true
    }
  })
  mainWindow.setMaximizable(false)
  mainWindow.setResizable(false)
  mainWindow.setMinimizable(true)

  // hide MenuBar
  mainWindow.setMenu(null)
  mainWindow.setAutoHideMenuBar(true)

  // console.log(winURL)
  // console.log(__static + '/src/renderer/worker/worker.js')

  mainWindow.loadURL(winURL)

  mainWindow.once('ready-to-show', () => {
    mainWindow.show()
    mainWindow.webContents.openDevTools({ mode: 'detach' })
  })

  mainWindow.on('closed', () => {
    unsetMainBgBridge()
    mainWindow = null
    // bgWindow = null
  })
}

app.on('ready', () => {
  // start communication
  setMainBgBridge()

  createWindow()
  // create only once
  createBackgroundWin()
})

app.on('activate-with-no-open-windows', () => {
  if (!mainWindow) {
    createWindow()
    setMainBgBridge()
  }
})

app.on('window-all-closed', () => {
  // on macos, the app will still live on docker
  // and will active again
  if (process.platform !== 'darwin') {
    app.quit()
  } else {
    // stop communication
    unsetMainBgBridge()
  }
})

app.on('activate', () => {
  if (mainWindow === null) {
    setMainBgBridge()
    createWindow()
  }
  if (bgWindow === null) {
    // TODO ?
    createBackgroundWin()
  }
})

/**
 * send message to main window process
 * @param {*} event: cause event
 * @param {*} payload: the message payload
 */
function toMainListener (event, payload) {
  mainWindow.webContents.send(mt.MsgToMain, payload)
}

/**
 * send messagt to background window process
 * @param {*} event: cause event
 * @param {*} payload: the passage event
 */
function toBgListener (event, payload) {
  bgWindow.webContents.send(mt.MsgToBackground, payload)
}

// restart bgmain
ipcMain.on('restartBG', () => {
  bgWindow = null
  setTimeout(() => {
    createBackgroundWin()
  }, 20)
})

/**
 * set Main window process and background window process communications
 */
function setMainBgBridge () {
  unsetMainBgBridge()
  ipcMain.on(mt.MsgToMain, toMainListener)
  ipcMain.on(mt.MsgToBackground, toBgListener)
}

/**
 * unset Main window process and background window process communications
 */
function unsetMainBgBridge () {
  ipcMain.removeListener(mt.MsgToMain, toMainListener)
  ipcMain.removeListener(mt.MsgToBackground, toBgListener)
}

/**
 * Auto Updater
 *
 * Uncomment the following code below and install `electron-updater` to
 * support auto updating. Code Signing with a valid certificate is required.
 * https://simulatedgreg.gitbooks.io/electron-vue/content/en/using-electron-builder.html#auto-updating
 */

/*
import { autoUpdater } from 'electron-updater'

autoUpdater.on('update-downloaded', () => {
  autoUpdater.quitAndInstall()
})

app.on('ready', () => {
  if (process.env.NODE_ENV === 'production') autoUpdater.checkForUpdates()
})
 */
