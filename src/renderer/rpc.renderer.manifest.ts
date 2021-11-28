import { ipcRenderer } from 'electron';
import { createByProc } from '@terasum/electron-call';
import { MessageApi } from '../main/apis/MessageApi';
import WorkerMessageAPI from '../worker/apis/WorkerMessageAPI';
import DictAPI from '../worker/apis/DictAPI';

const stubByRenderer = createByProc('renderer', 'error');

const workerMessageApi = stubByRenderer.use<WorkerMessageAPI>('worker', 'WorkerMessageAPI');
const workerDictApi = stubByRenderer.use<DictAPI>('worker', 'DictAPI');
const messageApi = stubByRenderer.use<MessageApi>('main', 'MessageApi');

(function init() {
  setTimeout(async () => {
    const indexed = await workerDictApi.loadAllIndexed()
    const unindexed = await workerDictApi.loadAllUnIndexed()

    console.log('[RENDERER] list all indexed dicts', indexed);
    console.log('[RENDERER] list all unindexed dicts', unindexed);
  }, 4000)
})();



(function errorListen() {
  window.onerror = function (error, url, line) {
    ipcRenderer.send('errorInWindow', { error, url, line });
  };
})();

function syncWrap(fnName: string) {
  return (args?: any) => {
    return ipcRenderer.sendSync(fnName, args);
  };
}

export const SyncMainAPI = {
  syncShowOpenDialog: syncWrap('syncShowOpenDialog'),
//   dictAddOne: syncWrap('dictAddOne'),
//   dictFindOne: syncWrap('dictFindOne'),
//   dictDeleteOne: syncWrap('dictDeleteOne'),
//   dictFindAll: syncWrap('dictFindAll'),
  syncShowMainLoggerPath: syncWrap('syncShowMainLoggerPath'),
  syncGetResourceRootPath: syncWrap('syncGetResourceRootPath'),
//   syncShowComfirmMessageBox: syncWrap('syncShowComfirmMessageBox'),
//   loadTranslateApiConfig: syncWrap('loadTranslateApiConfig'),
//   saveTranslateBaiduApiConfig: syncWrap('saveTranslateBaiduApiConfig'),
//   saveTranslateYoudaoApiConfig: syncWrap('saveTranslateYoudaoApiConfig'),
  syncGetWebviewPreliadFilePath: syncWrap('syncGetWebviewPreliadFilePath'),
//   clipboardWriteText:syncWrap('clipboardWriteText'),
};
