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

export const MainProcAsyncAPI = {
  asyncMessage: asyncWrap('asyncMessage'),
  createSubWindow: asyncWrap('createSubWindow'),
  entryLinkWord: asyncWrap('entryLinkWord'),
  suggestWord: asyncWrap('suggestWord'),
  findWordPrecisly: asyncWrap('findWordPrecisly'),
  loadDictResource: asyncWrap('loadDictResource'),
};

export const MainProcSyncAPI = {
  syncMessage: syncWrap('syncMessage'),
};
