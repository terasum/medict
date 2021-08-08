<template>
  <div class="row header">
    <div class="header-navigate-btns">
      <button type="button" class="btn btn-light btn-nav btn-nav-left">
        <b-icon-chevron-compact-left />
      </button>
      <button type="button" class="btn btn-light btn-nav btn-nav-right">
        <b-icon-chevron-compact-right />
      </button>
    </div>
    <div class="header-search-box">
      <div>
        <b-input-group>
          <template v-slot:prepend>
            <b-dropdown :text="currentDict.alias" variant="info">
              <b-dropdown-item
                v-for="item in selectDicts"
                :key="item.id"
                @click="selectDictItem(item)"
                >{{ item.alias }}</b-dropdown-item
              >
            </b-dropdown>
          </template>
          <b-form-input
            :disabled="displaySearchBox"
            v-model="searchWord"
            @keyup.enter.native="confirmSelect"
            @keyup.up.native="upSelect"
            @keyup.down.native="downSelect"
          ></b-form-input>
          <b-button variant="info"><b-icon-search /></b-button>
        </b-input-group>
      </div>
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
        v-bind:class="{ 'fn-box-active': currentTab === '翻译' }"
        v-on:click="clickTranslation"
      >
        <span class="fn-box-icon">
          <Translate :width="16" :height="16" />
        </span>
        <span class="fn-box-text">翻译</span>
      </div>

      <div
        class="fn-box"
        v-bind:class="{ 'fn-box-active': currentTab === '插件' }"
        v-on:click="clickPlugins"
      >
        <span class="fn-box-icon">
          <Plugins :width="16" :height="16" />
        </span>
        <span class="fn-box-text">插件</span>
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
</template>


<script lang="ts">
import Vue from 'vue';
import Store from '../store/index';
import Translate from '../components/icons/translate.icon.vue';
import Search from '../components/icons/search.icon.vue';
import Plugins from '../components/icons/plugins.icon.vue';
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
    };
  },
  computed: {
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
    selectDictItem(item: DictItem) {
      this.$store.commit('updateCurrentSelectDict', item);
      if (this.searchWord && this.searchWord.length > 0) {
        console.log(this.searchWord);
        this.$store.dispatch('asyncSearchWord', {
          dictid: item.id,
          word: this.searchWord,
        });
      }
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
});
</script>


<style lang="scss" scoped>
.header {
  height: 60px;
  // background-color: #d84042;
  // background-color: #fbfbfb;
  /* background-color: #f6f6f6; */
  background-image: linear-gradient(135deg , #325DFF , #529EFF);

  padding-top: 6px;
  border-bottom: 1px solid #d1d1d1;
  -webkit-app-region: drag;
  .header-navigate-btns {
    height: 54px;
    max-width: 80px;
    padding: 0;
    margin-left: 34px;

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
  .header-search-box {
    max-width: 360px;
    height: 54px;
    padding: 0;
    margin: 0;
       margin-top: 18px;
    &::v-deep {
      .btn-group,
      .btn-group-vertical {
        vertical-align: top;
      }
      // toggle button
      button:nth-child(1) {
        border: 1px solid #f1f1f1;
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
      }
      // search button
      button:nth-child(3) {
        // background-color: #fff;
        // border: 1px solid #fff;

        border: 1px solid #f1f1f1;
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
      }
      input {
        height: 26px;
        padding: 0;
        margin: 0;
        box-shadow: none;
        // border: 1px solid #fff;

        border: 1px solid #f1f1f1;
        background-color: #fff;


      }
      .form-control:disabled,
      .form-control[readonly] {
        background-color: #fff;
      }
    }
  }
  
  .header-functions {
    max-width: 306px;
    height: auto;
    padding: 0;
    margin: 0;
    margin-left: 20px;
    margin-top:12px;
    display: flex;
    flex-direction: row;
    .fn-box-active {
      // background-color: #bd3134;
      // background-color: #e1e1e1;
      background-color: #325effa6;
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
        background-color: #325eff60;
      }

      .fn-box-icon {
        width: 20px;
        height: 20px;
        // color: #f9dad9;
        color: #fff;
        display: block;
        font-size: 16px;
        line-height: 18px;
        margin-left: auto;
        margin-right: auto;
      }
      .fn-box-text {
        // color: #f9dad9;
        color: #fff;
        margin-left: auto;
        margin-right: auto;
        text-align: center;
        width: 100%;
        display: block;
        font-size: 12px;
        font-weight: lighter;
        user-select: none;
      }
    }
  }
}
</style>