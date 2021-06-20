import { BrowserWindow } from 'electron';

interface WindowOption {
    width: number,
    height: number,
    html: string,
    titleBarStyle?: ('default' | 'hidden' | 'hiddenInset' | 'customButtonsOnHover'),
    nodeIntegration: boolean,
    contextIsolation: boolean
}

const createSubWindow = (parent: BrowserWindow, options: WindowOption): void => {
    // Create the browser window.
    const mainWindow = new BrowserWindow({
      parent,
      height: options.height,
      width: options.width,
      titleBarStyle: options.titleBarStyle,
        // The lines below solved the issue
        webPreferences: {
          nodeIntegration: options.nodeIntegration,
          contextIsolation: options.contextIsolation,
      }
    });
  
    // and load the index.html of the app.
    mainWindow.loadURL(options.html);
  
    // Open the DevTools.
    mainWindow.webContents.openDevTools();
  
    // sub-windows
    mainWindow.webContents.setWindowOpenHandler(({ url }) => {
      // if (url.startsWith('https://github.com/')) {
      //   return { action: 'allow' }
      // }
      return { action: 'allow' }
    })
    
    mainWindow.webContents.on('did-create-window', (childWindow) => {
      // For example...
      // childWindow.webContents('will-navigate', (e) => {
      //   e.preventDefault()
      // })
    })


  };

export {createSubWindow};
export {WindowOption};
