<template>
  <div class="container-fluid" style="height: 100%">
    <Header :displaySearchBox="true" />
    <div class="row" style="height: 100%">Plugins</div>
    <FooterBar />
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import Header from '../components/Header.vue';
import FooterBar from '../components/FooterBar.vue';

export default Vue.extend({
  components: { Header, FooterBar },
});
</script>

<style lang="scss" scoped>
</style>
