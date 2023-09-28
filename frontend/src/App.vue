<!--

 Copyright (C) 2023 Quan Chen <chenquan_act@163.com>

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.

 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.
-->

<style lang="scss">
@import './style/variables.scss';

#app-root {
  height: 100%;
  width: 100%;
  padding: 0;
  margin: 0;
  display: block;
  overflow: hidden;

  .fake-title-bar {
    width: 100%;
    height: $fake-title-bar-height;
    display: block;
    --wails-draggable: drag;
    background: transparent;
    // background-color: #fafafa;

    background-color: $theme-top-header-background-color;
  }
  .x-space-provider {
    width: 100%;
    height: calc(100% - $fake-title-bar-height);
  }
}
</style>

<template>
  <div id="app-root" class="app-container">
    <div class="fake-title-bar" data-wails-drag></div>
    <n-config-provider
      :theme="theme"
      :locale="zhCN"
      :date-locale="dateZhCN"
      class="x-space-provider"
      :theme-overrides="themeOverrides"
    >
      <n-dialog-provider>
        <n-message-provider>
          <router-view></router-view>
        </n-message-provider>
      </n-dialog-provider>
      <n-global-style />
    </n-config-provider>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue';
import {
  NConfigProvider,
  NGlobalStyle,
  NDialogProvider,
  NMessageProvider,
} from 'naive-ui';
import { darkTheme as dark, lightTheme as light } from 'naive-ui';
import { zhCN, dateZhCN } from 'naive-ui';
import { GlobalThemeOverrides } from 'naive-ui';
import { useDictQueryStore } from './store/dict';

let isDark = ref(false);
let theme = reactive(light);

if (isDark.value) {
  theme = dark;
}

const themeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: '#326cb8',
  },
  Input: {
    borderFocus: '1px solid #326cb8',
    borderHover: '1px solid #326cc9',
  },
  Button: {
    textColor: '#333',
    // borderFocus: '1px solid #326cb8',
    // borderHover: '1px solid #326cc9',
    textColorHoverPrimary: '#326cb8',
    textColorPressedPrimary: '#326cb8',
    textColorFocusPrimary: '#326cb8',
    border: 'none',
    borderHover: 'none',
    borderPressed: 'none',
    borderFocus: 'none',
    borderDisabled: 'none',
  
  },
  Dialog: {
    // iconColor: string;
    // iconColorInfo: string;
    // iconColorSuccess: "#326cb8",
    iconSize: "0px",
    // iconColorWarning: string;
    // iconColorError: string;
  }
};




function listenStoreChange(store: any) {
  const unscribe = store.$onAction(
    ({
      name, // action 名称
      store, // store 实例，类似 `someStore`
      args, // 传递给 action 的参数数组
      after, // 在 action 返回或解决后的钩子
      onError, // action 抛出或拒绝的钩子
    }) => {
      let startTime = Date.now();
      console.debug(`[store-action] {${name}} triggered started, args: {${args}}`);

      // 这将在 action 成功并完全运行后触发。
      // 它等待着任何返回的 promise
      after((result) => {
        console.debug(
          `[store-action] {${name}} triggered success, after ${
            Date.now() - startTime
          }ms, with result ${result}.`
        );
      });

      // 如果 action 抛出或返回一个拒绝的 promise，这将触发
      onError((error) => {
        console.warn(
          `[store-action] {${name}} trigger faild, after ${
            Date.now() - startTime
          }ms.\nerror: ${error}.`
        );
      });
    }
  );
  return unscribe;
}

let unscribeDictQueryStore = null;
const dictQueryStore = useDictQueryStore();
onMounted(()=>{
  unscribeDictQueryStore = listenStoreChange(dictQueryStore);

})

onUnmounted(() =>{
  if (unscribeDictQueryStore) {
    unscribeDictQueryStore();
    unscribeDictQueryStore = null;
  }
})



</script>
