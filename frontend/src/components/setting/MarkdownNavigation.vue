<template>
  <nav class="markdown-navigation" aria-label="使用说明文档导航">
    <div class="markdown-navigation-list">
      <RouterLink
        v-for="item in items"
        :key="item.to"
        :to="item.to"
        class="markdown-navigation-item"
      >
        <span :class="['icon', item.icon]" aria-hidden="true" />
        <span>{{ item.label }}</span>
      </RouterLink>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router';

export interface MarkdownNavigationItem {
  label: string;
  icon: string;
  to: string;
}

defineProps<{ items: MarkdownNavigationItem[] }>();
</script>

<style lang="scss" scoped>
.markdown-navigation {
  flex: 0 0 42px;
  width: 100%;
  height: 42px;
  box-sizing: border-box;
  overflow-x: auto;
  overflow-y: hidden;
  background: var(--c-gray-50);
  border-bottom: 1px solid var(--c-gray-300);
}

.markdown-navigation-list {
  display: flex;
  align-items: stretch;
  width: max-content;
  min-width: 100%;
  height: 100%;
  padding: 0 12px;
  box-sizing: border-box;
}

.markdown-navigation-item {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  min-width: 104px;
  height: 42px;
  padding: 0 14px;
  box-sizing: border-box;
  color: var(--c-gray-700);
  border-bottom: 2px solid transparent;
  font-size: 13px;
  font-weight: 500;
  line-height: 16px;
  text-decoration: none;
  white-space: nowrap;
}

.markdown-navigation-item:hover {
  color: var(--c-gray-900);
  background: var(--c-gray-100);
}

.markdown-navigation-item.router-link-active {
  color: var(--c-gray-900);
  background: #fff;
  border-bottom-color: var(--c-primary);
  font-weight: 600;
}

.markdown-navigation-item:focus-visible {
  outline: 2px solid var(--c-primary);
  outline-offset: -3px;
}

.markdown-navigation-item .icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex: 0 0 16px;
  width: 16px;
  height: 16px;
  margin: 0;
  font-size: 16px;
  line-height: 16px;
}

.markdown-navigation-item .icon::before {
  position: static;
  line-height: 1;
}

@media (max-width: 720px) {
  .markdown-navigation-list {
    padding: 0 6px;
  }

  .markdown-navigation-item {
    min-width: auto;
    padding: 0 10px;
  }
}
</style>
