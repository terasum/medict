<template lang="html">
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
          <!-- 避免array的一些坑，所以使用对象进行处理 -->
          <li v-for="(item, index) in wordlist" :key="index" v-model="wordlist[index]" :value="item" @click="itemclick(item)">
            {{ item }}
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
          <!-- <AutoComplete v-model="word" @keyup.enter.native="submit" :data="word_data" @on-search="autoComplete" popup placeholder="search..." icon="ios-search" style="width:350px"></AutoComplete> -->
          <!-- <AutoComplete v-model="word" @keyup.enter.native="submit" :data="word_data" @on-search="autoComplete" popup placeholder="search..." icon="ios-search" style="width:350px"></AutoComplete> -->
          <Input v-model="word" @keyup.enter.native="submit" placeholder="Search..." :disabled="inputDisabled" icon="ios-search" style="width: 300px" />

        </div>
        <!-- me-btn-menu starts -->
        <div class="me-setting-menu">
          <ButtonGroup>
            <Button icon="ios-arrow-back"></Button>
            <Button icon="ios-arrow-forward"></Button>
            <!-- <Button icon="ios-switch" @click="toSetting"></Button> -->
            <Button :type="settingType" :icon="settingIcon" @click="clickSetting"></Button>
          </ButtonGroup>
        </div> 
        <!-- me-btn-menu ends -->
      </div>
      <div class="me-main-view">
        <router-view></router-view>
      <!-- <webview id="mainContent" :src="contentData" disablewebsecurity nodeintegration style="display:inline-flex; width:557px; height:452px" webpreferences="javascript=yes" ></webview> -->
      <!-- <iframe :src="contentData" width="640" height="480" style="border:0" ></iframe> -->
      </div>
    </div>
    <!-- me-main ends -->
</div>
</template>


<style lang="scss" scoped>

.wrapper{
  /* height should not larger than 451 px, plus 1px header border*/ 
  height: 451px;
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
      .me-main-view {
        height: 100%;
        width: 100%;
        overflow-y: auto;
        padding-bottom: 50px;
      }
  }
}
</style>


<script>
import Worker from '../../worker/main.worker.js'
// import DOMBuilder from '../../util/dombuiler.js'
import { ipcRenderer } from 'electron'
import mt from '../../common/msgType'
import CommuniMsg from '../../common/CommuiMsg'

// import path from 'path'

let innerDicts = ['朗文中字', '新世纪英汉大字典']
let defaultDict = innerDicts[0]

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
    this.$router.push('/webview')
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
      wordDefinition: '',
      // setting btn
      settingType: 'default',
      settingIcon: 'ios-switch',
      settingStatus: false,
      inputDisabled: false
    }
  },
  // 计算属性，渲染最终webview中的内容
  computed: {
    wordlist () {
      return this.$store.state.Query.samelist
    }
  },
  methods: {
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
        return
      }
      this.tempWord = word
      this.word = word
      ipcRenderer.send(mt.MsgToBackground, new CommuniMsg(mt.SubMsgQueryBackground, word))
    },
    clickSetting () {
      if (this.settingStatus) {
        this.settingIcon = 'ios-switch'
        this.settingType = 'default'
        this.$router.push('/webview')
      } else {
        this.settingIcon = 'md-close'
        this.settingType = 'warning'
        this.$router.push('/setting')
      }
      this.settingStatus = !this.settingStatus
      this.inputDisabled = !this.inputDisabled
      // console.log('setting...')
      // this.$router.push('setting')
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
    word: (newWord, oldWord) => {
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

