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
  display: flex;
  flex-direction: row;
  height: 100%;
  width: $layout-right-toolbar-width;
  background-color: #fafafa;
  .dictionaries {
    display: flex;
    flex-direction: column;
    width:60px;
    margin-top: $layout-main-content-functions-height;

    .dictionary-item {
      display: block;
      width: 48px;
      height: 48px;
      text-align: center;
      line-height: 48px;
      margin: 6px auto;
      border: 1px solid #ccc;
      border-radius: 8px;;
      font-size: 13px;
      cursor: pointer;
      user-select: none;
      -webkit-user-select: none;

      &:hover {
        background-color: #f1f1f1;
      }

    }

  }
}
</style>
<template>
  <div id="app-right-toolbar" class="app-right-toolbar">
      <div class="dictionaries" >
        <span class="dictionary-item" v-for="item in state.dictList" 
        :key="item.id" 
        data-dictid="item.id" @click="chooseDict(item)">{{ item.name }}</span>
    </div>
  </div>
</template>
<script setup>
import { useDictQueryStore } from '@/store/dict';
import {reactive, onMounted} from "vue";
const dictQueryStore = useDictQueryStore();

const state = reactive({
  dictList: [],
});

function chooseDict(item) {
  dictQueryStore.updateSelectDict(item);
  dictQueryStore.updateMainContent("");
  dictQueryStore.updatePendingList([]);
}

function loadDictionaries() {
    dictQueryStore.queryDictList().then((res) => {
      if (res.length > 0) {
        dictQueryStore.updateSelectDict(res[0]);
      }
      
      for (let i = 0; i < res.length; i++) {
        state.dictList.push(res[i]);
      }
    });
}

onMounted(() => {
  loadDictionaries();
});

</script>
