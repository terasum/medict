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
              <h5 class="nav-group-title">词典配置</h5>
              <span
                class="nav-group-item"
                :class="currentMenu === 0 ? 'active' : ''"
                @click="onClickPreferenceMenu(0)"
              >
                <span class="icon icon-book"></span>
                词典配置
              </span>
              <span
                class="nav-group-item"
                :class="currentMenu === 1 ? 'active' : ''"
                @click="onClickPreferenceMenu(1)"
              >
                <span class="icon icon-language"></span>
                翻译配置
              </span>
              <h5 class="nav-group-title">系统设置</h5>
              <span
                class="nav-group-item"
                :class="currentMenu === 2 ? 'active' : ''"
                @click="onClickPreferenceMenu(2)"
              >
                <span class="icon icon-cog"></span>
                偏好设置
              </span>
              <span
                class="nav-group-item"
                :class="currentMenu === 3 ? 'active' : ''"
                @click="onClickPreferenceMenu(3)"
              >
                <span class="icon icon-tools"></span>
                开发者工具
              </span>
              <span
                class="nav-group-item"
                :class="currentMenu === 4 ? 'active' : ''"
                @click="onClickPreferenceMenu(4)"
              >
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
    <FooterBar />
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import Header from '../components/Header.vue';
import FooterBar from '../components/FooterBar.vue';
import NewDictionary from '../components/preference/NewDictionary.vue';

// declare const MAIN_WINDOW_WEBPACK_ENTRY: string;
// declare const DICT_SETTINGS_WINDOW_WEBPACK_ENTRY: string;

const routerMap = {
  0: '/preference/dictSettings',
  1: '/preference/translateSettings',
  2: '/preference/settings',
  3: '/preference/debug',
  4: '/preference/about',
};

export default Vue.extend({
  components: { Header, NewDictionary, FooterBar },
  computed: {},
  data: () => {
    return {
      currentMenu: 0,
    };
  },
  methods: {
    onClickPreferenceMenu(id: number) {
      console.log(`click id ${id} router: ${routerMap[id]}`);
      if (this.currentMenu != id) {
        this.currentMenu = id;
        if (routerMap[id] && this.$router.currentRoute !== routerMap[id]) {
          this.$router.replace(routerMap[id]);
        }
      }
    },
  },
  mounted() {},
});
</script>

 <style scoped lang="css" src="../assets/css/photon.min.css"></style>

<style lang="scss" scoped>
.window-content {
  height: 100%;
  overflow: hidden;
}
.pane-body {
  width: 100%;
  height: 100%;
  overflow-y: auto;
}
</style>
