import Query from './Query'
import fs from 'fs-extra'
import path from 'path'
import electron from 'electron'
import CommuniMsg from '../common/CommuiMsg'
import msgType from '../common/msgType'
import Store from 'electron-store'

const userDataPath = electron.remote.app.getPath('userData')
// const userDataPath = electron.remote.app.getPath('temp')
const resourceDataPath = path.join(userDataPath, 'resources')
const store = new Store()
if (!fs.existsSync(resourceDataPath)) {
  fs.mkdirSync(resourceDataPath)
}

class MdictQuery extends Query {
  constructor (def, worker) {
    super('medict', def, undefined)
    this.worker = worker
  }

  wrapper () {
    super.decorate()
    // oalecss
    // TODO change this url
    const csspath = store.get('css')
    const jspath = store.get('js')
    console.log(csspath)
    console.log(jspath)
    if (!fs.existsSync(csspath)) {
      console.log('css not exist:' + csspath)
      return
    }
    const oalecss = fs.readFileSync(csspath)
    this.content.addStyleContent(oalecss.toString('utf8'))

    if (!fs.existsSync(jspath)) {
      console.log('js not exist:' + jspath)
      return
    }
    const oalejs = fs.readFileSync(jspath)
    this.content.addScriptContent(oalejs.toString('utf8'))
  }

  filter () {
    let tasks = []
    // jsdom's DOM
    const dom = this.content.dom
    // filter links
    let links = dom.window.document.getElementsByTagName('LINK')
    // console.log(links)
    for (let i = 0; i < links.length; i++) {
      let ele = links[i]
      if (ele.hasAttribute('data-medict')) {
        continue
      }
      // 删除该节点
      ele.parentNode.removeChild(ele)
    }

    // filter scripts
    let scripts = dom.window.document.getElementsByTagName('SCRIPT')
    for (let i = 0; i < scripts.length; i++) {
      let ele = scripts[i]
      if (ele.hasAttribute('data-medict')) {
        continue
      }
      if (ele.hasAttribute('src')) {
        // console.log(ele)
        // console.log(ele.getAttribute('src'))
        // ele.parentNode.removeChild(ele)
        let old = ele.getAttribute('src')
        let resq = old
        if (!old.startsWith('/') && !old.startsWith('\\')) {
          resq = '\\' + old
        }
        resq = resq.replace(/\//g, '\\')
        const task = this.worker.postMessage(new CommuniMsg(msgType.BGWorkerSubMsgQuery, {query: resq, name: old}))
        const filepath = path.join(resourceDataPath, old)
        const dirname = path.dirname(filepath)
        if (!fs.existsSync(dirname)) {
          fs.mkdirpSync(dirname)
        }
        task.then((rawdata) => {
          // console.log('query script')
          // console.log(rawdata)
          if (rawdata.code === 0) {
            fs.writeFileSync(filepath, Buffer.from(rawdata.data.resdata))
            console.log('writefile: ' + filepath)
          }
        })
        tasks.push(task)
        if (old.indexOf('jquery') >= 0) {
          ele.parentNode.removeChild(ele)
          continue
        }
        ele.setAttribute('src', 'file://' + filepath)
      }
    }

    // todo filter images
    let imgs = dom.window.document.getElementsByTagName('IMG')
    for (let i = 0; i < imgs.length; i++) {
      let ele = imgs[i]
      if (ele.hasAttribute('data-medict')) {
        continue
      }
      if (ele.hasAttribute('src')) {
        // console.log(ele)
        // console.log(ele.getAttribute('src'))
        // ele.parentNode.removeChild(ele)
        let old = ele.getAttribute('src')
        let resq = old
        if (!old.startsWith('/') && !old.startsWith('\\')) {
          resq = '\\' + old
        }
        resq = resq.replace(/\//g, '\\')
        const task = this.worker.postMessage(new CommuniMsg(msgType.BGWorkerSubMsgQuery, {query: resq, name: old}))
        const filepath = path.join(resourceDataPath, old)

        const dirname = path.dirname(filepath)
        if (!fs.existsSync(dirname)) {
          fs.mkdirpSync(dirname)
        }

        task.then((rawdata) => {
          // console.log('query images')
          // console.log(rawdata)
          if (rawdata.code === 0) {
            // console.log('==== beore write ====')
            // console.log(rawdata.data)
            fs.writeFileSync(filepath, Buffer.from(rawdata.data.resdata))
            // console.log('writefile: ' + filepath)
          }
        })
        tasks.push(task)
        ele.setAttribute('src', 'file://' + filepath)
      }
    }

    // todo filter images
    let alinks = dom.window.document.getElementsByTagName('A')
    for (let i = 0; i < alinks.length; i++) {
      let ele = alinks[i]
      if (ele.hasAttribute('data-medict')) {
        continue
      }
      if (ele.hasAttribute('href')) {
        // console.log(ele)
        // console.log(ele.getAttribute('src'))
        // ele.parentNode.removeChild(ele)
        let old = ele.getAttribute('href')
        let resq = old
        if (old.startsWith('sound:')) {
          old = old.replace('sound://', '')
        } else {
          // skip others
          continue
        }

        if (!old.startsWith('/') && !old.startsWith('\\')) {
          resq = '\\' + old
        }
        resq = resq.replace(/\//g, '\\')
        const task = this.worker.postMessage(new CommuniMsg(msgType.BGWorkerSubMsgQuery, {query: resq, name: old}))
        const filepath = path.join(resourceDataPath, old)

        const dirname = path.dirname(filepath)
        if (!fs.existsSync(dirname)) {
          fs.mkdirpSync(dirname)
        }

        task.then((rawdata) => {
          // console.log('query images')
          // console.log(rawdata)
          if (rawdata.code === 0) {
            // console.log('==== beore write ====')
            // console.log(rawdata.data)
            fs.writeFileSync(filepath, Buffer.from(rawdata.data.resdata))
            // console.log('writefile: ' + filepath)
          }
        })
        tasks.push(task)

        // generate a audio node
        let audio = dom.window.document.createElement('audio')
        let audioSource = dom.window.document.createElement('source')
        audioSource.setAttribute('src', 'file://' + filepath)
        audioSource.setAttribute('type', 'audio/mpeg')
        audio.appendChild(audioSource)
        ele.setAttribute('onclick', `javascript:(function(event){
          event.preventDefault(); 
          if (event.target.parentNode && event.target.parentNode.tagName === 'A'){
            if(event.target.parentNode.lastChild && event.target.parentNode.lastChild.tagName === 'AUDIO'){
              event.target.parentNode.lastChild.play();
            }
          }
        })(event)`)
        ele.appendChild(audio)
        // .audio.appendChild()

        // ele.setAttribute('src', 'file://' + filepath)
      }
    }
    Promise.all(tasks).then(() => {
      console.log('!!!!Query Finished!!!')
    })
  }

  serialize () {
    this.filter()
    return `data:text/html, ` + this.content.serialize()
  }
}

export default MdictQuery
