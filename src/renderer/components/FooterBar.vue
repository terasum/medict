<template>
  <div class="footer">
    <span
      class="hyperlink"
      data-href="https://github.com/terasum/medict/"
      @click="onclickHyperlink"
      >star项目</span
    >

    <span class="split-line"></span>
    <span
      class="hyperlink"
      data-href="https://github.com/terasum/medict/issues"
      @click="onclickHyperlink"
      >问题反馈与建议</span
    >
    <span class="split-line"></span>
    <span
      class="hyperlink"
      data-href="https://github.com/terasum/medict"
      @click="onclickHyperlink"
      >使用说明</span
    >
  </div>
</template>


<script lang="ts">
import Vue from 'vue';
import { AsyncMainAPI } from '../service.renderer.manifest';

export default Vue.extend({
  data() {
    return {};
  },
  computed: {},
  watch: {},
  methods: {
    onclickHyperlink(event: any) {
      if (event && event.target) {
        if (event.target.dataset && event.target.dataset.href) {
          console.log(`clicked ${event.target.dataset.href}`);
          AsyncMainAPI.openUrlOnBrowser(event.target.dataset.href);
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
  border-top: 1px solid #ccc;
  height: 20px;
  overflow: hidden;
  padding: 0;
  margin: 0;

  margin-left: -12px;
  margin-right: -12px;
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