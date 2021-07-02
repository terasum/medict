<template>
  <div class="container-fluid" style="height: 100%">
    <Header :displaySearchBox="false" />
    <div class="row" style="height: 100%">
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
            style="height: 100%; width: 100%; overflow-y: scroll"
            :preload="preload"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import Header from '../components/Header.vue';
import { listeners } from '../service.renderer.listener';
import { MainProcAsyncAPI } from '../service.renderer.manifest';
import Store from '../store/index';
import fs from 'fs';
import path from 'path';
// @ts-ignore
import tmp from 'tmp';

let tempPreloadPath = '';

(function init() {
  // preload 中定义了点击后处理的 message 逻辑
  // 以及 main-process 返回之后的监听逻辑
  const rawPreloadScript = fs.readFileSync(
    path.resolve('./src/renderer/preload/webview.preload.js')
  );
  const preloadScript = rawPreloadScript.toString('utf8');
  const tmpfile = tmp.fileSync({
    mode: 0o644,
    prefix: 'mdict',
    postfix: '.js',
  });
  tempPreloadPath = tmpfile.name;
  console.log('File: ', tmpfile.name);
  if (fs.existsSync(tempPreloadPath)) {
    fs.unlinkSync(tempPreloadPath);
  }
  fs.writeFileSync(tempPreloadPath, preloadScript);
})();

export default Vue.extend({
  components: { Header },
  data: () => {
    return {
      currentContent: '',
      preload: `file://${tempPreloadPath}`,
    };
  },
  computed: {
    currentWordIdx() {
      return (this.$store as typeof Store).state.sideBarData.selectedWordIdx;
    },
    suggestWords() {
      return (this.$store as typeof Store).state.suggestWords;
    },
  },
  methods: {
    lookupWord(item: any) {
      // api["findWordPrecisly"](item);
      this.$store.dispatch('FIND_WORD_PRECISLY', item.id);
    },
    findResource(dictid: string, resourceKey: string) {
      return MainProcAsyncAPI.loadDictResource({ dictid, resourceKey });
    },
    onDevtoolBtnClick() {
      // for webview
      const webview = document.getElementsByTagName('webview')[0];
      //@ts-ignore
      webview.openDevTools();
    },
  },
  mounted() {
    listeners.onSuggestWord((event: any, args: any) => {
      console.log(`[async:mainWindow] response suggest words:`);
      console.log(event);
      console.log(args);
      this.$store.dispatch('ASYCN_UPDATE_SIDE_BAR', {
        candidateWordNum: args.length,
      });
      this.$store.commit('suggestWords', args);
    });

    listeners.onFindWordPrecisly((event: any, args: any) => {
      console.log(`[async:mainWindow] response onFindWordPrecisly:`);
      console.log(args);
      this.currentContent = Buffer.from(args.definition, 'utf8').toString(
        'base64'
      );
      // this.currentContent = Buffer.from("中文测试", "utf8").toString("base64");
    });

    listeners.onLoadDictResource((event: any, args: any) => {
      console.log(`[async:mainWindow] response onLoadDictResource:`);
      console.log(args);
      // this.currentContent = args.definition.trim("\r\n\u0000");
    });
    // for webview
    const webview = document.getElementsByTagName('webview')[0];
    webview.addEventListener('ipc-message', (event) => {
      // 通过event.channel的值来判断webview发送的事件名
      // @ts-ignore
      if (event.channel === 'onFindWordPrecisly') {
        console.log(`[async:mainWindow] response onFindWordPrecisly:`);
        console.log(event);
        this.currentContent = Buffer.from(
          // @ts-ignore
          event.args[0].definition,
          'utf8'
        ).toString('base64');
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
</style>
