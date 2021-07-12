<template>
  <div class="container-fluid" style="height: 100%">
    <Header :displaySearchBox="true" />
    <div class="row" style="height: 100%">
      <!-- content goes inside .window-content -->
      <div class="window-content">debug</div>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import Header from '../components/Header.vue';
import '../assets/css/photon.min.css';

export default Vue.extend({
  components: { Header },
  methods: {},
});
</script>

<style lang="scss" scoped>
.parent-window {
  height: 100%;
  .window-content {
    height: 100%;
  }
}
</style>
