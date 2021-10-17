<template>
  <div class="container settings-view">
    <div class="settings-title">
      <h3><i class="fas fa-book"></i> 词典设置</h3>
    </div>
    <div class="settings-body">
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
            <div class="dictionary-list">
              <ul>
                  <li
                    v-for="item in dictionaries"
                    :key="item.id"
                    v-on:click="openDictionary(item.id)"
                    class="dictionary-list-item"
                  >
                   <div class="dictionary-list-item-icon-container">
                     <div class="dictionary-list-item-icon"><i class="fas fa-atlas"></i></div>
                   </div>
                   <div class="dictionary-list-item-infos">
                    <p>ID: {{ item.id }}</p>
                    <p>词典别名: {{ item.alias }}</p>
                    <p>词典全名: {{ item.name }}</p>
                    <!-- <p>mdx: {{ item.mdxpath }}</p> -->
                    <!-- <p>mdd: {{ item.mddpath }}</p> -->
                     </div>
                    
                </li>

              </ul>
            </div>
          </div>
    </div>

    <DictModal v-model="isShowDictModalActive">
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

    <DictModal v-model="isAddDictModalActive">
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
      this.isAddDictModalActive = false;
      this.isShowDictModalActive = false;
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
.settings-title {
  padding: 10px 0 10px 12px;
  border-bottom: 1px solid #dfdfdf;
  h3 {
    font-size: 20px;
    font-weight: 700px;
    color: #777;
  }
}
.settings-section {
  display: flex;
  flex-direction: column;
  border-bottom: 1px solid #c1c1c3;
  margin: 0 10px;
  padding: 12px 0;

  .section-title {
    width: 100%;
    padding: 2px 4px;
    font-size: 16px;
  }

  .section-body {
    padding-left: 16px;
  }

  .input-group {
    width: 100%;
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    .input-info {
      label {
        font-size: 15px;
        color: rgb(43, 103, 194);
      }
      p {
        font-size: 12px;
        color: #999;
      }
    }

    .input-container {
      display: flex;
      flex-direction: column;
      justify-content: center;
      max-width: 15rem;
      padding-right: 20px;
      input {
        height: 26px;
      }
    }
  }
}
.settings-body {
  min-height: 260px;
}
.settings-footer {
  display: flex;
  flex-direction: row-reverse;
  padding-right: 30px;
  margin-top: 10px;
}
.btn {
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
  height: 26px;
  min-width: 80px;
  &:active {
    background-color: #ddd;
    background-image: none;
  }
}

.dictionary-settings {
  display: flex;
  width: 100%;
  overflow: hidden;
  flex-direction: column;
  padding-left: 10px;
  padding-right: 10px;
  margin-top: 10px;

  .dictionary-list{
    padding: 8px 10px;
    max-height: 420px;
    overflow-y: auto;
    overflow-x: hidden;
    .dictionary-list-item {
      border: 1px solid #c1c1c1;
      display: flex;
      flex-direction: row;
      justify-content: space-between;
      padding: 4px 2px;
      border-radius: 4px;
      box-shadow: 1px 1px 3px #c1c1c1;
      cursor: pointer;
      margin-bottom: 5px;
      &:hover{
        background: #f5f5f5;
      }
      &:active{
        background: #f5f5f5;
      }

      .dictionary-list-item-icon-container {
        text-align: center;
        display: flex;
        flex-direction: column;
        justify-content: center;
.dictionary-list-item-icon{
        flex-direction: row;
        display: flex;
        justify-content: center;
        width: 64px;
        font-size: 32px;
        line-height: 64px;
        color: #666;
}
      }

      .dictionary-list-item-infos {
        width: calc(100% - 80px);
        font-size: 12px;

      }
    }
  }

  .dictionary-table {
    width: 100%;
    overflow-x: auto; //allows the div to scroll horizontally only when the div overflows.
    overflow-y: hidden; //does not allow for the div to scroll vertically
    &::v-deep table {
      border-spacing: 0;
      width: 100%;
      border: 0;
      border-collapse: separate;
      font-size: 12px;
      text-align: left;
      cursor: default;
      min-width:100%

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
      td {
        min-width: 50px;
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
  flex-direction: row;

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
    height: 26px;
    &:active {
      background-color: #ddd;
      background-image: none;
    }
  }
}
</style>
