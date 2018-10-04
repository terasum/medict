<template lang="html">
<div class="me-container">
  <div class="me-header dragable h22px w100p"></div>
  <div class="wrapper">
    <!-- me-sidebar start -->
    <div class="me-sidebar">
      <!-- me-search-container start -->
      <div class="me-select-container">
        <Select v-model="selectedDict" style="width:175px">
          <Option v-for="dict in dicts" :value="dict" :key="dict">{{ dict }}</Option>
        </Select>
      </div>
      <!-- me-search-container ends -->

      <!-- me-search-words-list -->
      <div class="me-search-words-list">
        <ul class="words-list">
          <li v-for="item in wordlist" :value="item" :key="item">
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
          <AutoComplete v-model="word" :data="word_data" @on-search="autoComplete" placeholder="input here" icon="ios-search" style="width:200px"></AutoComplete>
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
      <webview id="foo" :src="contentUrl" style="display:inline-flex; width:640px; height:480px"></webview>
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
import Medict from 'js-mdict'
import Worker from '../../worker/main.worker.js'

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

/**
 * load dictionary
 */
function loadDict () { // eslint-disable-line
  return new Promise((resolve, reject) => {
    let _medict = new Medict(__static + '/dicts/oale8.mdx')
    resolve(_medict)
  })
}

export default {
  mounted () {
    console.log('loading dictionary...')
    // multiThread()
    // mdict = loadDict()
    // mdict.then(() => {
    //   console.log('loaded.')
    // })
    // background...
    multiThread()
  },
  data () {
    return {
      wordlist: ['word1', 'word2', 'word3', 'word4', 'word5', '2ord6', 'word7', 'word8', '2ord9', 'word19', '2ord13', 'word112', '2ord143'],
      // pre-setted dictionaries
      dicts: innerDicts,
      // current used dictionary name
      selectedDict: defaultDict,
      // AutoComplete word data
      word_data: [],
      // currently actually query word
      word: '',
      // contentUrl: 'http://www.baidu.com',
      // current word definitions
      wordDefinition: ''
    }
  },
  // 计算属性，渲染最终webview中的内容
  computed: {
    contentUrl () {
      // TODO word definition filter
      return 'data:text/html, <html> ' + this.wordDefinition + '</html>'
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
      this.word_data = !value ? [] : [
        value,
        value + value,
        value + value + value
      ]
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
  }
}
</script>

