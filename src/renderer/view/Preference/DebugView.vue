<template>
  <div class="container-fluid debug-view">
    <div class="debug-section">
      <h3>开发者工具</h3>
      <p>可以通过以下按钮打开 chrome-dev-tool 面板(用于调试渲染进程)</p>
      <b-button
        class="btn-sm"
        variant="outline-secondary"
        @click="onClickDevBtn"
        >dev-tool</b-button
      >
    </div>

    <span class="split-line" />

    <div class="debug-section">
      <h3>运行日志</h3>
      <p>运行日志路径为:</p>
      <b-button
        class="btn-sm"
        variant="outline-secondary"
        @click="onClickMainProcessLog"
        >打开主进程日志</b-button
      >
    </div>

    <div class="debug-section">
      <h3>词典资源</h3>
      <p>词典资源路径为:</p>
      <b-button class="btn-sm" variant="outline-success" @click="onClickRescDir"
        >打开资源文件夹</b-button
      >
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import { AsyncMainAPI } from '../../service.renderer.manifest';
import BugFill from '../../components/icons/bug-fill.icon.vue';
export default Vue.extend({
  components: { BugFill },
  methods: {
    onClickDevBtn() {
      AsyncMainAPI.openDevTool();
    },
    onClickMainProcessLog() {
      AsyncMainAPI.openMainProcessLog();
    },
    onClickRescDir() {
      AsyncMainAPI.openResourceDir();
    },
  },
});
</script>

<style lang="scss" scoped>
.window-content {
  display: block;
  height: 100%;
  margin-top: 12px;
}
.debug-view {
}
.debug-section {
  margin-top: 20px;
  margin-bottom: 20px;
}
</style>
