<template lang="html">
<div class="me-container">
  <div class="me-header dragable h22px w100p"></div>
  <div class="wrapper">
    <!-- me-sidebar start -->
    <div class="me-sidebar">
      <!-- me-search-container start -->
      <div class="me-select-container">
        <img src="@/assets/images/medict.png" width="180" hegiht="40"/>
      </div>
      <!-- me-search-container ends -->

      <!-- me-search-words-list -->
      <div class="me-search-words-list">
        <ul class="words-list">
          <li v-for="(item, index) in wordlist" :key="index" v-model="wordlist[index]" :value="item" @click="itemclick(item)">
            <!-- TODO -->
            {{ item }}
            <Divider dashed />            
            </li>
        </ul>
      </div>
      <!-- me-search-words-list ends -->
    </div>
    <!-- me-sidebar ends -->

    <!-- me-main ends -->
    <div class="me-main">
      <!-- me-setting-menu starts -->
      <div class="me-settings">
        <div class="me-setting-tool">
          <!-- <Input suffix="ios-search" placeholder="search.." style="width: 240px" /> -->
          <!-- <AutoComplete v-model="word" @keyup.enter.native="submit" :data="word_data" @on-search="autoComplete" popup placeholder="search..." icon="ios-search" style="width:350px"></AutoComplete> -->
          <!-- <AutoComplete v-model="word" @keyup.enter.native="submit" :data="word_data" @on-search="autoComplete" popup placeholder="search..." icon="ios-search" style="width:350px"></AutoComplete> -->
          <!-- <input v-model="word" @keyup.enter.native="submit" type="text"></input> -->
          <Input v-model="word" @keyup.enter.native="submit" placeholder="Enter something..."  icon="ios-search" style="width: 300px" />

        </div>
        <!-- me-setting-menu starts -->
        <div class="me-setting-menu">
            <div class="btn-group">
              <button class="btn btn-large btn-default">
                <span class="icon icon-search" @click="searchClick"></span>
              </button>
              <button class="btn btn-large btn-default">
                <span class="icon icon-cog"></span>
              </button>
              <button class="btn btn-large btn-default">
                <span class="icon icon-info"></span>
              </button>
            </div>
        </div> 
        <!-- me-setting-menu ends -->
      </div>
      <!-- me-settings ends -->
      <webview id="mainContent" :src="contentData" disablewebsecurity nodeintegration style="display:inline-flex; width:557px; height:452px" webpreferences="javascript=yes" ></webview>
      <!-- <iframe :src="contentData" width="640" height="480" style="border:0" ></iframe> -->
    </div>
    <!-- me-main ends -->
  </div>
</div>
</template>


<style lang="scss" scoped>
.me-container{
  height: 473px;
  overflow: hidden;
  .wrapper{
    /* height should not larger than 452 px, plus 1px header border*/ 
    height: 452px;
    margin: 0;
    .me-sidebar{
      float: left;
      padding: 0;
      margin: 0;
      width: 185px;
      height: 100%;
      border-right: 1px solid #eee;
      .me-select-container{
        height: 40px;
        width: 100%;
        padding: 2px 5px;
        border-bottom: 1px solid #eee;
        background-color: #fefefe;
      }
    }

    /* me-search-words-list */
    .me-search-words-list{
      /* search container 40 + 1 px, and sidebar 452px (same as wrapper), so this is 411px */
      height: 411px;
      /* width same as sidebar 185px */
      width: 185px;
      overflow: hidden;
      .words-list{
        height: 411px;
        display: block;
        overflow-y: auto;
        overflow-x: hidden;
        li{
          display: block;
          height: 40px;
          // because of the scollbar (15px), so this width is 185-15 = 170px;
          width:170px;
          // 15px is the same as the scroll-bar's width
          padding-left:15px;
          // border-bottom: 1px solid #eee;
          &:hover{
            background-color: #ccc;
          }
        }

      }
    }

    .me-main{
      float: left;
      width: 557px;
      padding: 0;
      margin: 0; 
      height: 100%;
      // background-color: #1377ed;
      .me-settings{
        height: 40px;
        width: 100%;
        // padding: 2px 5px;
        border-bottom: 1px solid #eee;
        background-color: #fefefe;
        .me-setting-tool{
          width:400px;
          float: left;
          height: 40px;
          line-height: 40px;
          padding-left: 10px;
        }
        .me-setting-menu{
          width: 157px;
          float: right;
          height: 40px;
          line-height: 40px;
          text-align: center;

        }
      }  
    }
  }
  
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


<script>
import Worker from '../../worker/main.worker.js'
// import DOMBuilder from '../../util/dombuiler.js'
import { ipcRenderer } from 'electron'
import mt from '../../common/msgType'
import CommuniMsg from '../../common/CommuiMsg'

const {shell} = require('electron').remote

// import path from 'path'

let innerDicts = ['朗文中字', '新世纪英汉大字典']
let defaultDict = innerDicts[0]
let mdict = null

/**
 * worker event and multi-thread handler
 */
function multiThread () { // eslint-disable-line
  // var worker = new Worker('/worker.js') // eslint-disable-line no-undef
  const worker = new Worker()
  worker.postMessage('message')

  worker.onmessage = function (e) {
    console.log('Got message from Worker: ' + e.data)
    console.log(e.data)
  }
}

export default {
  mounted () {
    window.addEventListener('keyup', (event) => {
      // console.log(event)
      if (event.key && event.key === 'Enter' && event.code === 'Enter') {
        this.search()
      }
    }, true)

    // query test
    setTimeout(() => {
      ipcRenderer.send(mt.MsgToBackground, new CommuniMsg(mt.SubMsgQueryBackground, 'hello'))
      console.log(this.$store.state)
    }, 2000)
    // background tasks ...
    multiThread()

    // for webview
    const webview = document.getElementById('mainContent')
    webview.addEventListener('new-window', (e) => {
      console.log('new window event called')
      const protocol = require('url').parse(e.url).protocol
      if (protocol === 'http:' || protocol === 'https:') {
        shell.openExternal(e.url)
      }
    })
    if (process.env.NODE_ENV === 'development') {
      webview.addEventListener('dom-ready', (e) => {
        console.log('dom ready')
        webview.openDevTools()
      })
    }
  },
  data () {
    return {
      // pre-setted dictionaries
      dicts: innerDicts,
      // current used dictionary name
      selectedDict: defaultDict,
      // AutoComplete word data
      word_data: [],
      // currently actually query word
      word: '',
      // to prevent too quick query
      tempWord: '',
      // current word definitions
      wordDefinition: ''
    }
  },
  // 计算属性，渲染最终webview中的内容
  computed: {
    contentData () {
      return this.$store.state.Query.rawdef
    },
    wordlist () {
      console.log('wordlist get')
      // let samelist = this.$store.state.Query.samelist
      return this.$store.state.Query.samelist
      // return this.$store.state.Query.samelist
    }
  },
  methods: {
    searchClick () {
      console.log('clicked')
      // mdict = new Medict(__static + '/dicts/oale8.mdx')
      mdict
        .then((dict) => {
          let def = dict.lookup('hello')
          this.wordDefinition = def
        })
        .error((msg) => {
          console.error(msg)
        })
    },
    autoComplete (value) {
      // todo prefix or simword service
      this.word_data = !value ? ['a'] : [
        value
      ]
      // this.word_data = !value ? [] : [
      //   value,
      //   value + value,
      //   value + value + value
      // ]
    },
    search () {
      console.log('search')
      if (this.word === '' || this.word === this.tempWord) {
        console.log('same do nothing')
        return
      }
      this.tempWord = this.word
      ipcRenderer.send(mt.MsgToBackground, new CommuniMsg(mt.SubMsgQueryBackground, this.word))
      // console.log(this.$store.state)
    },
    submit () {
      console.log('submit')
      console.log(this.word)
    },
    itemclick (word) {
      if (word === '' || word === this.tempWord) {
        console.log('same do nothing')
        return
      }
      this.tempWord = word
      ipcRenderer.send(mt.MsgToBackground, new CommuniMsg(mt.SubMsgQueryBackground, word))
    },
    start () {
      this.$Loading.start()
    },
    finish () {
      this.$Loading.finish()
    },
    error () {
      this.$Loading.error()
    },
    open (nodesc) {
      this.$Notice.success({
        title: 'Notification title',
        desc: nodesc ? '' : 'Here is the notification description. Here is the notification description. '
      })
    }
  },
  // 侦听器
  watch: {
    // 如果 `word` 发生改变，这个函数就会运行
    word: function (newWord, oldWord) {
      // this.answer = 'Waiting for you to stop typing...'
      // this.debouncedGetAnswer()
      // if (this.tempWord === newWord){
      // }
      if (newWord === '') return
      console.log('word changed..')
      // ipcRenderer.send(mt.MsgToBackground, new mt.CommuniMsg(mt.SubMsgQueryBackground, newWord))
      // console.log(this.$store.state)
    }
  }
}
</script>

