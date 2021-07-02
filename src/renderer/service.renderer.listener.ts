import { ipcRenderer } from 'electron';

function wrap(eventName: string) {
  return (callback: (event: any, arg: any) => any) => {
    console.log(`🔧 register listener [${eventName}]`);
    ipcRenderer.on(eventName, (e: any, a: any) => {
      console.log(`------ renderer trigger [${eventName}] -----`);
      console.log(a);
      return callback(e, a);
    });
  };
}

export const listeners = {
  onAsyncMessage: wrap('onAsyncMessage'),
  onCreateSubWindow: wrap('onCreateSubWindow'), // disabled
  onEntryLinkWord: wrap('onEntryLinkWord'),
  onSuggestWord: wrap('onSuggestWord'),
  onFindWordPrecisly: wrap('onFindWordPrecisly'),
  onLoadDictResource: wrap('onLoadDictResource'),
};
