<template>
  <div class="container-fluid" style="height: 100%">
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
                <span class="icon icon-user"></span>
                用户设置
              </span>
              <span class="nav-group-item">
                <span class="icon icon-info-circled"></span>
                关于信息
              </span>
            </nav>
          </div>

          <div class="pane dictionary-settings">
            <!-- Mini button group -->
            <div class="toolbar toolbar-header">
              <div class="toolbar-actions">
                <div class="btn-group">
                  <button class="btn btn-default" v-on:click="addDictionary">
                    <span class="icon icon-plus-squared"></span>
                  </button>
                  <button class="btn btn-default">
                    <span class="icon icon-cw"></span>
                  </button>
                </div>
              </div>
            </div>
            <div class="dictionary-table">
              <table class="table-striped">
                <thead>
                  <tr>
                    <th>词典简称</th>
                    <th>词典文件名</th>
                    <th>词典描述</th>
                    <!-- <th>词典资源文件路径</th>
                    <th>词典增强样式路径</th>
                    <th>词典增强脚本路径</th> -->
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="item in dictionaries"
                    :key="item.id"
                    v-on:click="openDictionary(item.id)"
                  >
                    <td>{{ item.brief }}</td>
                    <td>{{ item.filename }}</td>
                    <td>{{ item.desc }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      <!-- modal start-->
      <div class="dictionary-modal">
        <b-modal
          id="dictionary-item-modal"
          title=""
          hide-footer
          scrollable
          button-size="sm"
          header-class="modal-header"
          footer-class="modal-footer"
        >
          <p>{{ selectedModalDict.brief }}</p>
          <p>{{ selectedModalDict.filename }}</p>
          <p>{{ selectedModalDict.desc }}</p>
          <p>{{ selectedModalDict.mddFilePath }}</p>
          <p>{{ selectedModalDict.mdxFilePath }}</p>
          <p>{{ selectedModalDict.cssFilePath }}</p>
          <p>{{ selectedModalDict.jsFilePath }}</p>
        </b-modal>
      </div>
      <!-- /modal start-->

      <div class="add-dictionary-modal">
        <b-modal
          id="add-dictionary-modal"
          title="新增词典"
          scrollable
          button-size="sm"
          header-class="modal-header"
          footer-class="modal-footer"
        >
          <form>
            <div class="form-group">
              <label>词典id(4英文字符)</label>
              <input type="text" class="form-control" disabled />
              <input
                type="text"
                class="form-control"
                placeholder="用户自定义"
              />
            </div>
            <div class="form-group">
              <label>词典名称</label>
              <input
                type="text"
                class="form-control"
                placeholder="英汉大辞典"
              />
            </div>
            <div class="form-group">
              <label>mdx词典文件</label>
              <input type="file" class="form-control" placeholder="oale.mdx" />
            </div>
            <div class="form-group">
              <label>mdd词典文件</label>
              <input type="file" class="form-control" placeholder="oale.mdd" />
            </div>
            <div class="form-group">
              <label>词典描述</label>
              <textarea class="form-control" rows="2"></textarea>
            </div>
          </form>
        </b-modal>
      </div>

      <!-- endof preference-->
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import Header from '../components/Header.vue';
import Store from '../store/index';

import '../assets/css/photon.min.css';

declare const MAIN_WINDOW_WEBPACK_ENTRY: string;
declare const DICT_SETTINGS_WINDOW_WEBPACK_ENTRY: string;

export default Vue.extend({
  components: { Header },
  computed: {
    dictionaries() {
      return (this.$store as typeof Store).state.dictionaries;
    },
  },
  data: () => {
    return {
      selectedModalDict: {},
    };
  },
  methods: {
    addDictionary() {
      console.log(MAIN_WINDOW_WEBPACK_ENTRY);
      console.log(DICT_SETTINGS_WINDOW_WEBPACK_ENTRY);
      this.$bvModal.show('add-dictionary-modal');
    },
    openDictionary(id: number) {
      console.log(MAIN_WINDOW_WEBPACK_ENTRY);
      console.log(DICT_SETTINGS_WINDOW_WEBPACK_ENTRY);
      this.$bvModal.show('dictionary-item-modal');
      this.selectedModalDict = this.dictionaries[id];

      // show window
      //   apis["createSubWindow"]({
      //     width: 200,
      //     height: 300,
      //     html: DICT_SETTINGS_WINDOW_WEBPACK_ENTRY,
      //     titleBarStyle: "default",
      //     nodeIntegration: true,
      //     contextIsolation: false,
      //   });
    },
  },
});
</script>

<style lang="scss" scoped>
.window-content {
  height: 100%;
}
.toolbar {
  height: 26px;
  .toolbar-actions {
    margin-top: 0px;
    margin-bottom: 0px;
  }
  .btn {
    height: 20px;
    line-height: 20px;
    font-size: 14px;
    padding: 0rem 0.75rem;
  }
}
.dictionary-settings {
  overflow: hidden;
  .dictionary-table {
    overflow-x: scroll;
  }
}
</style>

<style lang="scss">
.modal-header {
  padding: 0.3rem 0.5rem !important;
  button {
    border: 1px solid #ddd;
    outline: none;
    width: 26px;
    height: 26px;
    color: #666;
    text-align: center;
    font-size: 14px;
    padding: 0.1rem;
    border-radius: 5px;
  }
  .modal-title {
    font-size: 16px;
    font-weight: normal;
    padding-left: 10px;
  }
}
.form-control {
  padding: 0.15rem 0.375rem !important;
  font-weight: normal !important;
  &:focus {
    font-size: default !important;
    padding: 0.15rem 0.375rem !important;
    font-weight: default !important;
    outline: 0 !important;
    box-shadow: none !important;
    font-weight: normal !important;
  }
}
.modal-footer {
  padding: 0.3rem 0.5rem !important;
}
</style>
