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

.app-sidebar {
  display: flex;
  flex-direction: column;
  background-color: #fafafa;
  height: 100%;
  width: 100%;
  overflow: hidden;
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  cursor:default;

  
  .sidebar-content{
    border-right: 1px solid #dcdcdc;
    padding: 0;
    margin: 0;
    user-select: none;
    height: calc(100% -  $layout-sidebar-logo-height);
    ul{
      height: 100%;
      overflow-y: scroll;
      list-style: none;
      margin: 0;
      padding: 0 0 0 6px;;
      li{
       margin: 0 4px 0 4px;
       padding: 0 0px 0 6px;;
       border-bottom: 1px solid #f1f1f1;
       border-radius: 3px;
       user-select: none;
       font-size: 16px;
       -webkit-user-select: none;
       &:hover{
          background-color: #f1f1f1;
          cursor: pointer;
       }
      }
      .active{
          background-color: #f2f2f2;
       }
    }

  }
}
</style>
<template>
  <div id="app-sidebar" class="app-sidebar">
    <div class="sidebar-content">
        <slot></slot>
    </div>
  </div>
</template>

<script setup>
</script>
