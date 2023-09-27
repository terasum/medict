<template>
  <div class="setting-main-container" ref="containerRef">
    <div class="setting-main-container-header">
      <h3>软件设置</h3>
    </div>

    <div class="setting-main-container-content">
      <n-card class="setting-section">
        <SettingItem title="全文搜索引擎">
          <template #desc>
          用于进行全文搜索的内置引擎
          </template>
          bleve
        </SettingItem>
        <SettingItem title="全文索引超时时间">
          <template #desc>
          全文搜索的内置引擎建立索引的超时时间
          </template>
          <template #action>
            <button class="btn btn-default">修改</button>
          </template>
          10s
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
