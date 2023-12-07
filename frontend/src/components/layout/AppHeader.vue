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
  <div class="app-header" id="app-header">
    <div class="app-header-logo">
      <div class="logo">
        <h1>Medict</h1>
      </div>
    </div>
    <div class="app-header-func-box">
    <div class="app-header-func-solt">
        <slot></slot>
    </div>

    <div class="app-header-func-nav">
    <AppFunctions />
    </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import AppFunctions from './AppFunctions.vue';
</script>

<style lang="scss" scoped>
@import '@/style/variables.scss';
@import '@/style/photon/photon.scss';

.app-header {
  width: 100%;
  overflow: hidden;
  padding: 0 20px;
  margin: 0;
  display: flex;
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;

  background-color: $theme-top-header-background-color;
  border-bottom: 1px solid #dcdcdc;

  justify-content: space-between;
  height: $layout-header-height;
  max-height:$layout-header-height;


  .app-header-logo {
    display: flex;
    height: calc($layout-sidebar-logo-height - 1px);
    width: 120px;
    .logo {
      display: flex;
      flex-direction: row;
      align-items: center;
      justify-content: center;
      width: 100%;
      img {
        width: 100%;
      }
      h1 {
        color: $theme-logo-font-color;
        font-style: italic;
        font-size: 24px;
      }
    }
  }
  .app-header-func-box{
    display: flex;
    justify-content: space-between;
    width: calc(100% - 120px);
    .app-header-func-solt{
      display: flex;
      flex-direction: row;
      align-items: center;
      justify-content: center;
      min-width: 120px;
      height: 54px;
      padding: 0;
      margin: 0;
      min-width: 174px;
      margin-left: 40px;
    }
    .app-header-func-nav{
    }
  }
}
</style>
