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

.app-content {
  display: flex;
  flex-direction: column;
  height: 100%;
  .app-content-functions{
    height: $layout-main-content-functions-height;
    width: 100%;
    display: flex;
    flex-direction: row;
    background: #fafafa;
    .search{
      width:240px;
      display: flex;
      justify-content: center;
      flex-direction: column;
      margin-left: 15px;
    }
  }

  .app-content-main {
    width: 100%;
    height: calc(100% - $layout-main-content-functions-height);
    #app-content-main-iframe-wrapper{
      height: calc(100% - 16px);
      padding: 8px 4px;
      .app-content-main-iframe {
        width: 100%;
        height: 100%;
      }
    }
  }
}
</style>
<template>
  <div class="app-content" id="app-content">
    <div class="app-content-functions">
    <div class="search">
        <!--                 @change="searchWord" -->
        <n-input type="text" size="small" placeholder="搜索" @change="handleChange" v-model:value="inputWord">
          <template #suffix>
            <n-icon :component="Search" />
          </template>
        </n-input>
      </div>

      </div>

    <div class="app-content-main">
      <div id="app-content-main-iframe-wrapper"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useDictQueryStore } from '@/store/dict';

import { Search } from '@vicons/fa';
import { NIcon } from 'naive-ui';

import { SearchWord } from "@/apis/dicts-api"

const dictQueryStore = useDictQueryStore();

let inputWord = ref("");


function cerateIframe() {
  if (document.getElementById('app-content-main-iframe')) {
    document.getElementById('app-content-main-iframe').remove();
  }
  const iframe_container = document.getElementById(
    'app-content-main-iframe-wrapper'
  );
  const iframe = document.createElement('iframe');
  iframe.src = 'data:text/html;base64,' + dictQueryStore.mainContent;
  iframe.frameborder = 0;
  iframe.width = '100%';
  iframe.height = '100%';
  iframe.id = 'app-content-main-iframe';
  iframe.style = 'border: 0px;';
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
      if (name !== "updateMainContent" && name !== "updateMainContentURL" && name !== "updateSelectDict") {
        console.log("not updateMainContent, skipped");
        return;
      }

      // 这将在 action 成功并完全运行后触发。
      // 它等待着任何返回的 promise
      after((result) => {
        if (name === "updateMainContent") {
            const content = b64DecodeUnicode(store.mainContent)
            updateIframeContent(content, true);
        } else if (name === "updateMainContentURL") {
          updateIframeContent(store.mainContentURL, false);
        } else if (name === "updateSelectDict") {
          inputWord.value = ""
        }
      })

      // 如果 action 抛出或返回一个拒绝的 promise，这将触发
      onError((error) => {
        console.warn(
          `Failed "${name}" after ${Date.now() - startTime}ms.\nError: ${error}.`
        )
      })
    }
  )

// 手动删除监听器
//   unsubscribe()

}

function updateIframeContent(content, is_base64=true){
  const iframe = document.getElementById('app-content-main-iframe');
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


function searchWord(word) {
  if (dictQueryStore.selectDict.id === "") {
    console.log("skipped")
    return;
  }
  if (word === "") {
    console.log("empty word skipped")
  }

  SearchWord(dictQueryStore.selectDict.id, word).then((res) =>{
    console.log("=== SearchWord ===")
    console.log(dictQueryStore.selectDict.id)
    console.log(res)
    if (res.code === 200) {
      dictQueryStore.updatePendingList(res.data)
      if (res.data && res.data.length > 0) {
        dictQueryStore.locateWord(0)
      }
    }
  })
}

///----------------------------
// event listener function
///----------------------------

function handleChange (v) {
  console.info('[Event change]: ' + v)
  searchWord(v.trim())
}



onMounted(() => {
  cerateIframe();
  listenContentUpdate();
  dictQueryStore.setUpDictAPI()
});

///----------------------------
// utils function
///----------------------------

function b64DecodeUnicode(str) {
  return decodeURIComponent(atob(str).split('').map(function(c) {
    return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
  }).join(''));
}


</script>
