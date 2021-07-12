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
                  <button class="btn btn-default" v-on:click="refreshDicts">
                    <span class="icon icon-cw"></span>
                  </button>
                </div>
              </div>
            </div>
            <div class="dictionary-table">
              <table class="table-striped">
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>词典别名</th>
                    <th>词典全名</th>
                    <th>mdx</th>
                    <th>mdd</th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="item in dictionaries"
                    :key="item.id"
                    v-on:click="openDictionary(item.id)"
                  >
                    <td>{{ item.id }}</td>
                    <td>{{ item.alias }}</td>
                    <td>{{ item.name }}</td>
                    <td>{{ item.mdxpath }}</td>
                    <td>{{ item.mddpath }}</td>
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
          title="词典详情"
          hide-footer
          scrollable
          button-size="sm"
          header-class="modal-header"
          footer-class="modal-footer"
        >
          <NewDictionary
            :dictData="selectedModalDict"
            :showResourceBtn="true"
            :readOnly="true"
          />
        </b-modal>
      </div>
      <!-- /modal start-->

      <div class="add-dictionary-modal">
        <b-modal
          id="add-dictionary-modal"
          title="新增词典"
          hide-footer
          scrollable
          button-size="sm"
          header-class="modal-header"
          footer-class="modal-footer"
        >
          <NewDictionary
            :dictData="newDictData"
            :showResourceBtn="false"
            :readOnly="false"
          />
        </b-modal>
      </div>

      <!-- endof preference-->
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import Header from '../components/Header.vue';
import NewDictionary from '../components/preference/NewDictionary.vue';
import { SyncMainAPI } from '../service.renderer.manifest';
import '../assets/css/photon.min.css';
import { random_key } from '../../utils/random_key';
import { StorabeDictionary } from '../../model/StorableDictionary';
import Store from '../store/index';
// declare const MAIN_WINDOW_WEBPACK_ENTRY: string;
// declare const DICT_SETTINGS_WINDOW_WEBPACK_ENTRY: string;

export default Vue.extend({
  components: { Header, NewDictionary },
  computed: {
    dictionaries() {
      return (this.$store as typeof Store).state.dictionaries;
    },
  },
  data: () => {
    return {
      selectedModalDict: {},
      newDictData: {
        id: '',
      } as StorabeDictionary,
    };
  },
  methods: {
    refreshDicts() {
      //this.dictionaries = SyncMainAPI.dictFindAll(undefined);
    },
    addDictionary() {
      this.newDictData.id = random_key(6);
      this.$bvModal.show('add-dictionary-modal');
    },
    openDictionary(id: string) {
      // console.log(MAIN_WINDOW_WEBPACK_ENTRY);
      // console.log(DICT_SETTINGS_WINDOW_WEBPACK_ENTRY);
      this.$bvModal.show('dictionary-item-modal');
      this.selectedModalDict = SyncMainAPI.dictFindOne({ dictid: id });

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
  mounted() {
    // this.$root.$on('bv::modal::hide', (arg: any) => {
    //   console.log(arg);
    //   this.refreshDicts();
    // });
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
    border: none;
    outline: none;
    width: 22px;
    height: 22px;
    color: #666;
    text-align: center;
    font-size: 18px;
    line-height: 18px;
    padding: 0.1rem;
    border-radius: 5px;
    &:hover {
      color: #999;
      border: 1px solid #ddd;
    }
  }
  .modal-title {
    font-size: 16px;
    font-weight: normal;
    padding-left: 10px;
  }
}

.form-control {
  font-size: 13px !important;
  padding: 0.15rem 0.375rem !important;
  font-weight: normal !important;
  &:focus {
    font-size: 13px !important;
    padding: 0.15rem 0.375rem !important;
    font-weight: default !important;
    outline: 0 !important;
    box-shadow: none !important;
    font-weight: normal !important;
  }
  input {
    font-size: 13px !important;
  }
}
.modal-footer {
  padding: 0.3rem 0.5rem !important;
}
</style>
