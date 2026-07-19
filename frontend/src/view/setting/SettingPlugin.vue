<template>
  <SettingPage>
    <section class="settings-section">
      <h2 class="settings-section-title">访问权限</h2>
      <SettingItem title="允许访问的域名">
        <template #desc>每行一个域名，不包含协议和路径</template>
        <EditableSettingValue
          v-model="allowedDomains"
          aria-label="插件允许访问的域名"
          multiline
          :rows="5"
          :on-save="saveDomains"
        />
      </SettingItem>
    </section>
    <section class="settings-section">
      <h2 class="settings-section-title">插件目录</h2>
      <SettingItem title="本地插件">
        <template #desc>插件功能仍在完善中，目录将在启用后显示</template>
        <span class="muted-value">尚未启用</span>
      </SettingItem>
    </section>
  </SettingPage>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { getPreferences, savePreferences } from '@/apis/config';
import EditableSettingValue from '@/components/setting/EditableSettingValue.vue';
import SettingItem from '@/components/setting/SettingItem.vue';
import SettingPage from '@/components/setting/SettingPage.vue';

const allowedDomains = ref('localhost');
const saveDomains = (value: string) => savePreferences({ plugindomainallowlist: value });

onMounted(async () => {
  const preferences = await getPreferences();
  allowedDomains.value = String(preferences.plugindomainallowlist ?? 'localhost');
});
</script>

<style scoped>
.muted-value { color: var(--c-gray-500); }
</style>
