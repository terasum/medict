import { ipcRenderer } from 'electron';

function syncWrap(fnName: string) {
  return (args: any) => {
    return ipcRenderer.sendSync(fnName, args);
  };
}

function asyncWrap(fnName: string) {
  return (args: any) => {
    ipcRenderer.send(fnName, args);
  };
}

export const AsyncMainAPI = {
  asyncMessage: asyncWrap('asyncMessage'),
  createSubWindow: asyncWrap('createSubWindow'),
  entryLinkWord: asyncWrap('entryLinkWord'),
  suggestWord: asyncWrap('suggestWord'),
  findWordPrecisly: asyncWrap('findWordPrecisly'),
  loadDictResource: asyncWrap('loadDictResource'),
};

export const SyncMainAPI = {
  syncShowOpenDialog: syncWrap('syncShowOpenDialog'),
  syncMessage: syncWrap('syncMessage'),
  dictAddOne: syncWrap('dictAddOne'),
  dictFindOne: syncWrap('dictFindOne'),
  dictDeleteOne: syncWrap('dictDeleteOne'),
  dictFindAll: syncWrap('dictFindAll'),
};
