<template>
  <!-- content goes inside .window-content -->
  <div>
    <div class="row justify-content-center">
      <div class="selection-block col-md-3">
        <b-button variant="outline-secondary" @click="onClickDevBtn"
          >开发者工具</b-button
        >
      </div>
      <div class="col-md-3">
        <b-button variant="outline-secondary" @click="onClickMainProcessLog"
          >主进程日志</b-button
        >
      </div>
      <div class="col-md-3">
        <b-button variant="outline-success" @click="onClickRescDir"
          >资源文件夹</b-button
        >
      </div>
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
.selection-block {
  height: 56x;
}
</style>
