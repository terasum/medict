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
  syncShowOpenDirDialog: syncWrap('syncShowOpenDirDialog'),
  syncGetResourceServerPort:syncWrap('syncGetResourceServerPort'),
  syncClipboardWriteText:syncWrap('syncClipboardWriteText'),
  syncGetResourceRootPath:syncWrap('syncGetResourceRootPath'),


  syncGetUserResourceRootPath:syncWrap('syncGetUserResourceRootPath'),
  syncGetConfigJsonPath:syncWrap('syncGetConfigJsonPath'),
  syncGetLoggerFilePath:syncWrap('syncGetLoggerFilePath'),
  syncGetWebviewPreliadFilePath:syncWrap('syncGetWebviewPreliadFilePath'),
  syncWritePreloadFile: syncWrap('syncWritePreloadFile'),

};
