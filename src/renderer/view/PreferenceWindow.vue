<template>
  <div class="container-fluid" style="height: 100%; width: 100%">
    <Header :displaySearchBox="true" />
    <div class="row" style="height: 100%">
      <!--preference-->
      <!-- content goes inside .window-content -->
      <div class="window-content">
        <div class="pane-group">
          <div class="pane pane-sm sidebar">
            <div class="nav-group">
              <h5 class="nav-group-title">词典配置</h5>
              <span
                class="nav-group-item"
                :class="currentMenu === 0 ? 'item-active' : ''"
                @click="onClickPreferenceMenu(0)"
              >
                <span class="icon"><i class="fas fa-book"></i></span>
                词典配置
              </span>
              <span
                class="nav-group-item"
                :class="currentMenu === 1 ? 'item-active' : ''"
                @click="onClickPreferenceMenu(1)"
              >
                <span class="icon"><i class="fas fa-language"></i></span>
                翻译配置
              </span>
              <h5 class="nav-group-title">系统设置</h5>
              <span
                class="nav-group-item"
                :class="currentMenu === 2 ? 'item-active' : ''"
                @click="onClickPreferenceMenu(2)"
              >
                <span class="icon"><i class="fas fa-cog"></i></span>
                偏好设置
              </span>
              <span
                class="nav-group-item"
                :class="currentMenu === 3 ? 'item-active' : ''"
                @click="onClickPreferenceMenu(3)"
              >
                <span class="icon"><i class="fas fa-bug"></i></span>
                开发设置
              </span>
              <span
                class="nav-group-item"
                :class="currentMenu === 4 ? 'item-active' : ''"
                @click="onClickPreferenceMenu(4)"
              >
                <span class="icon"><i class="fas fa-info-circle"></i></span>
                关于信息
              </span>
            </div>
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


<style lang="scss" scoped>
.window-content {
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: row;
  .pane-group {
    width: 100%;
    display: flex;
  }
  .sidebar {
    width: 160px;
    height: 100%;
    border-right: 1px solid #e8e8e8;
    background: #f2f4f5;
    .nav-group {
      user-select: none;
      display: flex;
      flex-direction: column;
      padding: 10px 6px;

      .nav-group-title {
        display: none;
        width: 100%;
        margin: 0;
        padding-left:6px;
        font-size: 12px;
        font-weight: 500;
        color: #999;
        border-bottom: 1px solid #c1c1c1;
        margin-bottom: 5px;
      }

      .nav-group-item {
        height: 32px;
        width: 100%;
        display: flex;
        color: #777;
        text-decoration: none;
        font-size: 14px;
        line-height: 32px;
        cursor: pointer;

        &:active {
          background-color: #e8eaec;
        }

        .icon {
          height: 30px;
          width: 30px;
          display: inline-block;
          font-size: 14px;
          line-height: 30px;
          text-align: center;
          color: #777;
        }
      }
      .item-active {
        border-radius: 3px;
        background-color: #e8eaec;
        color: #222;
      }
    }
  }
}
.pane-body {
  width: calc(100% - 161px);
  height: 100%;
  overflow-y: auto;
}
</style>
