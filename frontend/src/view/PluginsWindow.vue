<template>
  <div class="container-fluid" style="height: 100%">
    <Header :displaySearchBox="true" />
    <div class="fake-plugin" style="height: 100%"><span>Plugins</span></div>
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
.fake-plugin{
  background:#fafafa;
  padding:0;
  margin:0 -12px;
  
  span{
    color: #fafafa;  /*文字颜色必须和背景同色*/
    font-size: 64px;
    margin: auto;
    display: block;
    text-align: center;
    padding-top: 140px;
    text-shadow: 1px 1px 1px rgb(109, 107, 107),-1px -1px 1px #fff;
    user-select: none;
  }
}
</style>
