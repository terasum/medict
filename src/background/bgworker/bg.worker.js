import JSMdict from 'js-mdict'
import fs from 'fs-extra'
import registerPromiseWorker from 'promise-worker/register'
import CommuniMsg from '../../common/CommuiMsg'
import msgType from '../../common/msgType'

let MDDDict = {
  jsmdict: null,
  loaded: false,
  cache: new Set()
}

registerPromiseWorker(function (message) {
  console.log('bgworker.js recevied message')
  console.log(message)
  if (!message.msgType) {
    return new CommuniMsg(msgType.BGWorkerSubMsgResponse, 'invalid msg type', -4)
  }
  if (message.msgType === msgType.BGWorkerSubMsgLoad) {
    loadMDD(message.data)
    return new CommuniMsg(msgType.BGWorkerSubMsgResponse, 'loaded', 2)
  } else if (message.msgType === msgType.BGWorkerSubMsgQuery) {
    if (!MDDDict.loaded) {
      console.debug('mdd not loaded..')
      return new CommuniMsg(msgType.BGWorkerSubMsgResponse, 'mdd not loaded yet', -1)
    }

    console.log(message.data)
    console.log('query message should be {query: word, name: name}')
    if (!message.data.query || message.data.query === undefined) {
      console.debug('invalid format')
      return new CommuniMsg(msgType.BGWorkerSubMsgResponse, 'invalid format', -2)
    }
    const key = message.data.query
    const name = message.data.name
    // find at cache
    if (MDDDict.cache.has(key)) {
      console.log('cache hit! ignore')
      return new CommuniMsg(msgType.BGWorkerSubMsgResponse, 'hit cache', 1)
    }

    console.log('query for: ' + key)
    const resData = MDDDict.jsmdict.lookup(key)
    // add to cache
    if (resData && resData !== 'NOTFOUND') {
      MDDDict.cache.add(key)
    }
    console.log(resData)
    // fs.writeFileSync(path.join(resourceDirPath, name), resData)
    if (!resData || resData === 'NOTFOUND') {
      return new CommuniMsg(msgType.BGWorkerSubMsgResponse, 'NOTFOUND', -3)
    }
    return new CommuniMsg(msgType.BGWorkerSubMsgResponse, {resdata: resData, query: key, name: name})
  } else {
    return new CommuniMsg(msgType.BGWorkerSubMsgResponse, 'invalid msg type', -4)
  }
})

function loadMDD (mddpath) {
  // TODO change mdd file path
  if (!fs.existsSync(mddpath) || !mddpath.endsWith('.mdd')) {
    console.log('bgworker mdd file path not exist return')
    return
  }
  // const mddpath = path.join(__static, 'dicts/oale8.mdd')
  console.log(mddpath)
  console.log(JSMdict)
  MDDDict.jsmdict = new JSMdict(mddpath)
  MDDDict.loaded = true
  console.log('mdd loaded...')
}

// main worker task
setTimeout(() => {
  self.postMessage('message from bgworker ')
}, 1000)
