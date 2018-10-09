import JSMdict from 'js-mdict'
import fs from 'fs-extra'
import registerPromiseWorker from 'promise-worker/register'
import CommuniMsg from '../../common/CommuiMsg'
import msgType from '../../common/msgType'
import Store from 'electron-store'

// console.log('RESOURCE:' + resourceDirPath)
// if (!fs.existsSync(resourceDirPath)) {
//   fs.mkdirSync(resourceDirPath)
//   console.log('created directory')
// }

// import electron from 'electron'
// const userDataPath = electron.remote.app.getPath('userDataPath')

// import os from 'os'
// const path = os.path
// eslint-disable-next-line
let MDDDict = {
  jsmdict: null,
  loaded: false,
  cache: new Set()
}

registerPromiseWorker(function (message) {
  console.log('bgworker.js recevied message')

  if (!MDDDict.loaded) {
    console.debug('mdd not loaded..')
    return new CommuniMsg(msgType.BGWorkerSubMsgResponse, 'mdd not loaded yet', -1)
  }
  console.log(message.data.query)
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
})

function loadMDD () {
  const store = new Store()
  // TODO change mdd file path
  if (!store) return
  const mddpath = store.get('mdd')
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
}, 10000)
loadMDD()
