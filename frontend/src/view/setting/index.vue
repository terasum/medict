<!--
 Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
 GPL-3.0.
-->

<template>
  <div class="x-space settings-shell">
    <div class="x-layout">
      <div class="x-layout-header"><AppHeader /></div>
      <div class="x-layout-main-area">
        <aside class="x-layout-sidebar">
          <AppSidebar>
            <nav class="nav-group settings-nav-menu" aria-label="设置菜单">
              <RouterLink
                v-for="item in settingsNavigation"
                :key="item.to"
                :to="item.to"
                class="nav-group-item settings-nav-item"
              >
                <span :class="['icon', item.icon]" aria-hidden="true" />
                <span class="settings-nav-label">{{ item.label }}</span>
              </RouterLink>
            </nav>
          </AppSidebar>
        </aside>
        <main class="x-layout-content">
          <AppMainContent><RouterView /></AppMainContent>
        </main>
      </div>
      <div class="n-layout-footer"><AppFooter /></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router';
import AppFooter from '@/components/layout/AppFooter.vue';
import AppHeader from '@/components/layout/AppHeader.vue';
import AppMainContent from '@/components/layout/AppMainContent.vue';
import AppSidebar from '@/components/layout/AppSidebar.vue';
import { useUIStore } from '@/store/ui';
import { settingsNavigation } from './settings-navigation';

useUIStore().updateCurrentTab('setting');
</script>

<style lang="scss" scoped>
@use '@/style/variables.scss' as *;

.x-space,
.x-layout {
  width: 100%;
  height: 100%;
  padding: 0;
  margin: 0;
}

.x-layout-main-area {
  display: flex;
  width: 100%;
  height: calc(100% - $layout-footer-height);
}

.x-layout-sidebar {
  width: $layout-left-sidebar-width;
  height: 100%;
  background: var(--c-gray-100);
  border-right: 1px solid var(--c-gray-300);
}

.x-layout-content {
  width: calc(100% - $layout-left-sidebar-width);
  height: 100%;
  min-width: 0;
}

.n-layout-footer {
  width: 100%;
  height: $layout-footer-height;
}

.settings-nav-menu {
  padding: 4px 0;
}

.settings-nav-item {
  display: flex;
  align-items: center;
  min-height: 27px;
  padding: 2px 10px;
  border-radius: 6px;
  margin: 1px 8px;
  white-space: nowrap;
}

.settings-nav-item.router-link-active {
  color: var(--c-gray-900);
  background: var(--c-gray-300);
}

.settings-nav-item .icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex: 0 0 16px;
  float: none;
  width: 16px;
  height: 16px;
  margin: 0 6px 0 0;
  font-size: 16px;
  line-height: 16px;
}

.settings-nav-item .icon::before {
  position: static;
  display: block;
  line-height: 1;
}

.settings-nav-label {
  display: inline-flex;
  align-items: center;
  min-width: 0;
  height: 16px;
  font-size: 13px;
  font-weight: 600;
  line-height: 16px;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
