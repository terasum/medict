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
      border-bottom: 1px solid var(--c-gray-100);
      border-radius: 3px;
      user-select: none;
      font-size: 16px;
      -webkit-user-select: none;
      &:hover {
        background-color: var(--c-gray-100);
        cursor: pointer;
      }
    }
    .active {
      background-color: var(--c-gray-100);
    }
  }
}
</style>
<template>
  <AppSidebar>
      <ul id="word-pending-list">
        <li
          v-for="(item, entryIndex) in dictQueryStore.queryPendingList"
          :data-id="entryIndex"
          :key="`${item.keyword}-${entryIndex}`"
          @click="selectItem(entryIndex)"
          :class="selectedIndex === entryIndex ? 'active' : ''"
        >
          <span>{{ item.keyword }}</span>
        </li>
      </ul>
  </AppSidebar>
</template>

<script setup>
import AppSidebar from "@/components/layout/AppSidebar.vue";

import { useDictQueryStore } from '@/store/dict';
import { ref, onMounted, onUnmounted, watch } from 'vue';
const dictQueryStore = useDictQueryStore();
const selectedIndex = ref(0);

watch(
  () => dictQueryStore.queryPendingList,
  () => {
    selectedIndex.value = 0;
  }
);

// 按 entry 定位:多词典模式更新主词典那节,单词典模式走 locateWord。
function locateEntry(entryIndex) {
  if (dictQueryStore.multiMode) {
    dictQueryStore.locateInMultiPrimary(entryIndex);
  } else {
    dictQueryStore.locateWord(entryIndex);
  }
}

function selectItem(entryIndex) {
  selectedIndex.value = entryIndex;
  locateEntry(entryIndex);
}

// 焦点守卫：当用户正在输入框 / 文本域 / contenteditable 中输入时，不劫持方向键
function isTextInputActive() {
  const ae = document.activeElement;
  if (!ae) {
    return false;
  }
  const tag = ae.tagName;
  return tag === 'INPUT' || tag === 'TEXTAREA' || ae.isContentEditable;
}

function onKeyDown(e) {
  if (e.key !== 'ArrowUp' && e.key !== 'ArrowDown') {
    return;
  }
  if (isTextInputActive()) {
    return;
  }
  const len = dictQueryStore.queryPendingList.length;
  if (!len || len <= 0) {
    return;
  }
  if (e.key == 'ArrowUp') {
    if (selectedIndex.value === 0) {
      selectedIndex.value = len - 1;
    } else {
      selectedIndex.value -= 1;
    }
    locateEntry(selectedIndex.value);
  } else if (e.key == 'ArrowDown') {
    if (selectedIndex.value === len - 1) {
      selectedIndex.value = 0;
    } else {
      selectedIndex.value += 1;
    }
    locateEntry(selectedIndex.value);
  }
}

onMounted(() => {
  document.addEventListener('keydown', onKeyDown);
});

onUnmounted(() => {
  document.removeEventListener('keydown', onKeyDown);
});

</script>
