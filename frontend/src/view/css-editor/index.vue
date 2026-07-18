<template>
  <main class="css-editor-window">
    <header class="css-editor-header">
      <div>
        <h1>词典 CSS 编辑器</h1>
        <p>{{ dictionary.name || dictionary.id }}</p>
      </div>
      <span class="save-state" :class="saveState">{{ saveStateText }}</span>
    </header>
    <section class="editor-shell">
      <Codemirror
        v-model="cssContent"
        :disabled="!loaded"
        :extensions="extensions"
        :style="{ height: '100%', fontSize: '13px' }"
        placeholder="/* 在此输入自定义 CSS；修改会自动保存并实时应用到主窗口 */"
        @change="scheduleSave"
      />
    </section>
    <footer class="css-editor-footer">
      <span>⌘/Ctrl + S 保存</span>
      <div class="actions">
        <n-button @click="closeWindow">关闭</n-button>
        <n-button type="primary" :disabled="!loaded" :loading="saveState === 'saving'" @click="save()">保存</n-button>
      </div>
    </footer>
  </main>
</template>

<script lang="ts" setup>
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { NButton, useMessage } from 'naive-ui';
import { Codemirror } from 'vue-codemirror';
import { css } from '@codemirror/lang-css';
import { EventsOn } from '../../../wailsjs/runtime/runtime';
import { getDictUserCSS, saveDictUserCSS } from '@/apis/dicts-api';

type SaveState = 'clean' | 'dirty' | 'saving' | 'error';
const dictionary = ref({ id: '', name: '' });
const cssContent = ref('');
const saveState = ref<SaveState>('clean');
const extensions = [css()];
const message = useMessage();
const loaded = ref(false);
let saveTimer: ReturnType<typeof setTimeout> | null = null;
let savePromise: Promise<boolean> | null = null;
let stopCloseRequestedListener: (() => void) | null = null;

const saveStateText = computed(() => ({
  clean: '已保存', dirty: '等待保存', saving: '正在保存…', error: '保存失败',
}[saveState.value]));

function scheduleSave() {
  if (!loaded.value) return;
  saveState.value = 'dirty';
  if (saveTimer) clearTimeout(saveTimer);
  saveTimer = setTimeout(() => void save(false), 300);
}

async function save(showConfirmation = true): Promise<boolean> {
  if (!loaded.value || !dictionary.value.id) return false;
  if (savePromise) return savePromise;
  if (saveTimer) {
    clearTimeout(saveTimer);
    saveTimer = null;
  }
  savePromise = (async () => {
    while (true) {
      const content = cssContent.value;
      saveState.value = 'saving';
      try {
        await saveDictUserCSS(dictionary.value.id, content);
      } catch (error) {
        saveState.value = 'error';
        message.error((error as Error)?.message || '保存失败');
        return false;
      }
      if (content === cssContent.value) {
        saveState.value = 'clean';
        if (showConfirmation) message.success('CSS 已保存');
        return true;
      }
    }
  })();
  try {
    return await savePromise;
  } finally {
    savePromise = null;
  }
}

function onKeydown(event: KeyboardEvent) {
  if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === 's') {
    event.preventDefault();
    void save();
  }
}

async function closeWindow() {
  if (saveTimer) {
    clearTimeout(saveTimer);
    saveTimer = null;
  }
  // A failed initial load must never write the empty editor buffer back over
  // an existing sidecar. It is still always safe to close without saving.
  if (!loaded.value) {
    await (window as any).go.main.App.CloseCSSWindow();
    return;
  }
  const saved = saveState.value === 'clean' ? true : await save(false);
  if (!saved) return;
  await (window as any).go.main.App.CloseCSSWindow();
}

onMounted(async () => {
  window.addEventListener('keydown', onKeydown);
  stopCloseRequestedListener = EventsOn('medict:css-editor-close-requested', () => {
    void closeWindow();
  });
  try {
    const context = await (window as any).go.main.App.CSSWindowDictionary();
    dictionary.value = context;
    cssContent.value = await getDictUserCSS(context.id);
    loaded.value = true;
  } catch (error) {
    saveState.value = 'error';
    message.error((error as Error)?.message || '无法加载词典 CSS');
  }
});

onUnmounted(() => {
  window.removeEventListener('keydown', onKeydown);
  stopCloseRequestedListener?.();
  if (saveTimer) clearTimeout(saveTimer);
});
</script>

<style lang="scss" scoped>
@use '@/style/variables.scss' as *;

.css-editor-window {
  box-sizing: border-box;
  display: grid;
  grid-template-rows: auto minmax(0, 1fr) auto;
  gap: 12px;
  width: 100%;
  height: 100%;
  padding: 18px 20px 16px;
  background: var(--c-gray-50, #f8f9fa);
}
.css-editor-header,
.css-editor-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.css-editor-header h1 { margin: 0; font-size: 18px; color: var(--c-gray-900); }
.css-editor-header p { margin: 4px 0 0; color: var(--c-gray-600); }
.save-state { color: var(--c-gray-600); font-size: 12px; }
.save-state.error { color: #d03050; }
.editor-shell {
  min-height: 0;
  overflow: hidden;
  border: 1px solid var(--c-gray-300);
  border-radius: 6px;
  background: white;
}
.editor-shell :deep(.cm-editor) { height: 100%; }
.css-editor-footer { color: var(--c-gray-600); font-size: 12px; }
.actions { display: flex; gap: 8px; }
</style>
