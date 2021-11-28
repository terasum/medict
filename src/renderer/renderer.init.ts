import { ipcRenderer } from 'electron';
import { createByProc } from '@terasum/electron-call';
import {DictAPI} from '../worker/apis/DictAPI';
import WorkerMessageAPI from '../worker/apis/WorkerMessageAPI';

console.log('👋 This message is being logged by "renderer.init.ts", included via webpack');



const stubByRenderer = createByProc('renderer', 'error');

const workerMessageApi = stubByRenderer.use<WorkerMessageAPI>('worker', 'WorkerMessageAPI');
const workerDictApi = stubByRenderer.use<DictAPI>('worker', 'DictAPI');

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

