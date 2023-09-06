<template>
  <div>
    <b-modal
      custom-class="dict-select-modal"
      :width="480"
      scroll="keep"
      v-model="isDictSelectModalActive"
    >
      <div class="dictionaries-container">
        <ul class="dictionaries">
          <li
            v-for="item in selectDicts"
            :key="item.id"
            @click="selectDictItem(item)"
            :class="
              item.id === currentDict.id
                ? 'dictionary-item dict-active'
                : 'dictionary-item'
            "
          >
            <span class="dict-icon"><i class="fas fa-book"></i></span>
            <div class="dict-info">
              <span class="dict-name">{{ item.name }}</span>
              <span class="dict-desc"> {{ item.id }}</span>
            </div>
          </li>
        </ul>
      </div>
    </b-modal>

    <div class="header">
      <div class="header-navigate-btns">
        <button type="button" class="button btn btn-light btn-nav btn-nav-left">
          <i class="fas fa-chevron-left"></i>
        </button>

        <button
          type="button"
          class="button btn btn-light btn-nav btn-nav-right"
        >
          <i class="fas fa-chevron-right"></i>
        </button>
      </div>
      <div class="header-search-box">
        <b-button
          icon-pack="fas"
          :label="currentDict != undefined ? currentDict.name : 'NULL'"
          icon-left="book"
          @click="onSelectDictBtnClick"
        />
        <b-input
          placeholder="search words..."
          :disabled="displaySearchBox"
          v-model="searchWord"
          @keyup.enter.native="confirmSelect"
          @keyup.up.native="upSelect"
          @keyup.down.native="downSelect"
          type="search"
        ></b-input>
        <b-button icon-pack="fas" icon-right="search" />
      </div>

      <div class="header-functions">
        <div
          class="fn-box"
          v-bind:class="{ 'fn-box-active': currentTab === '词典' }"
          v-on:click="clickDictionary"
        >
          <span class="fn-box-icon">
            <Search :width="16" :height="16" />
          </span>
          <span class="fn-box-text">词典</span>
        </div>

        <div
          class="fn-box"
          v-bind:class="{ 'fn-box-active': currentTab === '设置' }"
          v-on:click="clickPreference"
        >
          <span class="fn-box-icon">
            <Settings :width="16" :height="16" />
          </span>
          <span class="fn-box-text">设置</span>
        </div>
      </div>
    </div>
  </div>
</template>


<script lang="ts">
import Store from '../store/index';
import Search from '../components/icons/search.icon.vue';
import Settings from '../components/icons/settings.icon.vue';

interface DictItem {
  id: string;
  alias: string;
  name: string;
}
import { defineComponent } from 'vue';
export default defineComponent({

  components: {
    Search,
    Settings,
  },
  props: {
    displaySearchBox: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      searchWord: '',
      isDictSelectModalActive: false,
    };
  },

  computed: {
    avaliableDictCount() {
      // let dicts = (this.$store as typeof Store).state.dictionaries;
      // return !dicts ? 0 : dicts.length;
      return 0;
    },
    currentDict() {
      // return (this.$store as typeof Store).state.currentSelectDict;
      return "dict01";
    },
    selectDicts() {
      // let dicts: DictItem[] = [];
      // (this.$store as typeof Store).state.dictionaries.forEach((item) => {
      //   dicts.push({ id: item.id, alias: item.alias, name: item.name });
      // });
      // return dicts;
      return [];
    },
    currentTab() {
      // return (this.$store as typeof Store).state.headerData.currentTab;
      return 0;
    },
  },
  watch: {
    searchWord(word) {
      console.debug(`asyncSuggest ${word}`);
      this.$store.dispatch('asyncSuggest', {
        dictid: this.currentDict.id,
        word,
      });
    },
  },
  methods: {
    onSelectDictBtnClick() {
      if (this.avaliableDictCount > 0) {
        this.isDictSelectModalActive = true;
      }
    },
    selectDictItem(item: DictItem) {
      this.$store.commit('updateCurrentSelectDict', item);
      // async search

      if (this.searchWord && this.searchWord.length > 0) {
        console.log(this.searchWord);
        this.$store.dispatch('asyncSuggest', {
          dictid: item.id,
          word: this.searchWord,
        });
      }

      this.isDictSelectModalActive = false;
      this.$emit('close');
    },
    // keyup event methods
    confirmSelect(event: any) {
      const idx = (this.$store as typeof Store).state.sideBarData
        .selectedWordIdx;
      this.$store.dispatch('asyncLocateWord', idx);
    },
    upSelect(event: any) {
      const idx = (this.$store as typeof Store).state.sideBarData
        .selectedWordIdx;
      this.$store.commit('updateSelectedWordIdx', idx - 1);
    },
    downSelect(event: any) {
      const idx = (this.$store as typeof Store).state.sideBarData
        .selectedWordIdx;
      this.$store.commit('updateSelectedWordIdx', idx + 1);
    },
    clickDictionary(event: any) {
      console.log(event);
      this.$store.commit('updateTab', '词典');

      if (this.$router.currentRoute.path !== '/') {
        this.$router.replace({ path: '/' });
        this.$store.commit('updateSuggestWords', []);
      }
    },

    clickPreference(event: any) {
      this.$store.commit('updateTab', '设置');
      if (this.$router.currentRoute.path !== '/preference') {
        this.$router.replace({ path: '/preference' });
      }
    },
  },
  mounted() {},
});
</script>


<style lang="scss" scoped>
.dict-select-modal {
  background-color: #f3f3f3d3;
  position: fixed;
  width: 100%;
  height: 100%;
  top: 60px;
  z-index: 999;

  &:deep(.modal-background) {
    display: none;
  }

  &:deep(.modal-content) {
    background: #fff;
    border: 1px solid #cfcfcf;
    margin: 0 auto;
    padding: 5px 2px;
    border-radius: 3px;
    box-shadow: 1px 1px 3px #c1c1c1;
    overflow: hidden;
    max-height: calc(100% - 102px);
  }

  &:v-deep(.modal-close) {
    display: none;
    &::before {
      content: '确认';
    }
  }

  &:v-deep(.dictionaries-container) {
    height: 260px;
    width: 100%;
    overflow: hidden;
  }

  &:v-deep(.dictionaries){
    padding: 0;
    margin-bottom: 20px;
    min-height: 100px;
    height: 100%;
    overflow-y: auto;
    /* Works on Chrome, Edge, and Safari */
    &::-webkit-scrollbar {
      width: 5px;
    }

    &::-webkit-scrollbar-track {
      background: transparent;
    }

    &::-webkit-scrollbar-thumb {
      border-radius: 20px;
      background: #f1f1f160;
      border: 1px solid #cfcfcf60;
    }
  }

  &:v-deep(.dict-active) {
    background: #ddeef580;
    &:hover {
      background: #ddeef580 !important;
    }
  }

  &:v-deep(.dictionary-item) {
    display: flex;
    height: 46px;
    padding: 4px 10px;
    border-radius: 3px;
    &:hover {
      background: #f1f1f1;
      cursor: pointer;
    }

    .dict-icon {
      width: 38px;
      height: 38px;
      line-height: 38px;
      font-size: 26px;
      text-align: center;
    }
    .dict-info {
      display: flex;
      flex-direction: column;
      height: 45px;
    }
    .dict-name {
      height: 20px;
      font-size: 16px;
      line-height: 20px;
    }
    .dict-desc {
      height: 18px;
      font-size: 13px;
      line-height: 18px;
      color: #999;
    }
  }
}
.header {
  height: 60px;
  display: flex;
  // background-image: linear-gradient(135deg, #325dff, #529eff);
  // background-color: #eeebea;
  background-color: transparent;
  padding-top: 6px;
  border-bottom: .5px solid #d9d9d9;

  -webkit-user-select: none;
  user-select: none;

  .header-navigate-btns {
    height: 54px;
    max-width: 80px;
    padding: 0;
    margin-left: 34px;
    margin-right: 14px;

    .btn-nav {
      height: 24px;
      width: 26px;
      margin-top: 18px;
      padding: 0;
      text-align: center;
      font-size: 12px;
      color: #333;
      outline: none;
      border: 1px solid #fefefe;
      background-color: #fff;
      box-shadow: none;
      -webkit-user-select: none;
      user-select: none;

      &:active {
        box-shadow: none;
        // border: rgba(63, 80, 236, 0.452) 1px solid;
        border: #c1c1c1 1px solid;
        // background-color: #eeebea;
        // background-color: #d80034;
        background-color: #fff;
      }
    }
    .btn-nav-left {
      border-radius: 10px 0px 0px 10px;
      margin-left: 9px;
    }
    .btn-nav-right {
      border-radius: 0px 10px 10px 0px;
      margin-left: 0px;
    }
  }
  .header-search-box {
    display: flex;
    height: 54px;
    padding: 0;
    margin: 0;
    margin-top: 18px;
    &:v-deep(.btn-group,
      .btn-group-vertical)
       {
        vertical-align: top;
      }
      // toggle button
      button:nth-child(1) {
        border: none;
        border-right: 1px solid #f1f1f1;
        background-color: #fff;
        border-radius: 20px 0px 0px 20px;
        height: 26px;
        font-size: 12px;
        line-height: 26px;
        padding: 0;
        margin: 0;
        padding-left: 10px;
        padding-right: 10px;
        box-shadow: none;
        width: 100px;
        display: flex;
        justify-content: space-between;

        span:nth-child(1) {
          padding-left: 3px;
          padding-right: 5px;
          width: 26px;
        }

        span:nth-child(2) {
          padding-left: 2px;
          padding-right: 1px;
          width: 58px;
          height: 100%;
          text-overflow: ellipsis;
          overflow: hidden;
        }
      }
      // search button
      button:nth-child(3) {
        // background-color: #fff;
        // border: 1px solid #fff;
        border: none;
        border-right: 1px solid #f1f1f1;

        background-color: #fff;
        border-radius: 0px 20px 20px 0px;
        height: 26px;
        padding: 0;
        margin: 0;
        padding-left: 10px;
        padding-right: 10px;
        box-shadow: none;
        line-height: 26px;
        font-size: 12px;

        color: #999;
        cursor: pointer;
      }
      input {
        height: 26px;
        padding: 0 8px;
        margin: 0;
        box-shadow: none;
        font-size: 12px;
        // border: 1px solid #f1f1f1;
        border: none;
        outline: none;

        background-color: #fff;
        padding-left: 5px;
        &:active {
          outline: none;
        }
        &:focus {
          outline: none;
        }
      }
      .form-control:disabled,
      .form-control[readonly] {
        background-color: #fff;
      }
  }

  .header-functions {
    max-width: 306px;
    height: auto;
    padding: 0;
    margin: 0;
    margin-left: 20px;
    margin-top: 12px;
    display: flex;
    flex-direction: row;
    .fn-box-active {
      // background-color: #bd3134;
      background-color: #68686891;
      // background-color: #325effa6;
      .fn-box-icon {
        color: #fff !important;
      }
      .fn-box-text {
        color: #fff !important;
      }
    }
    .fn-box {
      width: 52px;
      height: 44px;
      margin-top: -6px;
      padding-top: 4px;
      cursor: pointer;
      border-radius: 4px;
      margin-left: 2px;

      &:hover {
        // background-color: #c73639;
        // background-color: #325eff60;
        background-color: #a0a0a069;
      }

      .fn-box-icon {
        width: 20px;
        height: 20px;
        // color: #f9dad9;
        color: #696969;
        display: block;
        font-size: 16px;
        line-height: 18px;
        margin-left: auto;
        margin-right: auto;
      }
      .fn-box-text {
        // color: #f9dad9;
        color: #696969;
        margin-left: auto;
        margin-right: auto;
        text-align: center;
        width: 100%;
        display: block;
        font-size: 12px;
        user-select: none;
      }
    }
  }
}
</style>