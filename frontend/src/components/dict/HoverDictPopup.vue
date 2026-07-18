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

<!-- 悬停弹窗(#780):光标附近浮出该词完整释义。试点:样式全用 UnoCSS 原子类
     (gray 标度 + 卡片 utility),不再有 <style scoped> 硬编码颜色。 -->
<template>
  <Teleport to="body">
    <div
      v-if="visible"
      class="fixed z-[1001] flex flex-col w-[380px] max-h-[320px] overflow-hidden bg-white border border-gray-200 rounded-lg shadow-md"
      :style="posStyle"
      @mouseleave="onLeave"
    >
      <div class="flex items-center gap-2 h-[26px] px-2.5 bg-gray-50 border-b border-gray-200 text-xs">
        <span class="font-semibold text-gray-900 truncate">{{ word }}</span>
        <span v-if="loading" class="text-gray-500">查询中…</span>
        <span v-else-if="empty" class="text-gray-400">无此词</span>
      </div>
      <div class="h-[280px] bg-white">
        <div
          v-if="empty || loading"
          class="h-full flex items-center justify-center text-gray-400 text-[13px]"
        >
          {{ empty ? '该词典无此词' : '查询中…' }}
        </div>
        <iframe
          v-else-if="url"
          ref="iframeRef"
          class="block w-full h-full"
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

// 贴边翻转:超出视口则翻到光标另一侧。
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
