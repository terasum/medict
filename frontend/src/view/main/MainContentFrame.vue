<!--

 Copyright (C) 2023 Quan Chen <chenquan_act@163.com>

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.

 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.
-->

<style lang="scss" scoped>
@import '@/style/variables.scss';
@import '@/style/photon/photon.scss';

.app-content-main {
  width: 100%;
  height: calc(100% - $layout-header-height);
  .app-content-main-toolbar {
    height: 26px;
    display: flex;
    flex-direction: row-reverse;
    background-color: #f6f8fa;
    box-shadow: inset 0 calc(max(1px, 0.0625rem) * -1) #d0d7de;
    .app-content-main-toolbar-box {
      display: block;
      height: 22px;
      width: 22px;
      border: 1px solid #d1d7dd;
      font-size: 16px;
      text-align: center;
      line-height: 22px;
      margin-left: 3px;
      margin-right: 3px;
      margin-top: 2px;
      border-radius: 3px;
      background-color: #f6f8fa;
      color: #596059;
      svg {
        cursor: pointer;
      }
    }
  }
  #app-content-main-iframe-wrapper {
    height: calc(100% - 16px);
    padding: 8px 4px;
    .app-content-main-iframe {
      width: 100%;
      height: 100%;
    }
  }
}
</style>
<template>
  <div class="app-content-main">
    <div class="app-content-main-toolbar">
      <span class="app-content-main-toolbar-box" @click="todo"
        ><NIcon><Bug16Regular /></NIcon
      ></span>
      <span class="app-content-main-toolbar-box" @click="todo"
        ><NIcon><DocumentCss20Regular /></NIcon
      ></span>
      <span class="app-content-main-toolbar-box" @click="refresh"
        ><NIcon><ArrowClockwise20Filled /></NIcon
      ></span>

      <span class="app-content-main-toolbar-box" @click="zoomIn"
        ><NIcon><ZoomOut16Regular /></NIcon
      ></span>
      <span class="app-content-main-toolbar-box" @click="zoomOut"
        ><NIcon><ZoomIn16Regular /></NIcon
      ></span>
    </div>
    <div id="app-content-main-iframe-wrapper"></div>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, onUnmounted } from 'vue';
import { useDictQueryStore } from '@/store/dict';
import { ZoomIn16Regular, ZoomOut16Regular,ArrowClockwise20Filled, Bug16Regular, DocumentCss20Regular } from '@vicons/fluent';
import { NIcon } from 'naive-ui';
import { useMessage } from 'naive-ui';

const dictQueryStore = useDictQueryStore();
const message = useMessage();

const TOP_WIN_MSG_ZOOM_OUT =  '__Medict_TOP_WIN_MSG_EVTP_ZOOM_OUT';
const TOP_WIN_MSG_ZOOM_IN =  '__Medict_TOP_WIN_MSG_EVTP_ZOOM_IN';
const TOP_WIN_MSG_REFRESH = '__Medict_TOP_WIN_MSG_EVTP_REFRESH';
const TOP_WIN_MSG_SETUP =  '__Medict_TOP_WIN_MSG__EVTY_SETUP__';
const INNER_FRAME_MSG_ENTRY_JUMP = '__Medict_INNER_FRAME_MSG_EVTP_ENTRY_JUMP';



function cerateIframe() {
  if (document.getElementById('app-content-main-iframe')) {
    document.getElementById('app-content-main-iframe').remove();
  }
  const iframe_container = document.getElementById(
    'app-content-main-iframe-wrapper'
  );
  const iframe = document.createElement('iframe');
  iframe.src = 'data:text/html;base64,' + dictQueryStore.mainContent;
  iframe.frameBorder = '0';
  iframe.width = '100%';
  iframe.height = '100%';
  iframe.id = 'app-content-main-iframe';
  iframe.setAttribute('style', 'border: 0px;');
  iframe_container.appendChild(iframe);
}

function updateIframeContent(content, is_base64 = true) {
  const iframe = document.getElementById(
    'app-content-main-iframe'
  ) as unknown as HTMLIFrameElement;
  if (!iframe) {
    return;
  }
  if (is_base64) {
    let encodedStr = unescape(encodeURIComponent(content));
    let src = 'data:text/html;charset=utf-8;base64,';
    iframe.src = src + btoa(encodedStr);
  } else {
    iframe.src = content;
  }
  setTimeout(() => {
    iframe.contentWindow.postMessage(TOP_WIN_MSG_SETUP, '*');
  }, 1000);
}

function listenInnerFrameMessage() {
  window.onmessage = function (e) {
    console.debug('[TOPWIN GOT INNERFRAME MSG] ', e);
    if (!e || !e.data || !e.data.evtype) {
      return;
    }

    switch (e.data.evtype) {
      // entry:// 跳转
      case INNER_FRAME_MSG_ENTRY_JUMP: {
        console.log('inner frame jump to entry: ', e.data);
        let keyWord = e.data.word;
        keyWord = keyWord.split('#')[0];
        dictQueryStore.updateInputSearchWord(keyWord);
        dictQueryStore.searchWord(keyWord);
        dictQueryStore.pushHistoryByEntryIDx(0);

        break;
      }
    }
  };
}

function todo() {
  message.info('功能开发中');
}

// 缩小
function zoomOut() {
  const evtype = TOP_WIN_MSG_ZOOM_OUT;
  const iframe = document.getElementById(
    'app-content-main-iframe'
  ) as unknown as HTMLIFrameElement;
  if (!iframe) {
    return;
  }
  iframe.contentWindow.postMessage(
    { evtype: evtype, ts: new Date().getTime() },
    '*'
  );
}

function refresh() {
  const iframe = document.getElementById(
    'app-content-main-iframe'
  ) as unknown as HTMLIFrameElement;
  if (!iframe) {
    return;
  }
  const evtype = TOP_WIN_MSG_REFRESH;
  iframe.contentWindow.postMessage(
    { evtype: evtype, ts: new Date().getTime() },
    '*'
  );
}

// 放大
function zoomIn() {
  const evtype = TOP_WIN_MSG_ZOOM_IN;
  const iframe = document.getElementById(
    'app-content-main-iframe'
  ) as unknown as HTMLIFrameElement;
  if (!iframe) {
    return;
  }
  iframe.contentWindow.postMessage(
    { evtype: evtype, ts: new Date().getTime() },
    '*'
  );
}

// devtools
function showInspector() {
  const wailsEvent = "wails:showInspector";
   // @ts-ignore
  if (window.WailsInvoke) {
    // @ts-ignore
    window.WailsInvoke(wailsEvent).then(resp =>{
    })
  }
}


let storeChangeUnscribe = null;
function listenContentUpdate() {
  storeChangeUnscribe = dictQueryStore.$onAction(({name, store, after}) => {
      after((result: any) => {
        switch (name) {
          case 'updateMainContent': {
            const content = b64DecodeUnicode(store.mainContent);
            updateIframeContent(content, true);
            break;
          }
          case 'updateMainContentURL': {
            if (store.mainContentURL === '') {
              const content = b64DecodeUnicode(store.mainContent);
              updateIframeContent(content, true);
            }

            updateIframeContent(store.mainContentURL, false);
            break;
          }
        }
      });
    }
  );
}

onMounted(() => {
  cerateIframe();
  if (storeChangeUnscribe) {
    storeChangeUnscribe();
    storeChangeUnscribe = null;
  }
  listenContentUpdate();
  listenInnerFrameMessage();
  setTimeout(function () {
    dictQueryStore.setUpAPIBaseURL();
  }, 1000);
});

onUnmounted(()=>{
  if(storeChangeUnscribe) {
    storeChangeUnscribe();
    storeChangeUnscribe = null;
  }
})

///----------------------------
// utils function
///----------------------------

function b64DecodeUnicode(str) {
  return decodeURIComponent(
    atob(str)
      .split('')
      .map(function (c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
      })
      .join('')
  );
}
</script>
