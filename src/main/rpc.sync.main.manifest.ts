import { FileOpenApi } from './apis/FileOpenApi';
import { getResourceServerPort } from '../main/init.resource.server';
// import DictAPI from '../worker/apis/DictAPI';
// import { createByProc } from '@terasum/electron-call';
// const stubByMain = createByProc('main', 'error');
const fileOpenApi = new FileOpenApi();

// const dictApi = stubByMain.use<DictAPI>('worker', 'DictAPI');

export const asyncfn = {
  // entryLinkWord: (event: any, arg: { keyText: string; dictid: string }) {
  //   dictApi.suggestWord(arg.dictid, arg.keyText);
  // },
};

export const syncfn = {
  syncShowOpenDialog: fileOpenApi.syncShowOpenDialog,
  syncShowMainLoggerPath: fileOpenApi.syncShowMainLoggerPath,
  syncGetResourceRootPath: fileOpenApi.syncGetResourceRootPath,
  syncGetWebviewPreliadFilePath: fileOpenApi.syncGetWebviewPreliadFilePath,
  syncGetResourceServerPort: getResourceServerPort,
};
