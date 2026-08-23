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
@use './style/variables.scss' as *;
@use '@/style/photon/photon.scss';
@use '@/style/settings.scss';

#app-root {
  height: 100%;
  width: 100%;
  padding: 0;
  margin: 0;
  display: block;
  overflow: hidden;

  // macOS hides the native title bar (mac.TitleBarHidden()), so the in-app
  // strip keeps the full height. Windows/Linux cannot hide the native title
  // bar, so the strip shrinks by 12px there (.os-windows / .os-linux).
  --fake-title-bar-height: #{$fake-title-bar-height};

  &.os-windows,
  &.os-linux {
    --fake-title-bar-height: #{$fake-title-bar-height - 12px};
  }

  .fake-title-bar {
    width: 100%;
    height: var(--fake-title-bar-height);
    display: block;
    --wails-draggable: drag;
    background: transparent;
    // background-color: var(--c-gray-100);

    background-color: $theme-top-header-background-color;
  }
  .x-space-provider {
    width: 100%;
    height: calc(100% - var(--fake-title-bar-height));
  }
}
</style>

<template>
  <div id="app-root" class="app-container" :class="platformClass">
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
          <CSSWindow v-if="windowMode === 'css-editor'" />
          <router-view v-else-if="windowMode === 'main'"></router-view>
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
import { BRAND, palette } from '@/style/tokens';
import CSSWindow from '@/view/css-editor/index.vue';
import { detectPlatform, platformClass as toPlatformClass, sniffPlatformSync } from '@/utils/platform';

let isDark = ref(false);
const windowMode = ref<'loading' | 'main' | 'css-editor'>('loading');
// synchronous guess first so the first paint already has the right height;
// refined with the authoritative runtime.GOOS answer in onMounted
const platformClass = ref(toPlatformClass(sniffPlatformSync()));
let theme = reactive(light);

if (isDark.value) {
  theme = dark;
}

const themeOverrides: GlobalThemeOverrides = {
  common: {
    primaryColor: BRAND,
  },
  Input: {
    borderFocus: `1px solid ${BRAND}`,
    borderHover: `1px solid ${palette.primaryHover}`,
  },
  Button: {
    textColor: palette.gray[900],
    textColorHoverPrimary: BRAND,
    textColorPressedPrimary: BRAND,
    textColorFocusPrimary: BRAND,
    border: 'none',
    borderHover: 'none',
    borderPressed: 'none',
    borderFocus: 'none',
    borderDisabled: 'none',
  },
  Dialog: {
    iconSize: '0px',
  },
};




function listenStoreChange(store: any) {
  const unscribe = store.$onAction(
    ({
      name, // action 名称
      store, // store 实例，类似 `someStore`
      args, // 传递给 action 的参数数组
      after, // 在 action 返回或解决后的钩子
      onError, // action 抛出或拒绝的钩子
    }: any) => {
      let startTime = Date.now();
      console.debug(`[store-action] {${name}} triggered started, args: {${args}}`);

      // 这将在 action 成功并完全运行后触发。
      // 它等待着任何返回的 promise
      after((result: any) => {
        console.debug(
          `[store-action] {${name}} triggered success, after ${
            Date.now() - startTime
          }ms, with result ${result}.`
        );
      });

      // 如果 action 抛出或返回一个拒绝的 promise，这将触发
      onError((error: any) => {
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

let unscribeDictQueryStore: (() => void) | null = null;
const dictQueryStore = useDictQueryStore();
onMounted(()=>{
  detectPlatform().then((os) => {
    platformClass.value = toPlatformClass(os);
  });
  const modeCall = (window as any)?.go?.main?.App?.WindowMode;
  if (typeof modeCall === 'function') {
    modeCall().then((mode: string) => {
      windowMode.value = mode === 'css-editor' ? 'css-editor' : 'main';
    }).catch(() => { windowMode.value = 'main'; });
  } else {
    windowMode.value = 'main';
  }
  unscribeDictQueryStore = listenStoreChange(dictQueryStore);

})

onUnmounted(() =>{
  if (unscribeDictQueryStore) {
    unscribeDictQueryStore();
    unscribeDictQueryStore = null;
  }
})



</script>
