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
          <AppSidebar>
            <div class="dict-groups">
              <div class="dict-group-list">
                <nav class="nav-group">
                  <h5 class="nav-group-title">APP 设置</h5>
                  <!-- <a class="nav-group-item active"> -->
                  <a class="nav-group-item">
                    <span class="icon icon-book"></span>
                    <router-link to="/setting/dict"> 词典设置 </router-link>
                  </a>
                  <span class="nav-group-item">
                    <span class="icon icon-cog"></span>
                    <router-link to="/setting/software"> 软件设置 </router-link>
                  </span>
                  <span class="nav-group-item">
                    <span class="icon icon-palette"></span>
                    <router-link to="/setting/theme"> 主题设置 </router-link>
                  </span>
                  <span class="nav-group-item">
                    <span class="icon icon-rocket"></span>
                    <router-link to="/setting/plugin"> 插件设置 </router-link>
                  </span>
                </nav>
                <nav class="nav-group">
                  <h5 class="nav-group-title">关于信息</h5>
                  <span class="nav-group-item">
                    <span class="icon icon-help-circled"></span>
                    <router-link to="/docs"> 使用说明</router-link>
                  </span>
                  <span class="nav-group-item">
                    <span class="icon icon-feather"></span>
                    <router-link to="/setting/terms"> 隐私声明</router-link>
                  </span>
                  <span class="nav-group-item">
                    <span class="icon icon-info-circled"></span>
                    <router-link to="/setting/about"> 关于信息</router-link>
                  </span>
                  <span class="nav-group-item">
                    <span class="icon icon-arrows-ccw"></span>
                    <router-link to="/setting/update"> 版本更新</router-link>
                  </span>
                  <span class="nav-group-item">
                    <span class="icon icon-cc"></span>
                    <router-link to="/setting/license"> 开源协议</router-link>
                  </span>
                </nav>
              </div>
            </div>
          </AppSidebar>
        </div>
        <div class="x-layout-content">
          <AppMainContent>
            <router-view></router-view>
          </AppMainContent>
        </div>
        <div class="x-layout-right-toolbar">
          <AppRightToolbar> <div class="dict-toolbar"></div> </AppRightToolbar>>
        </div>
      </div>
      <div class="n-layout-footer">
        <AppFooter> </AppFooter>
      </div>
    </div>
  </div>
</template>

<script setup>
import AppSidebar from '../../components/layout/AppSidebar.vue';
import AppFooter from '../../components/layout/AppFooter.vue';
import AppFunctions from '../../components/layout/AppFunctions.vue';
import AppMainContent from '../../components/layout/AppMainContent.vue';
import AppRightToolbar from '@/components/layout/AppRightToolbar.vue';
import { RouterView, RouterLink } from 'vue-router';

import { NIcon, NButton, NButtonGroup } from 'naive-ui';
import { useUIStore } from '@/store/ui';
import { useDictQueryStore } from '@/store/dict';
import { GetAllDicts } from '@/apis/dicts-api';

import { ref, onMounted } from 'vue';

const dictsList = ref([]);

const dictQueryStore = useDictQueryStore();
const uiStore = useUIStore();
uiStore.updateCurrentTab('setting');
</script>

<style lang="scss" scoped>
@import '@/style/variables.scss';
@import '@/style/photon/photon.scss';

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
      margin-right: -60px;

      .x-layout-sidebar {
        width: $layout-left-sidebar-width;
        height: 100%;
        padding: 0;
        margin: 0;
        background-color: #fafafa;
      }

      .x-layout-content {
        width: calc(
          100% - $layout-left-sidebar-width - $layout-right-toolbar-width
        );
        height: 100%;
        padding: 0;
        margin: 0;
      }

      .x-layout-right-toolbar {
        width: $layout-right-toolbar-width;
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
.nav-group-item {
  a {
    color: #666;
    text-decoration: none;
  }
}
</style>
<style lang="scss">
.markdown-body {
  height: 100%;
  width: 100%;
  font-size: 100%;
  overflow-y: scroll;
  -webkit-text-size-adjust: 100%;
  -ms-text-size-adjust: 100%;
  padding-bottom: 30px;

  color: #444;
  font-family: Georgia, Palatino, 'Palatino Linotype', Times, 'Times New Roman',
    serif;
  font-size: 14px;
  line-height: 20px;
  padding: 1em;
  margin: auto;
  max-width: 45em;
  a {
    color: #0645ad;
    text-decoration: none;
  }
  a:visited {
    color: #0b0080;
  }
  a:hover {
    color: #06e;
  }
  a:active {
    color: #faa700;
  }
  a:focus {
    outline: thin dotted;
  }
  a:hover,
  a:active {
    outline: 0;
  }

  ::-moz-selection {
    background: rgba(255, 255, 0, 0.3);
    color: #000;
  }
  ::selection {
    background: rgba(255, 255, 0, 0.3);
    color: #000;
  }

  a::-moz-selection {
    background: rgba(255, 255, 0, 0.3);
    color: #0645ad;
  }
  a::selection {
    background: rgba(255, 255, 0, 0.3);
    color: #0645ad;
  }

  p {
    margin: 1em 0;
  }

  img {
    max-width: 100%;
    margin: 1em 0;
  }

  h1,
  h2,
  h3,
  h4,
  h5,
  h6 {
    font-weight: normal;
    color: #343434;
    line-height: 1em;
    margin: 1em 0;
  }
  h4,
  h5,
  h6 {
    font-weight: bold;
  }
  h1 {
    font-size: 1.8em;
    text-align: center;
  }
  h2 {
    font-size: 1.5em;
  }
  h3 {
    font-size: 1.3em;
  }
  h4 {
    font-size: 1.1em;
  }
  h5 {
    font-size: 1em;
  }
  h6 {
    font-size: 1em;
  }

  blockquote {
    color: #666666;
    margin: 0;
    padding-left: 3em;
    border-left: 0.5em #eee solid;
  }
  hr {
    display: block;
    height: 0;
    border: 0;
    border-top: 1px solid #aaa;
    border-bottom: 1px solid #eee;
    margin: 1em 0;
    padding: 0;
  }
  pre,
  code,
  kbd,
  samp {
    color: #000;
    font-family: monospace, monospace;
    _font-family: 'courier new', monospace;
    font-size: 0.98em;
  }
  pre {
    white-space: pre;
    white-space: pre-wrap;
    word-wrap: break-word;
    code {
      width: 100%;
    }
  }

  b,
  strong {
    font-weight: bold;
  }

  dfn {
    font-style: italic;
  }

  ins {
    background: #ff9;
    color: #000;
    text-decoration: none;
  }

  mark {
    background: #ff0;
    color: #000;
    font-style: italic;
    font-weight: bold;
  }

  sub,
  sup {
    font-size: 75%;
    line-height: 0;
    position: relative;
    vertical-align: baseline;
  }
  sup {
    top: -0.5em;
  }
  sub {
    bottom: -0.25em;
  }

  ul,
  ol {
    margin: 1em 0;
    padding: 0 0 0 2em;
  }
  li p:last-child {
    margin: 0;
  }
  dd {
    margin: 0 0 0 2em;
  }

  img {
    border: 0;
    -ms-interpolation-mode: bicubic;
    vertical-align: middle;
  }

  table {
    border-collapse: collapse;
    border-spacing: 0;
  }
  td {
    vertical-align: top;
  }
}
</style>
