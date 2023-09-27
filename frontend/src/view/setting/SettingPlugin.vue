<template>
    <div class="setting-main-container" ref="containerRef">
      <div class="setting-main-container-header">
        <h3>插件设置</h3>
      </div>
  
      <div class="setting-main-container-content">

      <n-card class="setting-section">
          <SettingItem title="前端域白名单" :value="dictDir">
            <template #desc>
            前端插件允许访问的域名 (不支持IP地址)
            </template>
            <template #action>
            <button class="btn btn-default" @click="OpenDictDir">修改</button>
            <button class="btn btn-default" @click="OpenDictDir">删除</button>
            </template>
              <table class="table-striped table-wrapper">
                <thead><tr><th width="20">选择</th><th width="20">序号</th><th>允许域名</th></tr></thead>
                <tbody>
                  <tr><td><input width="20" type="checkbox"/></td><td width="20">1</td><td>localhost</td></tr>
                  <tr><td><input width="20" type="checkbox"/></td><td width="20">1</td><td>localhost</td></tr>
                  <tr><td><input width="20" type="checkbox"/></td><td width="20">1</td><td>localhost</td></tr>
                </tbody>
              </table>
          </SettingItem>

      </n-card>
      <n-card class="setting-section">
          <SettingItem title="插件目录">
             <template #desc>
              插件目录, 用于存储所有的插件
            </template>
            <template #action>
            <button class="btn btn-default" @click="OpenDictDir">打开</button>
            </template>
            {{ PluginDir }}
          </SettingItem>
      </n-card>


      </div>
    </div>
  </template>
  <script setup>
  import { NCard, NAffix, NTag } from 'naive-ui';
  import { ref } from 'vue';
  import SettingItem from '@/components/setting/SettingItem.vue';
  
  const containerRef = (ref < HTMLElement) | (undefined > undefined);

  let PluginDir = ref("C:\\Users\\Administrator\\AppData\\Roaming\\medict\\plugins");
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
          .table-wrapper{
            border: 1px solid #ccc;
          }
      }
    }
  }
  </style>
  