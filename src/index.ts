import { app, BrowserWindow, ipcMain } from 'electron';
import { createSubWindow, WindowOption } from './subwindow';
import { logger } from './utils/logger';

import './main/init';

// This allows TypeScript to pick up the magic constant that's auto-generated by Forge's Webpack
// plugin that tells the Electron app where to look for the Webpack-bundled app code (depending on
// whether you're running in development or production).
declare const MAIN_WINDOW_WEBPACK_ENTRY: string;

// Handle creating/removing shortcuts on Windows when installing/uninstalling.
if (require('electron-squirrel-startup')) {
  // eslint-disable-line global-require
  app.quit();
}

const createWindow = (): void => {
  logger.info('📃 userData path: %s', app.getPath('userData'));
  logger.info('📃 appData path: %s', app.getPath('appData'));
  logger.info('📃 temp path: %s', app.getPath('temp'));
  logger.info('📃 documents path: %s', app.getPath('documents'));
  logger.info('📃 logs path: %s', app.getPath('logs'));

  // Create the browser window.
  const mainWindow = new BrowserWindow({
    height: 600,
    width: 768,
    titleBarStyle: 'hidden',
    // The lines below solved the issue
    webPreferences: {
      nodeIntegration: true,
      contextIsolation: false,
      webviewTag: true,
    },
  });

  // and load the index.html of the app.
  mainWindow.loadURL(MAIN_WINDOW_WEBPACK_ENTRY);

  // sub-windows
  mainWindow.webContents.setWindowOpenHandler(({ url }) => {
    // if (url.startsWith('https://github.com/')) {
    //   return { action: 'allow' }
    // }
    return { action: 'allow' };
  });

  mainWindow.webContents.on('did-create-window', childWindow => {
    // For example...
    // childWindow.webContents('will-navigate', (e) => {
    //   e.preventDefault()
    // })
  });
  mainWindow.on('unresponsive', function() {
    logger.error('window crashed');
  });

  mainWindow.webContents.on('did-fail-load', function() {
    logger.error('window failed load');
  });
  // special ipcmain
  ipcMain.on('createSubWindow', function(event: any, args: WindowOption) {
    logger.info(event);
    createSubWindow(mainWindow, args);
  });
  if (
    process.env.NODE_ENV == 'development' ||
    process.env.NODE_ENV == 'production'
  ) {
    mainWindow.webContents.openDevTools();
  }
  // special ipcmain
  ipcMain.on('openDevTool', function(event: any, args: WindowOption) {
    mainWindow.webContents.openDevTools();
  });

  ipcMain.on('errorInWindow', function(event, data) {
    console.error(data);
  });
};

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.on('ready', createWindow);

// Quit when all windows are closed, except on macOS. There, it's common
// for applications and their menu bar to stay active until the user quits
// explicitly with Cmd + Q.
app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('activate', () => {
  // On OS X it's common to re-create a window in the app when the
  // dock icon is clicked and there are no other windows open.
  if (BrowserWindow.getAllWindows().length === 0) {
    createWindow();
  }
});

// In this file you can include the rest of your app's specific main process
// code. You can also put them in separate files and import them here.

// listen uncaught Exception
process.on('uncaughtException', function(error) {
  // Handle the error
  logger.error(error);
});

// for logger
import './utils/logger';
