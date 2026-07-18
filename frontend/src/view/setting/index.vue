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
            <nav v-for="group in settingsNavigation" :key="group.title" class="nav-group settings-nav-group">
              <h2 class="nav-group-title">{{ group.title }}</h2>
              <RouterLink
                v-for="item in group.items"
                :key="item.to"
                :to="item.to"
                class="nav-group-item settings-nav-item"
              >
                <span :class="['icon', item.icon]" aria-hidden="true" />
                <span class="settings-nav-copy">
                  <strong>{{ item.label }}</strong>
                  <small>{{ item.description }}</small>
                </span>
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
  width: 210px;
  height: 100%;
  background: var(--c-gray-100);
  border-right: 1px solid var(--c-gray-300);
}

.x-layout-content {
  width: calc(100% - 210px);
  height: 100%;
  min-width: 0;
}

.n-layout-footer {
  width: 100%;
  height: $layout-footer-height;
}

.settings-nav-group + .settings-nav-group {
  margin-top: 8px;
}

.settings-nav-group .nav-group-title {
  padding: 14px 14px 5px;
  font-size: 11px;
  font-weight: 600;
  letter-spacing: 0.04em;
  text-transform: uppercase;
}

.settings-nav-item {
  display: flex;
  align-items: center;
  min-height: 44px;
  padding: 5px 12px;
  border-radius: 6px;
  margin: 1px 8px;
  white-space: normal;
}

.settings-nav-item.router-link-active {
  color: var(--c-gray-900);
  background: var(--c-gray-300);
}

.settings-nav-item .icon {
  flex: 0 0 22px;
  float: none;
  margin: 0 8px 0 0;
}

.settings-nav-copy {
  display: flex;
  min-width: 0;
  flex-direction: column;
  line-height: 1.25;
}

.settings-nav-copy strong {
  font-size: 13px;
  font-weight: 600;
}

.settings-nav-copy small {
  margin-top: 2px;
  color: var(--c-gray-600);
  font-size: 11px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
