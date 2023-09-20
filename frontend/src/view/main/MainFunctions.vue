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
        padding-top: 6px;
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
          margin-top: 2px;
          margin-left: auto;
          margin-right: auto;
          text-align: center;
          width: 100%;
          display: block;
          font-size: 12px;
          user-select: none;
        }
      }

      .active {
        background-color: $theme-function-box-hover-bg-color;
        color: $theme-function-box-hover-font-color;
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
            @keydown.enter="handleChange"
            v-model:value="inputWord"
          >
            <template #suffix>
              <n-icon :component="Search" />
            </template>
          </n-input>
        </div>
      </div>

      <div class="header-functions">
        <div class="fn-box" @click="changeTab('search')"  :class="uiStore.currentTab == 'search'?'active':''">
          <span class="fn-box-icon">
            <Search />
          </span>
          <span class="fn-box-text">搜索</span>
        </div>

        <div class="fn-box" @click="changeTab('dict')" :class="uiStore.currentTab == 'dict'?'active':''">
          <span class="fn-box-icon">
            <Book />
          </span>
          <span class="fn-box-text">词典</span>
        </div>

        <div class="fn-box" @click="changeTab('setting')"  :class="uiStore.currentTab == 'setting'?'active':''">
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

import { useDictQueryStore } from '@/store/dict';
import { useUIStore } from '@/store/ui';
import { useRouter } from "vue-router";

const dictQueryStore = useDictQueryStore();
const uiStore = useUIStore();
const router = useRouter();

let inputWord = ref('');
let inputActive = ref(false);

///----------------------------
// event listener function
///----------------------------

const tabRouters = {
  search: '/',
  dict: '/dict',
  setting: '/setting',
};

function changeTab(tabName) {
  if (uiStore.currentTab != tabName && tabRouters[tabName]) {
    router.replace({ path: tabRouters[tabName] });
  }
  uiStore.updateCurrentTab(tabName);
}

function handleChange(v) {
  console.info('[event keydown enter]' + inputWord.value);
  if(!uiStore.isSearchInputActive()){
    console.log("input disabled, skipped")
    return;
  }
  dictQueryStore.updateInputSearchWord(inputWord.value.trim());
}
</script>
