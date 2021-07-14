<template>
  <div class="container-fluid" style="height: 100%">
    <Header :displaySearchBox="false" />
    <div class="row" style="height: calc(100%-26px)">
      <div class="col col-2 left-sidebbar-container">
        <div class="left-sidebbar">
          <ul class="left-sidebar-wordlist">
            <li
              v-for="item in suggestWords"
              :key="item.id"
              v-on:click="lookupWord(item)"
              v-bind:class="{ 'word-item-active': currentWordIdx === item.id }"
            >
              {{ item.keyText }}
            </li>
          </ul>
        </div>
      </div>
      <div class="col word-content-continer">
        <div class="word-content-header">
          <div class="header-btn header-btn-devtool" @click="onDevtoolBtnClick">
            devtool
          </div>
        </div>
        <div class="word-content">
          <!-- webpreferences="allowRunningInsecureContent=yes" -->
          <webview
            :src="'data:text/html;charset=utf-8;base64,' + currentContent"
            enableremotemodule="true"
            webpreferences="nodeIntegration=false,webSecurity=true,allowRunningInsecureContent=false,contextIsolation=true"
            style="height: 100%; width: 100%; overflow-y: auto"
            :preload="preload"
          />
        </div>
      </div>
    </div>
    <Footer />
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import Header from '../components/Header.vue';
import { AsyncMainAPI } from '../service.renderer.manifest';
import Store from '../store/index';
import fs from 'fs';
// @ts-ignore
import tmp from 'tmp';

let tempPreloadPath = '';

(function init() {
  // preload 中定义了点击后处理的 message 逻辑
  // 以及 main-process 返回之后的监听逻辑
  // const rawPreloadScript = fs.readFileSync(
  //   path.resolve('./src/renderer/preload/webview.preload.js')
  // );

  const rawPreloadScript = `
const { ipcRenderer } = require('electron');
console.warn('=== preload electron [sandbox] ===');

// 监听 main-process 发回来的 结果，格式是 {keyText:"", definition:""}
ipcRenderer.on('onFindWordPrecisly', (event, args) => {
  console.log('------ webview listener[onFindWordPrecisly] -----');
  console.log(args);
  return ipcRenderer.sendToHost('onFindWordPrecisly', args);
});

// 主要处理点击 entry://之后的逻辑
// 将会把需要查询的词发送到 main-process
window.addEventListener('message', function (event) {
  console.log('---- preload listenning message -----');
  console.log(event.data);
  if (event.data && event.data.channel && event.data.payload) {
    console.log(
      'send to main-process [\${event.data.channel}|\${event.data.payload}]'
    );
    ipcRenderer.send(event.data.channel, event.data.payload);
  }
});
  `;

  const preloadScript = rawPreloadScript;
  const tmpfile = tmp.fileSync({
    mode: 0o644,
    prefix: 'mdict',
    postfix: '.js',
  });
  tempPreloadPath = tmpfile.name;
  console.log('preload file: ', tmpfile.name);
  if (fs.existsSync(tempPreloadPath)) {
    fs.unlinkSync(tempPreloadPath);
  }
  fs.writeFileSync(tempPreloadPath, preloadScript);
})();

export default Vue.extend({
  components: { Header },
  data: () => {
    return {
      preload: `file://${tempPreloadPath}`,
    };
  },
  computed: {
    currentWordIdx() {
      return (this.$store as typeof Store).state.sideBarData.selectedWordIdx;
    },
    currentContent() {
      return (this.$store as typeof Store).state.currentContent;
    },
    suggestWords() {
      return (this.$store as typeof Store).state.suggestWords;
    },
  },
  methods: {
    lookupWord(item: any) {
      this.$store.dispatch('asyncFindWordPrecisly', item.id);
    },
    findResource(dictid: string, resourceKey: string) {
      return AsyncMainAPI.loadDictResource({ dictid, resourceKey });
    },
    onDevtoolBtnClick() {
      // for webview
      const webview = document.getElementsByTagName('webview')[0];
      //@ts-ignore
      webview.openDevTools();
    },
  },
  mounted() {
    // webview's content update, this listener
    // designed for @@ENTRY_LINK==
    const webview = document.getElementsByTagName('webview')[0];
    webview.addEventListener('ipc-message', (event) => {
      // 通过event.channel的值来判断webview发送的事件名
      // @ts-ignore
      if (event.channel === 'onFindWordPrecisly') {
        console.log(`[async:mainWindow] response onFindWordPrecisly:`);
        console.log(event);
        const newContent = Buffer.from(
          // @ts-ignore
          event.args[0].definition,
          'utf8'
        ).toString('base64');
        this.$store.commit('updateCurrentContent', newContent);
      }
    });
  },
  destroyed() {},
});
</script>

<style lang="scss" scoped>
.left-sidebbar-container {
  background-color: #f0e8e9;
  padding-left: 0;
  padding-right: 0;
  height: -webkit-calc(100% - 60px);
  height: -moz-calc(100% - 60px);
  height: calc(100% - 60px);

  ::-webkit-scrollbar-track {
    // -webkit-box-shadow: inset 0 0 6px;
    background-color: #f5f5f5;
  }

  ::-webkit-scrollbar {
    width: 5px;
    background-color: #f5f5f5;
  }

  ::-webkit-scrollbar-thumb {
    border-radius: 10px;
    background-color: #555;
  }
  .left-sidebbar {
    overflow-y: scroll;
    padding-bottom: 10px;
    height: 100%;
    .left-sidebar-wordlist {
      list-style: none;
      margin: 0;
      padding: 0;
      height: auto;
      li {
        padding: 0.2rem 0.2rem 0.2rem 0.2rem;
        border-bottom: 1px solid #c1c1c1;
        font-size: 0.9rem;
        &:hover {
          background: #dbdbdb;
        }
        &:active {
          background: #c0c0c0;
        }
      }
      .word-item-active {
        background: #dbdbdb;
      }
    }
  }
}
.word-content-header {
  height: 26px;
  background-color: #f0e8e9;
  // border-bottom: 1px solid #ccc;
  .header-btn {
    float: right;
    display: block;
    font-size: 12px;
    border: 1px solid #aaa;
    padding: 0 3px 0 3px;
    border-radius: 2px;
    height: 20px;
    line-height: 20px;
    background-color: #f1f1f1;
    margin: 3px 5px;
    text-align: center;
    &:hover {
      border: 1px solid #666;
    }
    &:active {
      background-color: #919191;
    }
  }
}
.word-content-continer {
  padding: 0;
  height: -webkit-calc(100% - 86px);
  height: -moz-calc(100% - 86px);
  height: calc(100% - 86px);
  .word-content {
    padding: 0;
    padding-bottom: 10px;
    height: 100%;
    * {
      user-select: all;
      -webkit-user-select: all;
    }
  }
}
webview {
  ::-webkit-scrollbar-track {
    // -webkit-box-shadow: inset 0 0 6px;
    background-color: #f5f5f5;
  }

  ::-webkit-scrollbar {
    width: 5px;
    background-color: #f5f5f5;
  }

  ::-webkit-scrollbar-thumb {
    border-radius: 10px;
    background-color: #555;
  }
}
</style>
