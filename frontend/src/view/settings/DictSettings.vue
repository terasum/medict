<template>
  <div class="parent-window" style="height: 100%">
    <!-- content goes inside .window-content -->
    <div class="window-content">settings</div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import "../assets/css/photon.min.css"
import { defineComponent } from 'vue';
export default defineComponent({


  components: {},
  methods: {},
});
//  <style scoped lang="css" src="../assets/css/photon.min.css"></style>
</script>

<style lang="scss" scoped>
.parent-window {
  height: 100%;
  .window-content {
    height: 100%;
  }
}
</style>
