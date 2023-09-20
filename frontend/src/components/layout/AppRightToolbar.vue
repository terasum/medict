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

.app-right-toolbar {
  height: 100%;
  width: $layout-right-toolbar-width;

  .toolbar-top {
    height: calc($layout-header-height - 1px);
    width: 100%;


    display: flex;
    padding: 0 10px;
    background-color: $theme-top-header-background-color;
    border-bottom: 1px solid #dcdcdc;
  }

  .toolbar-content {
    display: flex;
    flex-direction: row;
    height: calc(100% - $layout-header-height);
    border-left: 1px solid #dcdcdc;
  }
}
</style>
<template>
  <div id="app-right-toolbar" class="app-right-toolbar">
    <div class="toolbar-top"></div>
    <div class="toolbar-content">
      <slot></slot>
    </div>
  </div>
</template>

<script setup></script>
