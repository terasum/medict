<template>
    <webview id="mainContent" :src="contentData" disablewebsecurity nodeintegration style="display:inline-flex; width:557px; height:452px" webpreferences="javascript=yes" ></webview>
</template>

<script>
const {shell} = require('electron').remote
export default {
  mounted () {
  // for webview
    const webview = document.getElementById('mainContent')
    webview.addEventListener('new-window', (e) => {
      console.log('new window event called')
      const protocol = require('url').parse(e.url).protocol
      if (protocol === 'http:' || protocol === 'https:') {
        shell.openExternal(e.url)
      }
    })
    // if (process.env.NODE_ENV === 'development') {
    webview.addEventListener('dom-ready', (e) => {
      console.log('dom ready')
      if (process.env.NODE_ENV === 'development') {
        webview.openDevTools()
      }
    })
    // }
  },
  computed: {
    contentData () {
      return this.$store.state.Query.rawdef
    }
  }
}
</script>

