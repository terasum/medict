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

.sidebar-content {
  padding: 4px 2px;
  margin: 2px 0 0 0;
  user-select: none;
  height: calc(100% - $layout-sidebar-logo-height - $layout-footer-height);
  ul {
    height: 100%;
    overflow-y: scroll;
    list-style: none;
    margin: 0;
    padding: 0 0 0 6px;
    li {
      margin: 0 4px 0 4px;
      padding: 0 0px 0 6px;
      border-bottom: 1px solid #f1f1f1;
      border-radius: 3px;
      user-select: none;
      font-size: 16px;
      -webkit-user-select: none;
      &:hover {
        background-color: #f1f1f1;
        cursor: pointer;
      }
    }
    .active {
      background-color: #f2f2f2;
    }
  }
}
</style>
<template>
  <AppSidebar>
    <div class="sidebar-content">
      <ul>
        <li
          v-for="item in dictQueryStore.queryPendingList"
          :data-id="item.id"
          :key="item.id"
          @click="selectItem(item.id)"
          :class="selected_id == item.id ? 'active' : ''"
        >
          <span>{{ item.key_word }}</span>
        </li>
      </ul>
    </div>
  </AppSidebar>
</template>

<script setup>
import AppSidebar from "@/components/layout/AppSidebar.vue";

import { useDictQueryStore } from '@/store/dict';
import { ref } from 'vue';
const dictQueryStore = useDictQueryStore();
const selected_id = ref('0');

function selectItem(entry_id) {
  selected_id.value = entry_id;
  dictQueryStore.locateWord(entry_id);
}
</script>
