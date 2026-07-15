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

<!-- 多词典堆叠里的一节:词典名标题 + 单个释义 iframe。供 MainContentFrame 在多词典模式下 v-for。 -->
<style lang="scss" scoped>
.dict-section {
  margin: 0 0 10px 0;
  border: 1px solid #e3e3e3;
  border-radius: 6px;
  overflow: hidden;
  background: #fff;

  .dict-section-head {
    display: flex;
    align-items: center;
    gap: 8px;
    height: 26px;
    padding: 0 10px;
    background: #f6f8fa;
    border-bottom: 1px solid #e3e3e3;
    font-size: 12px;
    color: #555;
    user-select: none;

    .dict-section-name {
      font-weight: 600;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
    .dict-section-tag {
      color: #aaa;
      font-size: 11px;
    }
  }
  .dict-section-body {
    height: 320px; // 每节固定高度,iframe 内部自滚动;整体堆叠由外层滚动
    background: #fff;
    position: relative;

    .dict-section-iframe {
      width: 100%;
      height: 100%;
      display: block;
    }
    .dict-section-empty-body {
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #bbb;
      font-size: 13px;
    }
  }
}
</style>

<template>
  <div class="dict-section">
    <div class="dict-section-head">
      <span class="dict-section-name">{{ name }}</span>
      <span v-if="empty" class="dict-section-tag">无此词</span>
    </div>
    <div class="dict-section-body">
      <div v-if="empty" class="dict-section-empty-body">该词典无此词</div>
      <iframe
        v-else
        ref="iframeRef"
        class="dict-section-iframe"
        :src="url"
        frameborder="0"
        style="border: 0;"
        @load="onLoad"
      ></iframe>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';

// 与 MainContentFrame 的 setup 消息约定一致(内嵌脚本据此挂载 entry:// 跳转/双击选词)。
const SETUP_MSG = '__Medict_TOP_WIN_MSG__EVTY_SETUP__';

const props = defineProps<{
  name: string;
  url: string;
  empty: boolean;
}>();

const iframeRef = ref<HTMLIFrameElement | null>(null);

function onLoad() {
  iframeRef.value?.contentWindow?.postMessage(SETUP_MSG, '*');
}

// 暴露给父组件,用于多词典模式下广播缩放/刷新消息。
defineExpose({ iframeRef });
</script>
