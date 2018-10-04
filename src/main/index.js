'use strict'

import { app, BrowserWindow, ipcMain } from 'electron'
import MsgType from '../util/constant'

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
    show: true,
    height: 100,
    useContentSize: true,
    width: 200

  })

  bgWindow.loadURL(bgURL)
  // if bgwin Closed close main window too
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

function toMainListener (event, payload) {
  mainWindow.webContents.send(MsgType.MsgToMain, payload)
}

function toBgListener (event, payload) {
  bgWindow.webContents.send(MsgType.MsgToBackground, payload)
}

/**
 * mainWin <=> bgWin
 */
function setMainBgBridge () {
  unsetMainBgBridge()
  ipcMain.on(MsgType.MsgToMain, toMainListener)
  ipcMain.on(MsgType.MsgToBackground, toBgListener)
}

function unsetMainBgBridge () {
  ipcMain.removeListener(MsgType.MsgToMain, toMainListener)
  ipcMain.removeListener(MsgType.MsgToBackground, toBgListener)
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
