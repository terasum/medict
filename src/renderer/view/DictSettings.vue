<template>
  <div class="parent-window" style="height: 100%">
    <!-- content goes inside .window-content -->
    <div class="window-content">settings</div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';

export default Vue.extend({
  components: {},
  methods: {},
});
</script>

 <style scoped lang="css" src="../assets/css/photon.min.css"></style>
<style lang="scss" scoped>
.parent-window {
  height: 100%;
  .window-content {
    height: 100%;
  }
}
</style>
