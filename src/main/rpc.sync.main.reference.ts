// this file is reference by RENDERER process
// do not import in main-process
import {ipcRenderer} from 'electron';

function syncWrap(fnName: string) {
  return (args?: any) => {
    return ipcRenderer.sendSync(fnName, args);
  };
}

export const SyncMainAPI = {
  syncShowOpenDialog: syncWrap('syncShowOpenDialog'),
  syncShowMainLoggerPath: syncWrap('syncShowMainLoggerPath'),
  syncGetResourceRootPath: syncWrap('syncGetResourceRootPath'),
  syncGetResourceServerPort: syncWrap('syncGetResourceServerPort'),
  syncGetWebviewPreliadFilePath: syncWrap('syncGetWebviewPreliadFilePath'),

  clipboardWriteText:syncWrap('clipboardWriteText'),
};
