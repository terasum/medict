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

<template>
  <div class="x-space">
    <div class="x-layout">
      <div class="x-layout-main-area">
        <div class="x-layout-sidebar">
          <AppSidebar />
        </div>
        <div class="x-layout-content">
          <AppMainContent />
        </div>
        <div class="x-layout-right-toolbar">
          <AppRightToolbar />
        </div>
      </div>
      <div class="n-layout-footer">
        <AppFooter />
      </div>
    </div>
  </div>
</template>

<script setup>
import AppSidebar from './AppSidebar.vue';
import AppFooter from './AppFooter.vue';
import AppMainContent from './AppMainContent.vue';
import AppRightToolbar from './AppRightToolbar.vue';
</script>

<style lang="scss" scoped>
@import "@/style/variables.scss";

.x-space {
  width: 100%;
  height: 100%;
  padding: 0;
  margin: 0;

  .x-layout {
    width: 100%;
    height: 100%;
    padding: 0;
    margin: 0;
    .x-layout-main-area {
      display: flex;
      flex-direction: row;
      width: 100%;
      height: calc(100% - $layout-footer-height);
      padding: 0;
      margin: 0;

      .x-layout-sidebar {
        width: $layout-left-sidebar-width;
        height: 100%;
        padding: 0;
        margin: 0;
        background-color: #fafafa;
      }

      .x-layout-content {
        width: 100%;
        height: 100%;
        padding: 0;
        margin: 0;
      }
    }
    .n-layout-footer {
      width: 100%;
      height: $layout-footer-height;
      padding: 0;
      margin: 0;
    }
  }
}
</style>
