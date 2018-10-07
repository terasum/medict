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
  import mt from '../common/msgType'
  
  export default {
    mounted () {
      ipcRenderer.on(mt.MsgToMain, (event, payload) => {
        if (!payload || !payload.msgType) return
        switch (payload.msgType) {
          case mt.SubMsgQueryResponse: {
            // console.log('query response' + payload.data)
            this.$store.dispatch('updateRaw', payload.data)
            break
          }
          case mt.SubMsgQueryListResponse: {
            console.log('update list')
            console.log(payload.data)
            this.$store.dispatch('updateList', payload.data)
            break
          }
          default: {
            console.log('main receive bg message: ')
            console.log(payload)
          }
        }
      })
    },
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
