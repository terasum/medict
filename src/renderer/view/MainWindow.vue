<template>
  <div class="container-fluid" style="height: 100%">
    <Header :displaySearchBox=false />
    <div class="row" style="height: 100%">
      <div class="col col-2 left-sidebbar-container">
        <div class="left-sidebbar">
          <ul class="left-sidebar-wordlist">
          <li>hello</li>  
          </ul>
        </div>
      </div>
      <div class="col">2 of 2</div>
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
.left-sidebbar-container{
  background-color: #F0E8E9;
  .left-sidebbar {
    height: 100%;
    height: -webkit-calc(100% - 60px);
    height: -moz-calc(100% - 60px);
    height: calc(100% - 60px);
    .left-sidebar-wordlist{
      list-style: none;
      margin: 0;
      padding: 0;
      li{
        padding: 0.2rem 0.2rem 0.2rem 0.2rem;
        border-bottom: 1px solid #333;
      }
    }
  }
}

</style>
