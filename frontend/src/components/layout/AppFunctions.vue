<style lang="scss">
@import '@/style/variables.scss';

.app-content-functions {
  height: $layout-header-height;
  display: flex;
  width: 100%;
  flex-direction: row;
  background: #fafafa;

  .header {

    height: calc($layout-header-height - 1px);
    display: flex;
    width: 100%;
    flex-direction: row;
    justify-content: space-between;
    background-color: $theme-top-header-background-color;
    border-bottom: 1px solid #dcdcdc;

    .header-main-function {
      display: flex;
      height: 54px;
      padding: 0;
      margin: 0;
      min-width: 174px;
    }

    .header-nav-functions {
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
      <div class="header-main-function">
        <slot></slot>
      </div>

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
import { Search, Book, ToggleOn } from '@vicons/fa';
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
  setting: '/setting',
};

function changeTab(tabName) {
  if (uiStore.currentTab != tabName && tabRouters[tabName]) {
    router.replace({ path: tabRouters[tabName] });
  }
  uiStore.updateCurrentTab(tabName);
}
</script>
