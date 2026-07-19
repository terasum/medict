<template>
  <div class="settings-docs">
    <MarkdownNavigation :items="docsNavigation" />
    <main class="settings-docs-reader">
      <article class="markdown-body">
        <RouterView />
      </article>
    </main>
  </div>
</template>

<script setup lang="ts">
import { RouterView } from 'vue-router';
import MarkdownNavigation from '@/components/setting/MarkdownNavigation.vue';
import { docsNavigation } from './docs-navigation';
</script>

<style lang="scss" scoped>
.settings-docs {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  min-width: 0;
  overflow: hidden;
  background: #fff;
}

.settings-docs-reader {
  flex: 1 1 auto;
  min-width: 0;
  min-height: 0;
  overflow-y: auto;
  background: #fff;
}

.markdown-body {
  width: min(760px, calc(100% - 56px));
  margin: 0 auto;
  padding: 28px 28px 48px;
  box-sizing: border-box;
  color: var(--c-gray-800);
  font-family: Lato, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  font-size: 14px;
  line-height: 1.7;
  text-wrap: pretty;
}

.markdown-body :deep(h1),
.markdown-body :deep(h2),
.markdown-body :deep(h3),
.markdown-body :deep(h4) {
  color: var(--c-gray-900);
  line-height: 1.3;
}

.markdown-body :deep(h1) {
  margin: 0 0 20px;
  font-size: 24px;
  font-weight: 700;
}

.markdown-body :deep(h2) {
  margin: 28px 0 12px;
  padding-bottom: 7px;
  border-bottom: 1px solid var(--c-gray-200);
  font-size: 19px;
  font-weight: 650;
}

.markdown-body :deep(h3) {
  margin: 22px 0 8px;
  font-size: 16px;
  font-weight: 650;
}

.markdown-body :deep(p),
.markdown-body :deep(ul),
.markdown-body :deep(ol),
.markdown-body :deep(table),
.markdown-body :deep(blockquote) {
  margin: 0 0 14px;
}

.markdown-body :deep(a) {
  color: var(--c-primary);
  text-decoration: none;
}

.markdown-body :deep(a:hover) {
  text-decoration: underline;
}

.markdown-body :deep(a:focus-visible) {
  outline: 2px solid var(--c-primary);
  outline-offset: 2px;
  border-radius: 2px;
}

.markdown-body :deep(img) {
  display: block;
  max-width: 100%;
  height: auto;
  margin: 18px auto;
  border: 1px solid var(--c-gray-200);
  border-radius: 6px;
}

.markdown-body :deep(table) {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}

.markdown-body :deep(th),
.markdown-body :deep(td) {
  padding: 8px 10px;
  border: 1px solid var(--c-gray-300);
  text-align: left;
  vertical-align: top;
}

.markdown-body :deep(th) {
  background: var(--c-gray-100);
  font-weight: 650;
}

.markdown-body :deep(code) {
  padding: 1px 4px;
  background: var(--c-gray-100);
  border-radius: 3px;
  font-family: "Fira Code", ui-monospace, monospace;
  font-size: 0.92em;
}

.markdown-body :deep(blockquote) {
  padding: 8px 12px;
  color: var(--c-gray-700);
  background: var(--c-gray-50);
  border-left: 3px solid var(--c-gray-300);
}

@media (max-width: 820px) {
  .markdown-body {
    width: calc(100% - 32px);
    padding: 22px 16px 40px;
  }
}
</style>
