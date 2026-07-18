<template>
  <SettingPage title="版本更新" description="查看当前版本并检查可用更新。" back>
    <section class="settings-section">
      <h2 class="settings-section-title">版本信息</h2>
      <SettingItem title="当前版本">
        <template #desc>Medict 当前安装的应用版本</template>
        <div class="update-row">
          <span>{{ latestVersion }}</span>
          <button class="btn btn-default" type="button" @click="checkLatestVersion">检查更新</button>
        </div>
      </SettingItem>
    </section>
  </SettingPage>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { useDialog, useMessage } from 'naive-ui';
import SettingItem from '@/components/setting/SettingItem.vue';
import SettingPage from '@/components/setting/SettingPage.vue';

const latestVersion = ref('3.0.1-alpha');
const dialog = useDialog();
const message = useMessage();

function checkLatestVersion() {
  dialog.info({
    title: '检查更新',
    content: `当前版本 ${latestVersion.value}，暂未发现可用更新。`,
    positiveText: '知道了',
    onPositiveClick: () => message.success('已完成检查'),
  });
}
</script>

<style scoped>
.update-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}
</style>
