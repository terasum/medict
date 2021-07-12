<template>
  <div class="container-fluid" style="height: 100%; width: 100%">
    <Header :displaySearchBox="true" />
    <div class="row" style="height: 100%">
      <!--preference-->
      <!-- content goes inside .window-content -->
      <div class="window-content">
        <div class="pane-group">
          <div class="pane pane-sm sidebar">
            <nav class="nav-group">
              <h5 class="nav-group-title">词典库配置</h5>
              <span class="nav-group-item active">
                <span class="icon icon-book"></span>
                词典配置
              </span>
              <h5 class="nav-group-title">系统设置</h5>
              <span class="nav-group-item">
                <span class="icon icon-cog"></span>
                偏好设置
              </span>
              <span class="nav-group-item">
                <span class="icon icon-tools"></span>
                开发者工具
              </span>
              <span class="nav-group-item">
                <span class="icon icon-info-circled"></span>
                关于信息
              </span>
            </nav>
          </div>
          <div class="pane-body">
            <router-view></router-view>
          </div>
        </div>
      </div>

      <!-- endof preference-->
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import Header from '../components/Header.vue';
import NewDictionary from '../components/preference/NewDictionary.vue';
import '../assets/css/photon.min.css';
// declare const MAIN_WINDOW_WEBPACK_ENTRY: string;
// declare const DICT_SETTINGS_WINDOW_WEBPACK_ENTRY: string;

export default Vue.extend({
  components: { Header, NewDictionary },
  computed: {},
  data: () => {
    return {};
  },
  methods: {},
  mounted() {},
});
</script>

<style lang="scss" scoped>
.window-content {
  height: 100%;
  overflow: hidden;
}
.pane-body {
  width: 100%;
  overflow: scroll;
}
</style>
