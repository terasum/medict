import { createByProc } from '@terasum/electron-call';
import WorkerMessageAPI from './apis/WorkerMessageApi';
const stubByWorker = createByProc('worker', 'error');
const workerMessageApi = new WorkerMessageAPI();

stubByWorker.provide(['main','renderer'], 'WorkerMessageAPI', workerMessageApi);