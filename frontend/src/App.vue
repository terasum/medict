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
    -webkit-app-region: drag;
    background: transparent;
    cursor: grab;
    background-color: #fafafa;
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
      <router-view></router-view>
      <n-global-style />
    </n-config-provider>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive } from 'vue';
import { NConfigProvider, NGlobalStyle } from 'naive-ui';
import { darkTheme as dark, lightTheme as light } from 'naive-ui';
import { zhCN, dateZhCN } from 'naive-ui';
import { GlobalThemeOverrides } from 'naive-ui'

let isDark = ref(false);
let theme = reactive(light);

if (isDark.value){
   theme = dark;
}

  const themeOverrides: GlobalThemeOverrides = {
    common: {
      primaryColor: '#326cb8'
    },
    Input:{
      borderFocus: '1px solid #326cb8',
      borderHover: '1px solid #326cc9',
    },
    Button: {
      textColor: '#FF0000'
    }
  }
</script>