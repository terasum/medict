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

.app-right-toolbar {
  height: 100%;
  width: $layout-right-toolbar-width;

  .toolbar-top {
    height: $layout-header-height;
    width: 100%;

    display: flex;
    padding: 0 10px;
    background-color: $theme-top-header-background-color;
  }

  .toolbar-content {
    display: flex;
    flex-direction: row;
    background-color: #fafafa;
    height: calc(100% - $layout-header-height);

    .dictionaries {
      display: flex;
      flex-direction: column;
      width: 60px;

      .dictionary-item {
        display: block;
        width: 48px;
        height: 48px;
        text-align: center;
        line-height: 48px;
        margin: 6px auto;
        border: 1px solid #ccc;
        border-radius: 8px;
        font-size: 13px;
        cursor: pointer;
        user-select: none;
        -webkit-user-select: none;

        box-shadow: rgba(0, 0, 0, 0.1) 0px 10px 50px;

        &:hover {
          background-color: #f1f1f1;
        }
      }
    }
  }
}
</style>
<template>
  <div id="app-right-toolbar" class="app-right-toolbar">
    <div class="toolbar-top"></div>
    <div class="toolbar-content">
      <div class="dictionaries">


        <n-popover
          v-for="item in state.dictList"
          :overlap="false" placement="left" trigger="hover">
          <template #trigger>
            <span
              class="dictionary-item"
              :key="item.id"
              @click="chooseDict(item)"
              :style="getBackground(item)"
            ></span
            >
          </template>
          <div class="large-text">
            <div>{{ item.description && item.description.title  ? item.description.title : item.name}} </div>
            <div style='max-width: 260px; max-height: 200px; overflow-y: auto'><p v-html='item.description.description'></p> </div>
          </div>
        </n-popover>


      </div>
    </div>
  </div>
</template>
<script setup>
import { useDictQueryStore } from '@/store/dict';
import { reactive, onMounted } from 'vue';
import { BuildIndex } from '@/apis/dicts-api';

import {NPopover} from "naive-ui";


const dictQueryStore = useDictQueryStore();


const state = reactive({
  dictList: [],
});

function chooseDict(item) {
  dictQueryStore.updateSelectDict(item);
  dictQueryStore.updateMainContent('');
  dictQueryStore.updatePendingList([]);
}

function getBackground(item) {
  if (item.background) {
    let style = `background:url(data:image/jpg;base64,${item.background});`
    style += `background-size:cover;`
    style += `background-repeat:no-repeat;`
    style += `background-position:center;`;
    style += `color: #fff;`;
    return style
  }
  return "";
}

function loadDictionaries() {
  dictQueryStore.queryDictList().then((res) => {
    console.log(res);
    if (res.length > 0) {
      dictQueryStore.updateSelectDict(res[0]);
    }

    for (let i = 0; i < res.length; i++) {
      state.dictList.push(res[i]);
    }

    setTimeout(() => {
      // build-index
      BuildIndex().then((resp) => {
        console.log('building index success:', resp);
      });
    }, 1000);
  });
}


onMounted(() => {
  loadDictionaries();
});
</script>
