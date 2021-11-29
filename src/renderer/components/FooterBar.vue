<template>
  <div class="footer">
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
import Vue from 'vue';

import { createByProc } from '@terasum/electron-call';
import { WindowAPI } from '../../main/apis/WindowAPI';
const rendererStub = createByProc('renderer', 'error')
const windowApi = rendererStub.use<WindowAPI>('main','WindowApi')

export default Vue.extend({
  data() {
    return {};
  },
  computed: {},
  watch: {},
  methods: {
    onClickHyperLink(event: any) {
      if (event && event.target) {
        if (event.target.dataset && event.target.dataset.href) {
          windowApi.openExternalURL(event.target.dataset.href);
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
.footer {
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