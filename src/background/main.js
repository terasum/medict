import { ipcRenderer } from 'electron'
import MsgType from '../util/constant'

console.log('backgound.js')

setInterval(() => {
  ipcRenderer.send(MsgType.MsgToMain, 'send to main')
}, 1000)

ipcRenderer.on(MsgType.MsgToBackground, (event, payload) => {
  console.log('background receive main message: ' + payload)
})
