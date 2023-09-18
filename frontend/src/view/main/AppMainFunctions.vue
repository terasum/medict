<style lang="scss">
@import '@/style/variables.scss';

.app-content-functions {
  height: $layout-header-height;
  display: flex;
  width: 100%;
  flex-direction: row;
  background: #fafafa;

  .header {
    height: 60px;
    display: flex;
    width: 100%;
    flex-direction: row;
    justify-content: space-between;
    background-color: $theme-top-header-background-color;

    .header-search-box {
      display: flex;
      height: 54px;
      padding: 0;
      margin: 0;

      .header-navigate-btns {
        height: 54px;
        max-width: 120px;
        padding: 0;
        margin-left: 16px;
        margin-right: 14px;
        display: flex;

        .btn-nav {
          height: 26px;
          width: 26px;
          margin-top: 14px;
          padding: 0;
          text-align: center;
          font-size: 12px;
          color: #333;
          outline: none;
          border: 1px solid #fefefe;
          background-color: #fff;
          box-shadow: none;

          &:active {
            box-shadow: none;
            border: rgba(63, 80, 236, 0.452) 1px solid;
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
      .header-search-input {
        height: 54px;
        display: flex;
        flex-direction: column;
        justify-content: center;

        .n-input {
          height: 26px;
          padding: 0 8px;
          margin: 0;
          box-shadow: none;
          font-size: 15px;
          // border: 1px solid #f1f1f1;
          border: none;

          background-color: #fff;
          padding-left: 5px;
          &:active {
            outline: none;
          }
          &:focus {
            outline: none;
          }
        }
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
        background-color: $theme-function-box-active-color;
      }
      .fn-box {
        width: 52px;
        height: 44px;
        margin-top: -6px;
        padding-top: 4px;
        cursor: pointer;
        border-radius: 4px;
        margin-left: 2px;

        color: $theme-function-box-font-color;

        &:hover {
          background-color: $theme-function-box-hover-bg-color;
          color: $theme-function-box-hover-font-color;
        }

        .fn-box-icon {
          width: 20px;
          height: 20px;
          display: block;
          font-size: 16px;
          line-height: 18px;
          margin-left: auto;
          margin-right: auto;
        }
        .fn-box-text {
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
}
</style>
<template>
  <div class="app-content-functions">
    <div class="header">
      <div class="header-search-box">
        <div class="header-navigate-btns">
          <button
            type="button"
            class="button btn btn-light btn-nav btn-nav-left"
          >
            <n-icon><AngleLeft /></n-icon>
          </button>

          <button
            type="button"
            class="button btn btn-light btn-nav btn-nav-right"
          >
            <n-icon><AngleRight /></n-icon>
          </button>
        </div>
        <div class="header-search-input">
          <n-input
            type="text"
            size="small"
            placeholder="搜索"
            @change="handleChange"
            v-model:value="inputWord"
          >
            <template #suffix>
              <n-icon :component="Search" />
            </template>
          </n-input>
        </div>
      </div>

      <div class="header-functions">
        <!-- <div
          class="fn-box"
          v-bind:class="{ 'fn-box-active': currentTab === '词典' }"
          v-on:click="clickDictionary"
        > -->
        <div class="fn-box">
          <span class="fn-box-icon">
            <Search />
          </span>
          <span class="fn-box-text">搜索</span>
        </div>
        <div
          class="fn-box"
          v-bind:class="{ 'fn-box-active': currentTab === '词典' }"
          v-on:click="clickDictionary"
        >
          <span class="fn-box-icon">
            <Book />
          </span>
          <span class="fn-box-text">词典</span>
        </div>

        <div class="fn-box">
          <span class="fn-box-icon">
            <ToggleOn />
          </span>
          <span class="fn-box-text">设置</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { Search, AngleLeft, AngleRight, Book, ToggleOn } from '@vicons/fa';
import { Settings48Filled } from '@vicons/fluent';
import { NIcon } from 'naive-ui';

import { SearchWord } from '@/apis/dicts-api';
import { useDictQueryStore } from '@/store/dict';
const dictQueryStore = useDictQueryStore();

let inputWord = ref('');

function searchWord(word) {
  if (dictQueryStore.selectDict.id === '') {
    console.log('skipped');
    return;
  }
  if (word === '') {
    console.log('empty word skipped');
  }

  SearchWord(dictQueryStore.selectDict.id, word).then((res) => {
    console.log('=== SearchWord ===');
    console.log(dictQueryStore.selectDict.id);
    console.log(res);
    dictQueryStore.updatePendingList(res);
    if (res.data && res.data.length > 0) {
      dictQueryStore.locateWord(0);
    }
  });
}

///----------------------------
// event listener function
///----------------------------

function handleChange(v) {
  console.info('[Event change]: ' + v);
  searchWord(v.trim());
}
</script>

<!-- <script lang="ts">
import Vue from 'vue';
import Store from '../store/index';
import Translate from '../components/icons/translate.icon.vue';
import Search from '../components/icons/search.icon.vue';
import Plugins from '../components/icons/plugins2.icon.vue';
import Settings from '../components/icons/settings.icon.vue';

interface DictItem {
  id: string;
  alias: string;
  name: string;
}
export default Vue.extend({
  components: {
    Translate,
    Search,
    Plugins,
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
      let dicts = (this.$store as typeof Store).state.dictionaries;
      return !dicts ? 0 : dicts.length;
    },
    currentDict() {
      return (this.$store as typeof Store).state.currentSelectDict;
    },
    selectDicts() {
      let dicts: DictItem[] = [];
      (this.$store as typeof Store).state.dictionaries.forEach((item) => {
        dicts.push({ id: item.id, alias: item.alias, name: item.name });
      });
      return dicts;
    },
    currentTab() {
      return (this.$store as typeof Store).state.headerData.currentTab;
    },
  },
  watch: {
    searchWord(word) {
      console.debug(`asyncSearchWord ${word}`);
      this.$store.dispatch('asyncSearchWord', {
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
        this.$store.dispatch('asyncSearchWord', {
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
      this.$store.dispatch('asyncFindWordPrecisly', idx);
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
    clickTranslation(event: any) {
      console.log(event);
      this.$store.commit('updateTab', '翻译');

      if (this.$router.currentRoute.path !== '/translate') {
        this.$router.replace({ path: '/translate' });
      }
    },
    clickPlugins(event: any) {
      this.$store.commit('updateTab', '插件');

      if (this.$router.currentRoute.path !== '/plugins') {
        this.$router.replace({ path: '/plugins' });
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
}); -->
<!-- </script> -->
