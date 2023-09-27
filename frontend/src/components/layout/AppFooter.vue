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
  <div class="app-footer" id="app-footer">

    <span class="hyperlink">
        <b data-href="https://github.com/terasum/medict"
        @click="onClickHyperLink">
         <span class="icon icon-github"></span>
        Github</b>
    </span>

    <span class="split-line"></span>
    <span class="hyperlink">
      <b
        data-href="https://github.com/terasum/medict/issues"
        @click="onClickHyperLink"
        >
         <span class="icon icon-help-circled"></span>
        问题反馈</b
      >
    </span>

    <span class="split-line"></span>
    <span class="hyperlink">
      <b>
        <span class="icon icon-gauge"></span>
        性能分析
      </b>
    </span>
    <span class="split-line"></span>
    <span class="hyperlink">
      <b>
        <span class="icon icon-book"></span>
        词典下载
      </b>
    </span>

    <span class="split-line"></span>

    <slot></slot>
  </div>
</template>

<script lang="ts" setup>
import { Bug, Coffee } from '@vicons/fa';
import { NIcon } from 'naive-ui';
import { BrowserOpenURL } from '../../../wailsjs/runtime/runtime';
import { useRouter } from 'vue-router';

const router = useRouter();

function onClickHyperLink(event: any) {
  if (event && event.target) {
    if (event.target.dataset && event.target.dataset.href) {
      BrowserOpenURL(event.target.dataset.href);
    }
  }
  console.log(event);
}
function onClickInternalLink(event: any) {
  if (event && event.target) {
    if (event.target.dataset && event.target.dataset.href) {
      console.log('replace router, path', event.target.dataset.href);
      router.replace({ path: event.target.dataset.href });
    }
  }
  console.log(event);
}
</script>

<style lang="scss" scoped>
@import '@/style/variables.scss';
@import '@/style/photon/photon.scss';

.app-footer {

  position: absolute;
  bottom: 0;
  width: 100%;
  border-top: 1px solid #ccc;
  height: 20px;
  overflow: hidden;
  padding: 0;
  margin: 0;
  background: #fff;
  padding-left: 20px;

  display: flex;

  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  span.hyperlink {
    // padding-left: 5px;
    // padding-right: 5px;
    
    & > b {
      margin: 0;
      padding: 0;
      // height: 20px;
      line-height: 20px;
      display: block;
      font-size: 12px;
      color: #666;
    }

    &:hover {
      color: #333;
    }
  }
  
  span.split-line {
    &::before {
      content: '|';
      color: #ccc;
    }
    margin-left: 3px;
    margin-right: 3px;
  }
  .building-index-process {
    display: flex;
    justify-content: center;
    flex-direction: column;
    height: 20px;
    width: 120px;
  }
}
</style>
