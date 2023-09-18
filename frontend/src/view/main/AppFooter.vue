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
      <NIcon size="12"><Github /></NIcon>
      <b
        data-href="https://github.com/terasum/medict/"
        @click="onClickHyperLink"
        >Github</b
      >
    </span>

    <span class="split-line"></span>
    <span class="hyperlink">
      <NIcon size="12"><Bug /></NIcon>
      <b
        data-href="https://github.com/terasum/medict/issues"
        @click="onClickHyperLink"
        >问题反馈</b
      >
    </span>

    <span class="split-line"></span>
    <span class="hyperlink">
      <NIcon size="12"><Question /></NIcon>
      <b data-href="/docs" @click="onClickInternalLink">使用说明</b>
    </span>

    <span class="split-line"></span>
    <span class="hyperlink">
      <NIcon size="12"><Coffee /></NIcon>
      <b data-href="/about" @click="onClickInternalLink">关于</b>
    </span>
  </div>
</template>

<script lang="ts">
import { Github, Question, Bug, Coffee } from '@vicons/fa';
import { NIcon } from 'naive-ui';
import { BrowserOpenURL } from '../../../wailsjs/runtime/runtime';
import { defineComponent } from 'vue';

export default defineComponent({
  components: { NIcon, Github, Question, Bug, Coffee },
  data() {
    return {};
  },
  computed: {},
  watch: {},
  methods: {
    onClickHyperLink(event: any) {
      if (event && event.target) {
        if (event.target.dataset && event.target.dataset.href) {
          BrowserOpenURL(event.target.dataset.href);
        }
      }
      console.log(event);
    },
    onClickInternalLink(event: any) {
      if (event && event.target) {
        if (event.target.dataset && event.target.dataset.href) {
          console.log('replace router, path', event.target.dataset.href);
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

  display: flex;

  span.hyperlink {
    padding-left: 5px;
    padding-right: 5px;
    cursor: pointer;
    & > i {
      height: 20px;
      width: 20px;
    }
    & > b {
      margin-top: 0;
      padding-top: 0;
      height: 20px;
      width: 20px;
      line-height: 20px;

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
  }
}
</style>
