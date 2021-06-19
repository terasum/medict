<template>
  <div class="container-fluid" style="height: 100%">
    <Header :displaySearchBox=true />
    <div class="row" style="height: 100%">
        preference
    </div>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Header from "../components/Header.vue";

export default Vue.extend({
  components: { Header },
});
</script>

<style lang="scss" scoped>
</style>
