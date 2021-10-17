<template>
  <div class="window-content">
    <div class="dictionary-settings">
      <!-- Mini button group -->
      <div class="toolbar toolbar-header">
        <div class="toolbar-actions">
          <div class="btn-group">
            <button class="button toolbar-btn" v-on:click="addDictionary">
              <span class="icon"><i class="fas fa-plus-square" /></span>
              添加词典
            </button>
            <button class="button toolbar-btn" v-on:click="refreshDicts">
              <span class="icon"><i class="fas fa-sync" /></span> 刷新列表
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


    <DictModal 
      v-model="isShowDictModalActive"
    >
      <template v-slot:modal-header>
        <h3><i class="fas fa-plus-square"></i> 词典信息</h3>
      </template>
      <NewDictionary
        :dictData="selectedModalDict"
        :showResourceBtn="true"
        v-on:modal-should-close="closeModal"
        :readOnly="true"
      />
    </DictModal>

    <DictModal 
      v-model="isAddDictModalActive"
    >
      <template v-slot:modal-header>
        <h3><i class="fas fa-plus-square"></i> 新增词典</h3>
      </template>
      <NewDictionary
        :dictData="newDictData"
        :showResourceBtn="false"
        v-on:modal-should-close="closeModal"
        :readOnly="false"
      />
    </DictModal>

  </div>
  <!-- endof preference-->
</template>

<script lang="ts">
import Vue from 'vue';
import Header from '../../components/Header.vue';
import NewDictionary from '../../components/preference/NewDictionary.vue';
import DictModal from '../../components/default/DictModal.vue';
import { SyncMainAPI } from '../../service.renderer.manifest';
import { random_key } from '../../../utils/random_key';
import { StorabeDictionary } from '../../../model/StorableDictionary';
import Store from '../../store/index';
// declare const MAIN_WINDOW_WEBPACK_ENTRY: string;
// declare const DICT_SETTINGS_WINDOW_WEBPACK_ENTRY: string;

export default Vue.extend({
  components: { Header, NewDictionary, DictModal },
  computed: {
    dictionaries() {
      return (this.$store as typeof Store).state.dictionaries;
    },
  },
  data: () => {
    return {
      isAddDictModalActive: false,
      isShowDictModalActive: false,
      selectedModalDict: {},
      newDictData: {
        id: '',
      } as StorabeDictionary,
    };
  },
  methods: {
    closeModal() {
      this.isAddDictModalActive= false;
      this.isShowDictModalActive= false;
    },
    refreshDicts() {
      //this.dictionaries = SyncMainAPI.dictFindAll(undefined);
    },
    addDictionary() {
      this.newDictData.id = random_key(6);
      // this.$bvModal.show('add-dictionary-modal');
      this.isAddDictModalActive = true;
    },
    openDictionary(id: string) {
      // console.log(MAIN_WINDOW_WEBPACK_ENTRY);
      // console.log(DICT_SETTINGS_WINDOW_WEBPACK_ENTRY);
      // this.$bvModal.show('dictionary-item-modal');
      this.isShowDictModalActive = true;
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
  width: 100%;
}
.dictionary-settings {
  display: flex;
  width: 100%;
  overflow: hidden;
  flex-direction: column;
  padding-left: 10px;
  padding-right: 10px;

  .dictionary-table {
    width: 100%;
    overflow-x: auto;
    &::v-deep table {
      border-spacing: 0;
      width: 100%;
      border: 0;
      border-collapse: separate;
      font-size: 12px;
      text-align: left;
      cursor: default;
    }

    &::v-deep thead {
      display: table-header-group;
      vertical-align: middle;
      border-color: inherit;
      tr {
        &:first-child:hover {
          background: #fff;
        }
      }
    }

    &::v-deep td,
    th {
      padding: 2px 15px;
    }

    &::v-deep tr {
      &:hover {
        background: #ddd;
      }
    }
    &::v-deep th {
      border-bottom: 1px solid #aaa;
      text-align: center;
    }
  }
}

.toolbar {
  padding: 0px 10px;
  height: 30px;
  width: 100%;
  display: flex;
  background: #fff;
  flex-direction: row-reverse;

  .toolbar-actions {
    margin-top: 0px;
    margin-bottom: 0px;
  }

  .toolbar-btn {
    background-color: rgb(252, 252, 252);
    background-image: linear-gradient(to bottom, #fcfcfc 0, #f1f1f1 100%);
    border-radius: 4px;
    border: 1px solid transparent;
    box-shadow: 0 1px 1px rgb(0 0 0 / 6%);
    border-color: #c2c0c2 #c2c0c2 #a19fa1;
    text-align: center;
    color: #666;
    font-size: 12px;
    display: inline-block;
    padding: 3px 8px;
    margin-bottom: 0;
    white-space: nowrap;
    &:active {
      background-color: #ddd;
      background-image: none;
    }
  }
}

</style>
