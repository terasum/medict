<template>
    <div class="setting-main-container" ref="containerRef">
      <div class="setting-main-container-header">
        <h3>检查更新</h3>
      </div>
  
      <div class="setting-main-container-content">
        <n-card class="setting-section">
            <SettingItem title="最新版本">
              <template #desc>
              </template>
              <template #action>
              <button class="btn btn-default" @click="checkLatestVersion">检查新版本</button>
              </template>
                {{latestVersion}}
            </SettingItem>
  
        </n-card>
      </div>
    </div>
  </template>
  <script lang="ts" setup>
  import { NCard } from 'naive-ui';
  import { useDialog,useMessage  } from 'naive-ui'

  import { ref,onMounted } from 'vue';
  import SettingItem from "@/components/setting/SettingItem.vue";
  
  const containerRef = ref<any>();

  const latestVersion = ref("3.0.1-alpha");

  const dialog = useDialog();
  const message = useMessage();


  function checkLatestVersion() {
    dialog.success({
          title: '检查更新',
          content: 'message',
          positiveText: '确认',
          negativeText: "取消",
          onPositiveClick: () => {
            message.success('确认')
          }
        })
  }
 
  onMounted(()=>{
 
  })
  
  
  
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
  