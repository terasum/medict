<template>
  <div id="app">
    <router-view></router-view>
  </div>
</template>

<script>
  import '@/assets/styles/App.scss'
  import '@/assets/css/photon.min.css'
  import 'iview/dist/styles/iview.css'
  import { ipcRenderer } from 'electron'
  import MsgType from '../util/constant.js'

// import SuiVue from 'semantic-ui-vue'
// import 'semantic-ui-css/semantic.min.css'

  setTimeout(() => {
    ipcRenderer.send(MsgType.MsgToBackground, 'send to background')
  }, 3000)

  ipcRenderer.on(MsgType.MsgToMain, (event, payload) => {
    console.log('main receive bg message: ' + payload)
  })

  export default {
    name: 'medict'
  }
</script>

<style>
html,
body {
    height: 100%;
}
 #app{
   height: 100%;
 }
</style>
