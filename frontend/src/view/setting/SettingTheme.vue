<template>
    <div class="setting-main-container" ref="containerRef">
      <div class="setting-main-container-header">
        <h3>主题设置</h3>
      </div>
  
      <div class="setting-main-container-content">
        <n-card class="setting-section">
        </n-card>
      </div>
    </div>
  </template>
  <script setup>
  import { NCard, NAffix, NTag } from 'naive-ui';
  import { ref } from 'vue';
  import SettingItem from '@/components/setting/SettingItem.vue';
  
  const containerRef = (ref < HTMLElement) | (undefined > undefined);
  </script>
  
  <style lang="scss" scoped>
  @import '@/style/variables.scss';
  @import '@/style/photon/photon.scss';
  
  .setting-main-container {
    height: 100%;
    width: 100%;
    overflow-y: auto;
    margin: 0;
    padding: 0;
    .setting-main-container-header {
      display: flex;
      z-index: 99;
      background: #fff;
      padding-left: 10px;
    }
    .setting-main-container-content {
      padding: 10px;
      .setting-section{
          margin: 6px auto;
      }
    }
  }
  </style>
  