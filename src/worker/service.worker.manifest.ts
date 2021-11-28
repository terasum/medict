import { createByProc } from '@terasum/electron-call';
import WorkerMessageAPI from './apis/WorkerMessageAPI';
import { TranslateAPI } from './apis/TranslateAPI';
import { DictAPI } from './apis/DictAPI';

const stubByWorker = createByProc('worker', 'error');

stubByWorker.provide(['main','renderer'], 'WorkerMessageAPI', new WorkerMessageAPI());
stubByWorker.provide(['main','renderer'], 'DictAPI', new DictAPI());
stubByWorker.provide(['main','renderer'], 'TranslateAPI', new TranslateAPI());
