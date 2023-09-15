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
      <span
        class="hyperlink"
        data-href="https://github.com/terasum/medict/"
        @click="onClickHyperLink"
        >star项目</span
      >
  
      <span class="split-line"></span>
      <span
        class="hyperlink"
        data-href="https://github.com/terasum/medict/issues"
        @click="onClickHyperLink"
        >问题反馈与建议</span
      >
      <span class="split-line"></span>
      <span class="hyperlink" data-href="/docs" @click="onClickInternalLink"
        >使用说明</span
      >
    </div>
  </template>
  
  
  <script lang="ts">
  import { defineComponent } from 'vue';
  export default defineComponent({
    data() {
      return {};
    },
    computed: {},
    watch: {},
    methods: {
      onClickHyperLink(event: any) {
        if (event && event.target) {
          if (event.target.dataset && event.target.dataset.href) {
            // TODO FIX
            // windowApi.openExternalURL(event.target.dataset.href);
          }
        }
        console.log(event);
      },
      onClickInternalLink(event: any) {
        if (event && event.target) {
          if (event.target.dataset && event.target.dataset.href) {
            this.$router.replace({ path: event.target.dataset.href });
          }
        }
        console.log(event);
      },
    },
    mounted() {},
  });
  </script>
  
  
  <style lang="scss" scoped>
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
    color: #666;
    font-size: 12px;
    font-weight: 400;
    line-height: 18px;
    span.hyperlink {
      cursor: pointer;
      &:hover {
        color: #333;
      }
    }
    span.split-line {
      &::before {
        content: '|';
        color: #ccc;
      }
    }
  }
  </style>