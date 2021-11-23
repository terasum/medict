import { ipcRenderer } from 'electron';

function syncWrap(fnName: string) {
  return (args?: any) => {
    return ipcRenderer.sendSync(fnName, args);
  };
}

function asyncWrap(fnName: string) {
  return (args?: any) => {
    ipcRenderer.send(fnName, args);
  };
}

export const AsyncMainAPI = {
  createSubWindow: asyncWrap('createSubWindow'),
  openDevTool: asyncWrap('openDevTool'),
  openDictResourceDir: asyncWrap('openDictResourceDir'),
  openResourceDir: asyncWrap('openResourceDir'),
  openMainProcessLog: asyncWrap('openMainProcessLog'),
  openUrlOnBrowser: asyncWrap('openUrlOnBrowser'),
  entryLinkWord: asyncWrap('entryLinkWord'),
  suggestWord: asyncWrap('suggestWord'),
  findWordPrecisly: asyncWrap('findWordPrecisly'),
  loadDictResource: asyncWrap('loadDictResource'),
  asyncBaiduTranslate: asyncWrap('asyncBaiduTranslate'),
  asyncGoogleTranslate: asyncWrap('asyncGoogleTranslate'),
  asyncYoudaoTranslate: asyncWrap('asyncYoudaoTranslate'),
  
};

export const SyncMainAPI = {
  syncShowOpenDialog: syncWrap('syncShowOpenDialog'),
  dictAddOne: syncWrap('dictAddOne'),
  dictFindOne: syncWrap('dictFindOne'),
  dictDeleteOne: syncWrap('dictDeleteOne'),
  dictFindAll: syncWrap('dictFindAll'),
  syncShowMainLoggerPath: syncWrap('syncShowMainLoggerPath'),
  syncGetResourceRootPath: syncWrap('syncGetResourceRootPath'),
  syncUserGetResourceRootPath: syncWrap('syncGetUserResourceRootPath'),
  syncShowComfirmMessageBox: syncWrap('syncShowComfirmMessageBox'),
  loadTranslateApiConfig: syncWrap('loadTranslateApiConfig'),
  saveTranslateBaiduApiConfig: syncWrap('saveTranslateBaiduApiConfig'),
  saveTranslateYoudaoApiConfig: syncWrap('saveTranslateYoudaoApiConfig'),
  syncGetWebviewPreliadFilePath: syncWrap('syncGetWebviewPreliadFilePath'),
  clipboardWriteText:syncWrap('clipboardWriteText'),
};
