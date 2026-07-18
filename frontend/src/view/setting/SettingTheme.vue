<template>
  <SettingPage title="外观设置" description="调整 Medict 的显示方式。">
    <section class="settings-section">
      <h2 class="settings-section-title">主题</h2>
      <SettingItem title="界面主题">
        <template #desc>当前跟随系统外观；更多主题选项将在后续版本提供</template>
        <span class="theme-status">跟随系统</span>
      </SettingItem>
    </section>
  </SettingPage>
</template>

<script setup lang="ts">
import SettingItem from '@/components/setting/SettingItem.vue';
import SettingPage from '@/components/setting/SettingPage.vue';
</script>

<style scoped>
.theme-status { color: var(--c-gray-700); }
</style>
