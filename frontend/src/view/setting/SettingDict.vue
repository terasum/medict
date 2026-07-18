<template>
  <SettingPage title="词典设置" description="管理词典目录、加载规则与内容模板。">
    <section class="settings-section">
      <h2 class="settings-section-title">存储位置</h2>
      <SettingItem title="词典目录">
        <template #desc>Medict 扫描和存放词典文件的位置</template>
        <div class="setting-value-with-action">
          <span class="setting-path">{{ dictDir || '正在读取…' }}</span>
          <button class="btn btn-default" type="button" :disabled="!dictDir" @click="openDictDir">在访达中打开</button>
        </div>
      </SettingItem>
    </section>

    <section class="settings-section">
      <h2 class="settings-section-title">加载规则</h2>
      <SettingItem title="允许软链接词典">
        <template #desc>扫描词典目录时包含通过软链接引用的词典</template>
        <EditableSettingValue
          v-model="allowSymlinks"
          aria-label="允许软链接词典"
          :options="booleanOptions"
          :on-save="saveValue('dictionaryallowsymlinks')"
        />
      </SettingItem>
      <SettingItem title="单个词典组上限">
        <template #desc>限制每组词典数量，避免一次建立过多索引</template>
        <EditableSettingValue
          v-model="groupLimit"
          aria-label="单个词典组上限"
          input-type="number"
          :min="1"
          :on-save="saveValue('dictionarygroupmax')"
        />
      </SettingItem>
    </section>

    <section class="settings-section">
      <h2 class="settings-section-title">高级内容</h2>
      <SettingItem title="词典内容预置 CSS / JS">
        <template #desc>注入词典内容页 head 的模板；修改后保存到应用配置</template>
        <EditableSettingValue
          v-model="presetContent"
          aria-label="词典内容预置 CSS 和 JavaScript"
          multiline
          :rows="6"
          :on-save="saveValue('dictionarypresetcontent')"
        />
      </SettingItem>
    </section>
  </SettingPage>
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue';
import { BaseDictDirectory, OpenDirOrFile } from '@/apis/apis';
import { getPreferences, savePreferences } from '@/apis/config';
import EditableSettingValue from '@/components/setting/EditableSettingValue.vue';
import SettingItem from '@/components/setting/SettingItem.vue';
import SettingPage from '@/components/setting/SettingPage.vue';

const defaultPreset = `<link href="#{DictName}.css?dict_id=#{DictID}" rel="stylesheet" />\n<script async src="#{DictName}.js?dict_id=#{DictID}"><\/script>`;
const dictDir = ref('');
const allowSymlinks = ref('true');
const groupLimit = ref('10');
const presetContent = ref(defaultPreset);
const booleanOptions = [
  { label: '允许', value: 'true' },
  { label: '不允许', value: 'false' },
];

function saveValue(key: string) {
  return (value: string) => savePreferences({ [key]: value });
}

async function openDictDir() {
  if (dictDir.value) await OpenDirOrFile(dictDir.value);
}

onMounted(async () => {
  const [dir, preferences] = await Promise.all([BaseDictDirectory(), getPreferences()]);
  dictDir.value = dir;
  allowSymlinks.value = String(preferences.dictionaryallowsymlinks ?? 'true');
  groupLimit.value = String(preferences.dictionarygroupmax ?? '10');
  presetContent.value = String(preferences.dictionarypresetcontent ?? defaultPreset);
});
</script>

<style lang="scss" scoped>
.setting-value-with-action {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.setting-path {
  min-width: 0;
  color: var(--c-gray-700);
  overflow-wrap: anywhere;
}
</style>
