import { ipcRenderer } from 'electron'
// import path from 'path'
import fs from 'fs'
import Mdict from 'mdict'
import MdictQuery from './MdictQuery'

import mt from '../common/msgType'
import CommuniMsg from '../common/CommuiMsg'
import YDDict from '../util/yd'
import YDQuery from './YDQuery'

import PromiseWorker from 'promise-worker'
import Worker from './bgworker/bg.worker.js'

import Store from 'electron-store'

const store = new Store()

console.log('========= backgound.js ==========')
store.set('bgtest', 'bgtest')
console.log(store.get('bgtest'))

const youdao = new YDDict()

// worker
let worker = new Worker()
let promiseWorker = new PromiseWorker(worker)

/**************************
 * send message functinos *
 **************************/
// eslint-disable-next-line
function sendQueryRespToMain (def) {
  ipcRenderer.send(mt.MsgToMain, new CommuniMsg(mt.SubMsgQueryResponse, def))
}

function sendQueryListRespToMain (list) {
  console.log('sendwords to main')
  console.log(list.map((w) => { return w.toString() }))
  ipcRenderer.send(mt.MsgToMain, new CommuniMsg(mt.SubMsgQueryListResponse, list.map((w) => { return w.toString() })))
}

// eslint-disable-next-line
function sendQueryToWorker (res) {
  return promiseWorker.postMessage(new CommuniMsg(mt.BGWorkerSubMsgQuery, res, 0))
}

function sendLoadToWorker (mddpath) {
  return promiseWorker.postMessage(new CommuniMsg(mt.BGWorkerSubMsgLoad, mddpath, 0))
}

/*************************
 *   message listener    *
 *************************/

/**
 * listening to message main send to background
 */
ipcRenderer.on(mt.MsgToBackground, (event, payload) => {
  console.log('bgmain receive from main ')
  console.log(payload)
  if (!payload || !payload.msgType) return
  let word = payload.data

  switch (payload.msgType) {
    case mt.SubMsgQueryBackground: {
      searchMdict(word)
      break
    }
    default: {
      console.log('background receive main message: ' + payload)
    }
  }
})

/**
 * search by mdict
 * @param {string} word which word want's to query
 * @param {*} dict the dictionary object
 * @param {*} cb callback function(definition)
 */
function searchMdict (word) {
  const spath = store.get('mdx')
  if (!fs.existsSync(spath)) {
    console.log('bg mdx not exist')
    return
  }
  // let spath = path.join(__static, '/dicts/oale8.mdx')
  const stat = fs.statSync(spath)
  if (stat.isFile()) {
    Mdict.dictionary(spath).then((dictionary) => {
      // dictionary is loaded
      dictionary.search({
        phrase: word, // '*' and '?' supported
        max: 10 // maximum results
      }).then(function (foundWords) {
        console.log('Found words:')
        console.log(foundWords) // foundWords is array

        sendQueryListRespToMain(foundWords)

        var fword = '' + foundWords[0]
        console.log('Loading definitions for: ' + fword)
        return dictionary.lookup(fword) // typeof word === string
      }).then(function (definitions) {
        console.log('definitions:') // definition is array

        const mq = new MdictQuery(definitions[0], promiseWorker)
        mq.wrapper()

        // console.log(mq.serialize())
        let def = mq.serialize()
        sendQueryRespToMain(def)
      })
    })
  }
}

/**
 * use youdao to query word
 * @param {string} word which word want's to query
 * @param {*} cb callback
 */
function searchYouDao (word, cb) { // eslint-disable-line
  youdao.lookup(word.trim()).then((def) => {
    const youdaoQuery = new YDQuery(youdao.formatHTML(def))
    // wrapper basic style and script
    youdaoQuery.wrapper()
    cb(youdaoQuery.serialize())
  })
}

/**
 * worker event and multi-thread handler
 */
function multiThread () {
  console.log('!!!! start test worker !!!!')
  // var worker = new Worker('/worker.js') // eslint-disable-line no-undef
  setTimeout(() => {
    sendQueryToWorker('querymsg').then((resp) => {
      console.log('bg main receive')
      console.log(resp)
    })
  }, 20000)
}
// start mutithread
multiThread()
// start load work
const mddpath = store.get('mdd')
if (mddpath && fs.existsSync(mddpath) && mddpath.endsWith('.mdd')) {
  sendLoadToWorker(mddpath)
}
