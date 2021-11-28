import { ipcRenderer } from 'electron';

function syncWrap(fnName: string) {
  return (args?: any) => {
    return ipcRenderer.sendSync(fnName, args);
  };
}

export const SyncMainAPI = {
  syncShowOpenDialog: syncWrap('syncShowOpenDialog'),
  syncShowMainLoggerPath: syncWrap('syncShowMainLoggerPath'),
  syncGetResourceRootPath: syncWrap('syncGetResourceRootPath'),
  syncUserGetResourceRootPath: syncWrap('syncGetUserResourceRootPath'),
  syncShowComfirmMessageBox: syncWrap('syncShowComfirmMessageBox'),
  loadTranslateApiConfig: syncWrap('loadTranslateApiConfig'),
  saveTranslateBaiduApiConfig: syncWrap('saveTranslateBaiduApiConfig'),
  saveTranslateYoudaoApiConfig: syncWrap('saveTranslateYoudaoApiConfig'),
  syncGetWebviewPreliadFilePath: syncWrap('syncGetWebviewPreliadFilePath'),
  clipboardWriteText:syncWrap('clipboardWriteText'),
  syncGetResourceServerPort:syncWrap('syncGetResourceServerPort'),
};
