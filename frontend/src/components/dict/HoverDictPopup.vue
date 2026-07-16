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

<!-- 悬停弹窗(#780):在光标附近浮出该词的完整释义(GoldenDict 式)。
     内容 = SearchWord 首匹配 → buildEntryURL → iframe(复用内容管线)。
     经 <Teleport to="body"> + fixed 高 z-index,盖在主 iframe 之上。
     弹窗内的 entry:// 跳转会以 postMessage 冒泡到父窗(走 ENTRY_JUMP 通路)。 -->
<style lang="scss" scoped>
.hover-popup {
  position: fixed;
  z-index: 1001;
  width: 380px;
  max-height: 320px;
  display: flex;
  flex-direction: column;
  background: #fff;
  border: 1px solid #d0d7de;
  border-radius: 8px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.18);
  overflow: hidden;

  .hover-popup-head {
    flex: 0 0 auto;
    display: flex;
    align-items: center;
    gap: 8px;
    height: 26px;
    padding: 0 10px;
    background: #f6f8fa;
    border-bottom: 1px solid #e3e3e3;
    font-size: 12px;
    .hover-popup-word { font-weight: 600; color: #24292f; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
    .hover-popup-status { color: #999; }
  }
  .hover-popup-body {
    flex: 1 1 auto;
    height: 280px;
    background: #fff;
    .hover-popup-iframe { width: 100%; height: 100%; display: block; }
    .hover-popup-placeholder { height: 100%; display: flex; align-items: center; justify-content: center; color: #bbb; font-size: 13px; }
  }
}
</style>

<template>
  <Teleport to="body">
    <div v-if="visible" class="hover-popup" :style="posStyle" @mouseleave="onLeave">
      <div class="hover-popup-head">
        <span class="hover-popup-word">{{ word }}</span>
        <span v-if="loading" class="hover-popup-status">查询中…</span>
        <span v-else-if="empty" class="hover-popup-status">无此词</span>
      </div>
      <div class="hover-popup-body">
        <div v-if="empty || loading" class="hover-popup-placeholder">
          {{ empty ? '该词典无此词' : '查询中…' }}
        </div>
        <iframe
          v-else-if="url"
          ref="iframeRef"
          class="hover-popup-iframe"
          :src="url"
          frameborder="0"
          style="border: 0;"
          @load="onIframeLoad"
        ></iframe>
      </div>
    </div>
  </Teleport>
</template>

<script lang="ts" setup>
import { computed, ref, watch } from 'vue';
import { SearchWord } from '@/apis/dicts-api';
import { buildEntryURL, useDictQueryStore } from '@/store/dict';

const props = defineProps<{
  visible: boolean;
  word: string;
  dictId: string;
  x: number;
  y: number;
}>();
const emit = defineEmits<{ (e: 'close'): void }>();

const dictQueryStore = useDictQueryStore();
const iframeRef = ref<HTMLIFrameElement | null>(null);
const url = ref('');
const loading = ref(false);
const empty = ref(false);

const SETUP_MSG = '__Medict_TOP_WIN_MSG__EVTY_SETUP__';

// 贴边翻转:预估尺寸,超出视口则翻到光标另一侧。
const posStyle = computed(() => {
  const W = 380, H = 320, margin = 12;
  let left = props.x + margin;
  let top = props.y + margin;
  if (left + W > window.innerWidth) left = Math.max(8, props.x - W - margin);
  if (top + H > window.innerHeight) top = Math.max(8, props.y - H - margin);
  return { left: left + 'px', top: top + 'px' };
});

// 词/词典变化时取释义 URL。visible 关闭时清空。
watch(
  () => [props.visible, props.word, props.dictId] as const,
  async ([vis, word, dictId]) => {
    if (!vis || !word || !dictId) {
      url.value = ''; empty.value = false; loading.value = false;
      return;
    }
    loading.value = true; empty.value = false; url.value = '';
    try {
      const res: any = await SearchWord(dictId, String(word));
      const list = Array.isArray(res) ? res : [];
      if (list.length === 0) {
        empty.value = true;
      } else {
        url.value = buildEntryURL(dictQueryStore.dictApiBaseURL, dictId, list[0], 0);
      }
    } catch (e) {
      empty.value = true;
    }
    loading.value = false;
  },
  { immediate: true }
);

function onIframeLoad() {
  // 弹窗内禁用 hover 取词(避免弹窗叠弹窗 + 省掉无谓的 mousemove 处理)。
  iframeRef.value?.contentWindow?.postMessage(
    { evtype: SETUP_MSG, hoverEnabled: false } as any,
    '*'
  );
}

function onLeave() {
  emit('close');
}

defineExpose({ iframeRef });
</script>
