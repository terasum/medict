import { BrowserWindow, shell } from 'electron';

interface WindowOption {
  width: number,
  height: number,
  html: string,
  titleBarStyle?: ('default' | 'hidden' | 'hiddenInset' | 'customButtonsOnHover'),
  nodeIntegration: boolean,
  contextIsolation: boolean,
  show: boolean
}

const createSubWindow = (parent: BrowserWindow | undefined, options: WindowOption): BrowserWindow => {
  // Create the browser window.
  
  const subWindow = new BrowserWindow({
    parent: parent,
    height: options.height,
    width: options.width,
    titleBarStyle: options.titleBarStyle,
    // The lines below solved the issue
    webPreferences: {
      nodeIntegration: options.nodeIntegration,
      contextIsolation: options.contextIsolation,
    },
    show: options.show
  });

  // and load the index.html of the app.
  subWindow.loadURL(options.html);

  // Open the DevTools.
  subWindow.webContents.openDevTools();

  // sub-windows
  // mainWindow.webContents.setWindowOpenHandler(({ e, url }) => {
  //   return { action: 'allow' }
  // })

  subWindow.webContents.on('will-navigate', (event: Event, url: string) => {
    /* If url isn't the actual page */
    if (url != subWindow.webContents.getURL()) {
      event.preventDefault();
      shell.openExternal(url);
    }
  });

  subWindow.webContents.on('did-create-window', (childWindow) => {
    // For example...
    // childWindow.webContents('will-navigate', (e) => {
    //   e.preventDefault()
    // })
  })

  return subWindow;
};

export { createSubWindow };
export { WindowOption };
