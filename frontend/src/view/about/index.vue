<template>
  <SettingPage>
    <section class="settings-section">
      <h2 class="settings-section-title">软件信息</h2>
      <SettingItem title="软件简介">
        <template #desc>Medict 是一款现代、跨平台的本地词典应用</template>
        <span class="about-value">支持 MDX、MDD 与 StarDict 词典格式</span>
      </SettingItem>
      <SettingItem title="当前版本">
        <template #desc>检查当前安装版本及可用更新</template>
        <div class="about-action-row">
          <span class="about-value">{{ latestVersion }}</span>
          <NButton size="small" secondary :loading="checking" @click="checkLatestVersion">
            {{ checking ? '检查中…' : '检查更新' }}
          </NButton>
        </div>
      </SettingItem>
      <SettingItem title="项目主页">
        <template #desc>查看源代码、版本发布和问题反馈</template>
        <NButton size="small" secondary @click="openExternal('https://github.com/terasum/medict')">
          <template #icon><span class="icon icon-github" aria-hidden="true" /></template>
          GitHub
        </NButton>
      </SettingItem>
    </section>

    <section class="settings-section">
      <h2 class="settings-section-title">项目成员</h2>
      <SettingItem title="开发人员">
        <template #desc>Medict 的开发与维护</template>
        <button class="about-link" type="button" @click="openExternal('https://github.com/terasum')">Chen, Quan</button>
      </SettingItem>
      <SettingItem title="设计人员">
        <template #desc>产品视觉与交互设计</template>
        <span class="about-value">Zhang, Mingjiao</span>
      </SettingItem>
      <SettingItem title="特别鸣谢">
        <template #desc>感谢对项目提供帮助的贡献者</template>
        <span class="about-value">Song, Xing</span>
      </SettingItem>
    </section>
  </SettingPage>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { NButton, useDialog, useMessage } from 'naive-ui';
import { BrowserOpenURL } from '../../../wailsjs/runtime/runtime';
import SettingItem from '@/components/setting/SettingItem.vue';
import SettingPage from '@/components/setting/SettingPage.vue';

const latestVersion = ref('3.0.1-alpha');
const checking = ref(false);
const dialog = useDialog();
const message = useMessage();

function openExternal(url: string) {
  BrowserOpenURL(url);
}

async function checkLatestVersion() {
  checking.value = true;
  await Promise.resolve();
  checking.value = false;
  dialog.info({
    title: '检查更新',
    content: `当前版本 ${latestVersion.value}，暂未发现可用更新。`,
    positiveText: '知道了',
    onPositiveClick: () => message.success('已完成检查'),
  });
}
</script>

<style scoped>
.about-value {
  color: var(--c-gray-800);
}

.about-action-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.about-action-row :deep(.n-button),
:deep(.n-button) {
  background-image: none;
  box-shadow: none;
  font-size: 12px;
}

.about-link {
  padding: 0;
  color: var(--c-primary);
  background: transparent;
  border: 0;
  font: inherit;
  cursor: pointer;
}

.about-link:hover {
  text-decoration: underline;
}
</style>
