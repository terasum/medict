import { ipcRenderer } from 'electron'
import { dictService } from './mainsvc/main.service'
import { dictContentService } from './mainsvc/dict-content.service'
import _ from 'lodash'

export const asyncfn = {
  asyncMessage: (event: any, arg: any) => {
    event.sender.send('asynchronous-reply', arg)
  },
  createSubWindow: (event: any, arg: any) => {
    event.sender.send('createSubWindow', arg)
  },
  asyncSearchWord: (event: any, arg: any) => {
    // main process handled and re-notify render process
    console.log(
      "[main-process] asyncSearchWord event.sender.send('onAsyncSearchWord')",
    )
    event.sender.send('onAsyncSearchWord', {
      response: 'main-process resp',
      arg,
    })
  },
  suggestWord: (event: any, arg: any) => {
    console.log("[main-process] suggestWord event.sender.send('onSuggestWord')")
    const result = dictService.associate(arg.word)
    event.sender.send('onSuggestWord', result)
  },
  findWordPrecisly: (event: any, arg: any) => {
    console.log(
      "[main-process] suggestWord event.sender.send('onFindWordPrecisly')",
    )
    const result = dictService.findWordPrecisly(
      arg.dictid,
      arg.keyText,
      arg.recordStartOffset,
    )

    const resFn = (resKey:string) =>{
        return dictService.loadDictResource(arg.dictid, resKey);
    }

    event.sender.send('onFindWordPrecisly', {
      keyText: arg.keyText,
      definition: dictContentService.definitionReplace(
        arg.dictid,
        result.definition,
        resFn,
      ),
    })
  },
  loadDictResource: (event: any, arg: any) => {
    const result = dictService.loadDictResource(arg.dictid, arg.resourceKey)
    event.sender.send('onLoadDictResource', result)
  },
}

export const asyncfnListener = {
  onAsyncSearchWord: (callback: (event: any, arg: any) => any) => {
    ipcRenderer.on('onAsyncSearchWord', (e: any, a: any) => {
      console.log('------ renderer listener[onAsyncSearchWord]-----')
      console.log(a)
      return callback(e, a)
    })
  },
  onSuggestWord: (callback: (event: any, arg: any) => any) => {
    ipcRenderer.on('onSuggestWord', (e: any, a: any) => {
      console.log('------ renderer listener[onSuggestWord]-----')
      console.log(a)
      return callback(e, a)
    })
  },
  onFindWordPrecisly: (callback: (event: any, arg: any) => any) => {
    ipcRenderer.on('onFindWordPrecisly', (e: any, a: any) => {
      console.log('------ renderer listener[onFindWordPrecisly]-----')
      console.log(a)
      return callback(e, a)
    })
  },
  onLoadDictResource: (callback: (event: any, arg: any) => any) => {
    ipcRenderer.on('onLoadDictResource', (e: any, a: any) => {
      console.log('------ renderer listener[onLoadDictResource]-----')
      console.log(a)
      return callback(e, a)
    })
  },
}

export const syncfn = {
  syncMessage: (arg: any) => {
    console.log(arg)
    return 'pong'
  },
  lookupWordExctly: (arg: any) => {
    return dictService.lookup(arg.dictid, arg.word)
  },
}
