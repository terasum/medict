import { ipcRenderer } from 'electron';
import { createByProc } from '@terasum/electron-call';
import { MessageApi } from '../main/apis/MessageApi';
import WorkerMessageAPI from '../worker/apis/WorkerMessageApi';

const stubByRenderer = createByProc('renderer', 'error');

const workerMessageApi = stubByRenderer.use<WorkerMessageAPI>('worker', 'WorkerMessageAPI');
const messageApi = stubByRenderer.use<MessageApi>('main', 'MessageApi');

  setInterval(() => {
    console.log('invoke messageApi.syncMessage("hello")')
    messageApi.syncMessage('hello').then((v) => {
      console.log("RRRR", v);
    })
  }, 3000);


(function errorListen() {

  window.onerror = function (error, url, line) {
    ipcRenderer.send('errorInWindow', { error, url, line });
  };

})();

// function syncWrap(fnName: string) {
//   return (args?: any) => {
//     return ipcRenderer.sendSync(fnName, args);
//   };
// }

// function asyncWrap(fnName: string) {
//   return (args?: any) => {
//     ipcRenderer.send(fnName, args);
//   };
// }

// export const AsyncMainAPI = {
//   createSubWindow: asyncWrap('createSubWindow'),
//   openDevTool: asyncWrap('openDevTool'),
//   openDictResourceDir: asyncWrap('openDictResourceDir'),
//   openResourceDir: asyncWrap('openResourceDir'),
//   openMainProcessLog: asyncWrap('openMainProcessLog'),
//   openUrlOnBrowser: asyncWrap('openUrlOnBrowser'),
//   entryLinkWord: asyncWrap('entryLinkWord'),
//   suggestWord: asyncWrap('suggestWord'),
//   findWordPrecisly: asyncWrap('findWordPrecisly'),
//   loadDictResource: asyncWrap('loadDictResource'),
//   asyncBaiduTranslate: asyncWrap('asyncBaiduTranslate'),
//   asyncGoogleTranslate: asyncWrap('asyncGoogleTranslate'),
//   asyncYoudaoTranslate: asyncWrap('asyncYoudaoTranslate'),

// };

// export const SyncMainAPI = {
//   syncShowOpenDialog: syncWrap('syncShowOpenDialog'),
//   dictAddOne: syncWrap('dictAddOne'),
//   dictFindOne: syncWrap('dictFindOne'),
//   dictDeleteOne: syncWrap('dictDeleteOne'),
//   dictFindAll: syncWrap('dictFindAll'),
//   syncShowMainLoggerPath: syncWrap('syncShowMainLoggerPath'),
//   syncGetResourceRootPath: syncWrap('syncGetResourceRootPath'),
//   syncShowComfirmMessageBox: syncWrap('syncShowComfirmMessageBox'),
//   loadTranslateApiConfig: syncWrap('loadTranslateApiConfig'),
//   saveTranslateBaiduApiConfig: syncWrap('saveTranslateBaiduApiConfig'),
//   saveTranslateYoudaoApiConfig: syncWrap('saveTranslateYoudaoApiConfig'),
//   syncGetWebviewPreliadFilePath: syncWrap('syncGetWebviewPreliadFilePath'),
//   clipboardWriteText:syncWrap('clipboardWriteText'),
// };
