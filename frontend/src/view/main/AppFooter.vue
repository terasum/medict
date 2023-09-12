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