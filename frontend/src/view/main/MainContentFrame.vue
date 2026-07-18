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
  .app-content-main-toolbar {
    height: 30px;
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    box-shadow: inset 0 calc(max(1px, 0.0625rem) * -1) var(--c-gray-300);
    background-color: var(--c-gray-100);
    .toolbar-dicts{
      width: calc(100% - 30px);
    }

    .toolbar-boxes{
    display: flex;
    .app-content-main-toolbar-box {
      display: block;
      height: 24px;
      width: 24px;
      border: 1px solid var(--c-gray-300);
      font-size: 16px;
      text-align: center;
      line-height: 24px;
      margin-left: 3px;
      margin-right: 3px;
      margin-top: 2px;
      border-radius: 3px;
      background-color: var(--c-gray-100);
      color: var(--c-gray-600);
      svg {
        cursor: pointer;
      }
    }
  }
  }
  #app-content-main-iframe-wrapper {
    height: calc(100% - 80px);
    padding: 8px 4px;
    .app-content-main-iframe {
      width: 100%;
      height: 100%;
    }
    // 多词典模式:各词典释义节纵向堆叠,整体滚动
    .multi-stack {
      height: 100%;
      overflow-y: auto;
    }
  }
}
</style>
<template>
  <div class="app-content-main">
  
    <div class="app-content-main-toolbar">
      <div class="toolbar-dicts">
         <MainDictsToolbar/>
      </div>
      <div class="toolbar-boxes">
      <span class="app-content-main-toolbar-box" title="导出当前词条 HTML(调试)" @click="onExportEntry"
        ><NIcon><Bug16Regular /></NIcon
      ></span>
      <span class="app-content-main-toolbar-box" title="编辑词典 CSS(#783)" @click="onEditCSS"
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
    </div>
    <div id="app-content-main-iframe-wrapper" ref="wrapperRef">
      <div v-if="showMulti" class="multi-stack">
        <MainDictSection
          v-for="(r, i) in dictQueryStore.multiResults"
          :ref="el => { if (el) sectionRefs[i] = el as any; }"
          :key="r.dictId"
          :name="r.dictName"
          :url="r.url"
          :empty="r.empty"
        />
      </div>
      <iframe
        v-else
        ref="iframeRef"
        class="app-content-main-iframe"
        :src="iframeSrc"
        frameborder="0"
        style="border: 0;"
        @load="onIframeLoad"
      ></iframe>
    </div>
    <HoverDictPopup
      ref="popupRef"
      :visible="popupVisible"
      :word="popupWord"
      :dict-id="popupDictId"
      :x="popupX"
      :y="popupY"
      @close="popupVisible = false"
    />
    <n-modal v-model:show="cssEditorVisible" preset="card" title="词典 CSS 编辑器" style="width: 760px;">
      <Codemirror
        :model-value="cssContent"
        @update:model-value="onCSSChange"
        :extensions="[cssLang()]"
        :style="{ height: '380px', fontSize: '13px' }"
        placeholder="/* 在此输入自定义 CSS,实时预览(防抖 300ms 注入 iframe);保存后下次查词自动生效 */"
      />
      <template #footer>
        <div style="display: flex; justify-content: flex-end; gap: 8px;">
          <n-button @click="cssEditorVisible = false">关闭</n-button>
          <n-button type="primary" @click="onSaveCSS">保存</n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { useDictQueryStore } from '@/store/dict';
import { ZoomIn16Regular, ZoomOut16Regular,ArrowClockwise20Filled, Bug16Regular, DocumentCss20Regular } from '@vicons/fluent';
import { NIcon, NModal, NButton } from 'naive-ui';
import { useMessage } from 'naive-ui';
import { Codemirror } from 'vue-codemirror';
import { css as cssLang } from '@codemirror/lang-css';
import { exportCurrentEntry, getDictUserCSS, saveDictUserCSS } from '@/apis/dicts-api';
import { getPreferences } from '@/apis/config';

import MainDictsToolbar from "./MainDictsToolbar.vue";
import MainDictSection from "./MainDictSection.vue";
import HoverDictPopup from "@/components/dict/HoverDictPopup.vue";

const dictQueryStore = useDictQueryStore();
const message = useMessage();

const TOP_WIN_MSG_ZOOM_OUT =  '__Medict_TOP_WIN_MSG_EVTP_ZOOM_OUT';
const TOP_WIN_MSG_ZOOM_IN =  '__Medict_TOP_WIN_MSG_EVTP_ZOOM_IN';
const TOP_WIN_MSG_REFRESH = '__Medict_TOP_WIN_MSG_EVTP_REFRESH';
const TOP_WIN_MSG_SETUP =  '__Medict_TOP_WIN_MSG__EVTY_SETUP__';
const INNER_FRAME_MSG_ENTRY_JUMP = '__Medict_INNER_FRAME_MSG_EVTP_ENTRY_JUMP';
const INNER_FRAME_MSG_DBLCLICK_LOOKUP = '__Medict_INNER_FRAME_MSG_EVTP_DBLCLICK_LOOKUP';
const INNER_FRAME_MSG_HOVER_LOOKUP = '__Medict_INNER_FRAME_MSG_EVTP_HOVER_LOOKUP';
const INNER_FRAME_MSG_HOVER_LEAVE = '__Medict_INNER_FRAME_MSG_EVTP_HOVER_LEAVE';

// 声明式 iframe：通过 Vue 响应式驱动 src，避免命令式 createElement / 手动设 .src
const iframeRef = ref<HTMLIFrameElement | null>(null);
const wrapperRef = ref<HTMLDivElement | null>(null);

// 多词典模式：堆叠里各词典节的组件实例（暴露 iframeRef），用于广播缩放/刷新消息。
const sectionRefs = ref<any[]>([]);
const showMulti = computed(
  () => dictQueryStore.multiMode && dictQueryStore.multiResults.length > 0
);
// 把消息广播给当前生效的 iframe：多模式发给所有 section，单模式发给唯一 iframe。
function broadcast(evtype: string) {
  const payload = { evtype, ts: new Date().getTime() };
  if (showMulti.value) {
    for (const s of sectionRefs.value) {
      s?.iframeRef?.contentWindow?.postMessage(payload, '*');
    }
  } else {
    iframeRef.value?.contentWindow?.postMessage(payload, '*');
  }
}

// 悬停弹窗(#780):设置(从 preferences 读,默认开/400ms)+ 弹窗状态。
const hoverEnabled = ref(true);
const hoverDelayMs = ref(400);
const popupVisible = ref(false);
const popupWord = ref('');
const popupDictId = ref('');
const popupX = ref(0);
const popupY = ref(0);
const popupRef = ref<any>(null);

// 弹窗用主/激活词典:多词典模式取激活集首项,否则当前选中词典。
function activeDictId(): string {
  if (dictQueryStore.multiMode && dictQueryStore.multiSelectedDicts.length > 0) {
    return dictQueryStore.multiSelectedDicts[0]?.id || '';
  }
  return dictQueryStore.selectDict?.id || '';
}
// 找到消息来源 iframe(单 iframe 或多词典某节)的屏幕位置,用于坐标映射。
function resolveSourceRect(e: MessageEvent): DOMRect | null {
  const single = iframeRef.value;
  if (single && e.source === single.contentWindow) {
    return single.getBoundingClientRect();
  }
  for (const s of sectionRefs.value) {
    const f = s?.iframeRef as HTMLIFrameElement | null;
    if (f && e.source === f.contentWindow) {
      return f.getBoundingClientRect();
    }
  }
  return null;
}

// iframe 的 src：优先使用 mainContentURL（释义查询 URL），否则用 mainContent 构造 data URL。
// 由 Vue 响应式驱动，store 中 mainContent / mainContentURL 变化时自动更新。
const iframeSrc = computed(() => {
  if (dictQueryStore.mainContentURL) {
    return dictQueryStore.mainContentURL;
  }
  // mainContent 为 base64，经 unicode 安全的编解码往返后构造 data URL
  const decoded = b64DecodeUnicode(dictQueryStore.mainContent);
  return 'data:text/html;charset=utf-8;base64,' + btoa(unescape(encodeURIComponent(decoded)));
});

// iframe 加载完成后发送 setup 消息（替代原来的 1s setTimeout），并下发 hover 设置(#780)
function onIframeLoad() {
  const win = iframeRef.value?.contentWindow;
  if (!win) {
    return;
  }
  win.postMessage(
    { evtype: TOP_WIN_MSG_SETUP, hoverEnabled: hoverEnabled.value, hoverDelayMs: hoverDelayMs.value } as any,
    '*'
  );
}

// 监听内嵌 iframe 的消息（entry:// 跳转等）。
// 使用具名函数 + addEventListener，便于在 onUnmounted 中精确移除，避免 window.onmessage 覆盖与泄漏。
function onInnerFrameMessage(e: MessageEvent) {
  console.debug('[TOPWIN GOT INNERFRAME MSG] ', e);
  if (!e || !e.data || !e.data.evtype) {
    return;
  }
  // 弹窗自身 iframe 发来的消息:entry:///双击照常处理(并关弹窗),hover 一律忽略(防递归)
  const popupIframe = popupRef.value?.iframeRef as HTMLIFrameElement | null;
  const fromPopup = !!popupIframe && e.source === popupIframe.contentWindow;

  switch (e.data.evtype) {
    // entry:// 跳转
    case INNER_FRAME_MSG_ENTRY_JUMP: {
      console.log('inner frame jump to entry: ', e.data);
      let keyWord = e.data.word;
      keyWord = keyWord.split('#')[0];
      dictQueryStore.updateInputSearchWord(keyWord);
      dictQueryStore.searchWord(keyWord);
      dictQueryStore.pushHistoryByEntryIDx(0);
      if (fromPopup) popupVisible.value = false;
      break;
    }
    // 双击选词查词（#258）
    case INNER_FRAME_MSG_DBLCLICK_LOOKUP: {
      let keyWord = (e.data.word || '').split('#')[0];
      if (keyWord) {
        dictQueryStore.updateInputSearchWord(keyWord);
        dictQueryStore.searchWord(keyWord);
        dictQueryStore.pushHistoryByEntryIDx(0);
      }
      if (fromPopup) popupVisible.value = false;
      break;
    }
    // 悬停取词(#780):坐标映射后弹窗
    case INNER_FRAME_MSG_HOVER_LOOKUP: {
      if (fromPopup || !hoverEnabled.value) break;
      const word = String(e.data.word || '').trim();
      const dictId = activeDictId();
      if (!word || !dictId) break;
      const r = resolveSourceRect(e);
      const baseX = r ? r.left : 0;
      const baseY = r ? r.top : 0;
      popupWord.value = word;
      popupDictId.value = dictId;
      popupX.value = baseX + Number(e.data.clientX || 0);
      popupY.value = baseY + Number(e.data.clientY || 0);
      popupVisible.value = true;
      break;
    }
    // 悬停离开(移出/滚动/ESC):关弹窗
    case INNER_FRAME_MSG_HOVER_LEAVE: {
      if (fromPopup) break;
      popupVisible.value = false;
      break;
    }
  }
}

// 导出当前词条为自包含 HTML(资源内联)到本地文件,用于调试复杂词条(#784)。
// 用当前选中词典 + 输入词;多词典模式下导出主词典的释义。
async function onExportEntry() {
  const dict = dictQueryStore.selectDict;
  const word = dictQueryStore.inputSearchWord;
  if (!dict?.id || !word?.trim()) {
    message.warning('请先查一个词');
    return;
  }
  try {
    const path = await exportCurrentEntry(dict.id, word.trim());
    if (path) {
      message.success(`已导出：${path}`);
    } // 用户取消保存对话框时不提示
  } catch (e) {
    message.error((e as Error)?.message || '导出失败');
  }
}

function todo() {
  message.info('功能开发中');
}

// ===== 词典 CSS 编辑器(#783)=====
const TOP_WIN_MSG_APPLY_USER_CSS = '__Medict_TOP_WIN_MSG_EVTP_APPLY_USER_CSS';
const cssEditorVisible = ref(false);
const cssContent = ref('');
let cssPreviewTimer: ReturnType<typeof setTimeout> | null = null;

async function onEditCSS() {
  const dict = dictQueryStore.selectDict;
  if (!dict?.id) {
    message.warning('请先选择一个词典');
    return;
  }
  cssEditorVisible.value = true;
  try {
    cssContent.value = await getDictUserCSS(dict.id);
  } catch {
    cssContent.value = '';
  }
}

// CodeMirror 内容变化 → 防抖 300ms → 客户端注入 iframe(实时预览,无 IPC 往返)
function onCSSChange(value: string) {
  cssContent.value = value;
  if (cssPreviewTimer) clearTimeout(cssPreviewTimer);
  cssPreviewTimer = setTimeout(() => applyUserCSS(value), 300);
}

function applyUserCSS(css: string) {
  const payload = { evtype: TOP_WIN_MSG_APPLY_USER_CSS, ts: Date.now(), css };
  if (showMulti.value) {
    for (const s of sectionRefs.value) {
      s?.iframeRef?.contentWindow?.postMessage(payload, '*');
    }
  } else {
    iframeRef.value?.contentWindow?.postMessage(payload, '*');
  }
}

async function onSaveCSS() {
  const dict = dictQueryStore.selectDict;
  if (!dict?.id) return;
  try {
    await saveDictUserCSS(dict.id, cssContent.value);
    message.success('CSS 已保存(下次查词自动生效)');
  } catch (e) {
    message.error((e as Error)?.message || '保存失败');
  }
}

// 缩小
function zoomOut() {
  broadcast(TOP_WIN_MSG_ZOOM_OUT);
}

function refresh() {
  broadcast(TOP_WIN_MSG_REFRESH);
}

// 放大
function zoomIn() {
  broadcast(TOP_WIN_MSG_ZOOM_IN);
}

// Ctrl/Cmd + =/- 与 Ctrl/Cmd + 滚轮缩放词典 iframe（#260）
function onZoomKey(e: KeyboardEvent) {
  if (!(e.ctrlKey || e.metaKey)) return;
  if (e.key === '=' || e.key === '+') { e.preventDefault(); zoomIn(); }
  else if (e.key === '-' || e.key === '_') { e.preventDefault(); zoomOut(); }
}
function onIframeWheel(e: WheelEvent) {
  if (!(e.ctrlKey || e.metaKey)) return;
  e.preventDefault();
  if (e.deltaY < 0) zoomIn(); else zoomOut();
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

onMounted(() => {
  window.addEventListener('message', onInnerFrameMessage);
  window.addEventListener('keydown', onZoomKey);
  wrapperRef.value?.addEventListener('wheel', onIframeWheel, { passive: false });
  setTimeout(function () {
    dictQueryStore.setUpAPIBaseURL();
  }, 1000);
  // 读取悬停弹窗偏好(var(--c-gray-600)):hoverpopup 默认 true,hoverdelayms 默认 400
  getPreferences().then((p: any) => {
    hoverEnabled.value = p?.hoverpopup !== false;
    const d = Number(p?.hoverdelayms);
    hoverDelayMs.value = d > 0 ? d : 400;
  }).catch(() => { /* 用默认值 */ });
});

onUnmounted(() => {
  window.removeEventListener('message', onInnerFrameMessage);
  window.removeEventListener('keydown', onZoomKey);
  wrapperRef.value?.removeEventListener('wheel', onIframeWheel);
})

///----------------------------
// utils function
///----------------------------

function b64DecodeUnicode(str: string) {
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
