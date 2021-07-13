import { ipcRenderer } from 'electron';

function wrap(eventName: string) {
  return (callback: (event: any, arg: any) => any) => {
    console.log(`🔧 register listener [${eventName}]`);
    ipcRenderer.on(eventName, (e: any, a: any) => {
      return callback(e, a);
    });
  };
}

export const listeners = {
  onCreateSubWindow: wrap('onCreateSubWindow'), // disabled
  onOpenDevTool: wrap('onOpenDevTool'), // disabled
  onEntryLinkWord: wrap('onEntryLinkWord'),
  onSuggestWord: wrap('onSuggestWord'),
  onFindWordPrecisly: wrap('onFindWordPrecisly'),
  onLoadDictResource: wrap('onLoadDictResource'),
};
