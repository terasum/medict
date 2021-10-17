<template>
  <div class="container" style="height: 100%">
    <Header :displaySearchBox="false" />
    <div class="main-content-container">
      <div class="left-sidebbar-container">
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

      <div class="word-content-continer">
        <div class="word-content-header">
          <div
            class="header-word-tab"
            :class="
              currentShowWord && currentShowWord.length > 0
                ? 'header-word-tab-with-content'
                : ''
            "
          >
            {{ tabWord }}
          </div>
          <div class="header-search">
            <input
              class="header-search-input"
              type="text"
              placeholder="resource key..."
              v-model="lookupResourceKey"
            />
            <button class="button header-search-btn" @click="onLookupResource">
              <span class="icon">
                <i class="fas fa-search"></i>
              </span>
            </button>
          </div>

          <button class="button header-btn header-btn-devtool" @click="onDevtoolBtnClick">
            <span class="icon">
              <i class="fas fa-bug"></i>
            </span>
          </button>

          <div class="header-btn header-btn-devtool" @click="onResourceDir">
            <span class="icon">
              <i class="fas fa-folder"></i>
            </span>
          </div>
          <div class="header-dict-info">
            <span>{{ currentDict.name }}</span>
          </div>
        </div>
        <div class="word-content">
          <!-- webpreferences="allowRunningInsecureContent=yes" -->
          <webview
            :src="'data:text/html;charset=utf-8;base64,' + currentContent"
            enableremotemodule="true"
            webpreferences="nodeIntegration=false,webSecurity=true,allowRunningInsecureContent=false,contextIsolation=true"
            :preload="preload"
          />
        </div>
      </div>
    </div>
    <FooterBar />
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import Header from '../components/Header.vue';
import FooterBar from '../components/FooterBar.vue';
import { AsyncMainAPI, SyncMainAPI } from '../service.renderer.manifest';
import { listeners } from '../service.renderer.listener';
import Store from '../store/index';

export default Vue.extend({
  components: { Header, FooterBar },
  data() {
    return {
      lookupResourceKey: '',
    };
  },
  computed: {
    preload() {
      return `file://${SyncMainAPI.syncGetWebviewPreliadFilePath()}`;
    },
    currentWordIdx() {
      return (this.$store as typeof Store).state.sideBarData.selectedWordIdx;
    },
    currentShowWord() {
      return (this.$store as typeof Store).state.currentLookupWord;
    },
    tabWord() {
      const searchWord = (this.$store as typeof Store).state.currentLookupWord;
      if (!searchWord) {
        return '';
      }

      // const actualWord = (this.$store as typeof Store).state.currentActualWord;
      // if (!actualWord || actualWord === '') {
      //   return '';
      // }

      // if (actualWord === searchWord) {
      //   return actualWord;
      // }

      // return searchWord + ' › ' + actualWord;
      return searchWord;
    },
    currentContent() {
      return (this.$store as typeof Store).state.currentContent;
    },
    currentDict() {
      return (this.$store as typeof Store).state.currentSelectDict;
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
    onLookupResource() {
      if (!this.lookupResourceKey || this.lookupResourceKey == '') {
        return;
      }
      if (
        !this.currentDict ||
        !this.currentDict.id ||
        this.currentDict.id === ''
      ) {
        return;
      }

      AsyncMainAPI.loadDictResource({
        dictid: this.currentDict.id,
        resourceKey: this.lookupResourceKey,
      });
    },
    onResourceDir() {
      AsyncMainAPI.openDictResourceDir(this.currentDict.id);
    },
  },
  mounted() {
    // webview's content update, this listener
    // designed for @@ENTRY_LINK==
    const webview = document.getElementsByTagName('webview')[0];
    webview.addEventListener('ipc-message', (event) => {
      console.log('====== webview post event ========');
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
        // @ts-ignore
        this.$store.commit('updateCurrentLookupWord', event.args[0].keyText);
      }

      // @ts-ignore
      if (event.channel === 'entryLinkWord') {
        console.log(`[async:mainWindow] webview entryLinkWord clicked`);
        console.log(event);
      }
    });

    // onloadDictResource listener
    listeners.onLoadDictResource((event, arg) => {
      console.log(arg);
    });
  },
  destroyed() {},
});
</script>

<style lang="scss" scoped>
.main-content-container {
  display: flex;
  height: -webkit-calc(100% - 80px);
  height: -moz-calc(100% - 80px);
  height: calc(100% - 80px);
  width: 100%;
}

.left-sidebbar-container {
  width: 160px;
  background-color: #f2f4f5;
  padding-left: 0;
  padding-right: 0;
  border-right: 1px solid #d1d1d1;
  height: 100%;
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
    overflow-y: auto;
    padding-bottom: 10px;
    height: 100%;
    .left-sidebar-wordlist {
      list-style: none;
      margin: 0;
      padding: 0;
      height: auto;
      li {
        padding: 0.5rem 0.2rem 0.5rem 0.4rem;
        border-bottom: 1px solid #e1e1e1;
        font-size: 0.9rem;
        &:hover {
          background: #4A8EFF;
          color: #fff;
        }
        &:active {
          background: #4A8EFF;
          color: #fff;
        }
      }
      .word-item-active {
        color: #fff;
        background: #4A8EFF;
      }
    }
  }
}

.word-content-continer {
  padding: 0;
  height: 100%;
  width: calc(100% - 160px);
  overflow-y: hidden;
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

.word-content-header {
  height: 26px;
  padding-left: 5px;
  background-color: #f5f5f5;
  border-bottom: 1px solid #d1d1d1;
  // border-bottom: 1px solid #ccc;
  .header-word-tab {
    display: inline-block;
  }
  .header-word-tab-with-content {
    padding-left: 3px;
    padding-right: 3px;
    border-top: 1px solid #ccc;
    border-left: 1px solid #ccc;
    border-right: 1px solid #ccc;
    height: 24px;
    margin-top: 2px;
    background: #fff;
    border-radius: 3px 3px 0px 0px;
    border-bottom: #fff;
  }

  .header-search {
    float: right;
    display: flex;
    font-size: 12px;

    .header-search-input {
      display: flex;
      height: 20px;
      margin: 3px 5px;
    }

    .header-search-btn {
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

.header-dict-info {
  float: right;
  display: flex;
  height: 20px;
  margin: 3px 5px;
  min-width: 50px;
  background: #e9e9e9;
  justify-content: center;
  border-radius: 2px;
  span {
    font-weight: 400;
    padding: 0 4px;
    color: #666;
    font-size: 12px;
    line-height: 20px;
    -webkit-user-select: none;
  }
}


webview {
  height: -webkit-calc(100% - 20px);
  height: -moz-calc(100% - 20px);
  height: calc(100% - 20px);
  width: 100%;
  overflow-y: auto;
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
