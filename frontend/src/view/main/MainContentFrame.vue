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

  .app-content-main {
    width: 100%;
    height: calc(100% - $layout-header-height);
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
      <div id="app-content-main-iframe-wrapper"></div>
    </div>
</template>

<script lang="ts" setup>

import { onMounted } from 'vue';
import { useDictQueryStore } from '@/store/dict';


const dictQueryStore = useDictQueryStore();


function cerateIframe() {
  if (document.getElementById('app-content-main-iframe')) {
    document.getElementById('app-content-main-iframe').remove();
  }
  const iframe_container = document.getElementById(
    'app-content-main-iframe-wrapper'
  );
  const iframe = document.createElement('iframe');
  iframe.src = 'data:text/html;base64,' + dictQueryStore.mainContent;
  iframe.frameBorder = "0";
  iframe.width = '100%';
  iframe.height = '100%';
  iframe.id = 'app-content-main-iframe';
  iframe.setAttribute("style", 'border: 0px;');
  iframe_container.appendChild(iframe);
}

function listenContentUpdate() {
  const unsubscribe = dictQueryStore.$onAction(
    ({
      name, // action 名称
      store, // store 实例，类似 `someStore`
      args, // 传递给 action 的参数数组
      after, // 在 action 返回或解决后的钩子
      onError, // action 抛出或拒绝的钩子
    }) => {
      let startTime = Date.now();

      switch (name) {
        case 'updateMainContent':
          break;
        case 'updateMainContentURL':
          break;
        case 'updateSelectDict':
          break;
        case 'updateInputSearchWord':
          break;
        default: {
          console.log(
            `[event] not recognized event, skipped, event name: ${name}`
          );
          return;
        }
      }

      // 这将在 action 成功并完全运行后触发。
      // 它等待着任何返回的 promise
      after((result) => {
        switch (name) {
          case 'updateMainContent': {
            const content = b64DecodeUnicode(store.mainContent);
            updateIframeContent(content, true);
            break;
          }
          case 'updateMainContentURL': {
            if (store.mainContentURL === "") {
              const content = b64DecodeUnicode(store.mainContent);
              updateIframeContent(content, true);  
            }

            updateIframeContent(store.mainContentURL, false);
            break;
          }
          case 'updateInputSearchWord' :{
            // skip for now
            break
          }
          case 'updateSelectDict': {
            // inputWord.value = '';
            break;
          }
        }
      });

      // 如果 action 抛出或返回一个拒绝的 promise，这将触发
      onError((error) => {
        console.warn(
          `Failed "${name}" after ${
            Date.now() - startTime
          }ms.\nError: ${error}.`
        );
      });
    }
  );

  // 手动删除监听器
  //   unsubscribe()
}

function updateIframeContent(content, is_base64 = true) {
  const iframe = document.getElementById('app-content-main-iframe') as unknown as HTMLIFrameElement;
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
}


onMounted(() => {
  cerateIframe();
  listenContentUpdate();
  setTimeout(function () {
    dictQueryStore.setUpAPIBaseURL();
  }, 1000);
});

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
