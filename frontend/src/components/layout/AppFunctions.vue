<style lang="scss">
@import '@/style/variables.scss';
@import '@/style/photon/photon.scss';

.app-content-functions {
  height: $layout-header-height;
  display: flex;
  width: 100%;
  flex-direction: row;

  .header {
    height: 100%;
    display: flex;
    width: 100%;
    flex-direction: row;
    justify-content: space-between;
    // background-color: $theme-top-header-background-color;
    // border-bottom: 1px solid var(--c-gray-300);

    .header-nav-functions {
      max-width: 280px;
      height: 100%;
      padding: 0;
      margin: 0;
      margin-left: 8px;
      margin-top: 0;
      display: flex;
      align-items: center;
      flex-direction: row;
      .fn-box-active {
        background-color: $theme-function-box-active-color;
      }
      .fn-box {
        width: 46px;
        height: 44px;
        margin-top: 0;
        padding: 4px 0 5px;
        box-sizing: border-box;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 2px;
        cursor: pointer;
        border-radius: 4px;
        margin-left: 2px;

        color: $theme-function-box-font-color;

        &:hover {
          background-color: $theme-function-box-hover-bg-color;
          color: $theme-function-box-hover-font-color;
        }

        .fn-box-icon {
          width: 18px;
          height: 18px;
          display: block;
          font-size: 16px;
          line-height: 18px;
          flex: 0 0 18px;
        }
        .fn-box-text {
          margin: 0;
          text-align: center;
          width: 100%;
          display: block;
          font-size: 12px;
          line-height: 14px;
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


      <div class="header-nav-functions">
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

        <div class="fn-box" @click="changeTab('bookmarks')" :class="uiStore.currentTab == 'bookmarks'?'active':''">
          <span class="fn-box-icon">
            <Star />
          </span>
          <span class="fn-box-text">生词</span>
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
import { Search, Book, ToggleOn, Star } from '@vicons/fa';
import { useDictQueryStore } from '@/store/dict';
import { useUIStore } from '@/store/ui';
import { useRouter } from "vue-router";

const uiStore = useUIStore();
const router = useRouter();

///----------------------------
// event listener function
///----------------------------

const tabRouters = {
  search: '/',
  dict: '/dict',
  bookmarks: '/bookmarks',
  setting: '/setting',
};

function changeTab(tabName) {
  if (uiStore.currentTab != tabName && tabRouters[tabName]) {
    router.replace({ path: tabRouters[tabName] });
  }
  uiStore.updateCurrentTab(tabName);
}
</script>
