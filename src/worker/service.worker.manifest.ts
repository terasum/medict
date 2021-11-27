import { createByProc } from '@terasum/electron-call';
import WorkerMessageAPI from './apis/WorkerMessageAPI';
import DictAPI from './apis/DictAPI';

const stubByWorker = createByProc('worker', 'error');

stubByWorker.provide(['main','renderer'], 'WorkerMessageAPI', new WorkerMessageAPI());
stubByWorker.provide(['main','renderer'], 'DictAPI', new DictAPI());
