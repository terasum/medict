<template>
  <div class="container-fluid about-view">
    <div class="px-4 py-5 my-5 text-center">
      <h1 class="display-5 fw-bold">Medict</h1>
      <div class="col-lg-6 mx-auto">
        <p class="lead mb-4">现代跨平台词典App</p>
        <div class="d-grid gap-2 d-sm-flex justify-content-sm-center">
          <button
            type="button"
            class="btn btn-sm btn-outline-primary px-4 gap-3"
            @click="onClickGithub"
          >
            访问 github
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import { AsyncMainAPI } from '../../service.renderer.manifest';

export default Vue.extend({
  components: {},
  methods: {
    onClickGithub() {
      AsyncMainAPI.openUrlOnBrowser('https://github.com/terasum/medict');
    },
  },
});
</script>

<style lang="scss" scoped>
</style>
