<template>
  <div id="app">
    <div class="me-header dragable h22px w100p"></div>
    <div class="me-container">
      <router-view></router-view>
  </div>
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
            this.$store.commit('UPDATE_LIST', payload.data)
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


<style lang="scss" scoped>
html,
body {
    height: 100%;
}
 #app{
   height: 100%;
 }
.me-container{
  height: 451px;
  overflow: hidden;
}

.dragable{
  display: block;
  // background: #1377ed ;
  -webkit-app-region: drag;
  overflow: hidden;
  border-bottom:1px solid #eee;
  min-height: 22px;
  box-shadow: inset 0 1px 0 #f5f4f5;
  background-color: #e8e6e8;
  background-image: -webkit-gradient(linear,left top,left bottom,color-stop(0,#e8e6e8),color-stop(100%,#d1cfd1));
  background-image: -webkit-linear-gradient(top,#e8e6e8 0,#d1cfd1 100%);
  background-image: linear-gradient(to bottom,#e8e6e8 0,#d1cfd1 100%);
}
.h22px{
  height: 22px;
}
.w100p{
  width: 100%;
}
</style>
