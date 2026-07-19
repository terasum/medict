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
      <SettingItem title="完整英汉词库">
        <template #desc>当前内置 5 万高频词；可从官方 ECDICT（MIT）下载完整数据并在本地离线查询</template>
        <div class="setting-value-with-action">
          <span class="setting-path">{{ ecdictStatusText }}</span>
          <button
            class="btn btn-default"
            type="button"
            :disabled="installingECDICT || ecdictStatus?.edition === 'full'"
            @click="installECDICT"
          >
            {{ installingECDICT ? '下载并安装中…' : ecdictStatus?.edition === 'full' ? '已安装' : '下载完整词库' }}
          </button>
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
import { computed, onMounted, ref } from 'vue';
import { useMessage } from 'naive-ui';
import { BaseDictDirectory, OpenDirOrFile } from '@/apis/apis';
import { getPreferences, savePreferences } from '@/apis/config';
import { getECDICTStatus, installFullECDICT, type ECDICTStatus } from '@/apis/ecdict';
import EditableSettingValue from '@/components/setting/EditableSettingValue.vue';
import SettingItem from '@/components/setting/SettingItem.vue';
import SettingPage from '@/components/setting/SettingPage.vue';

const defaultPreset = `<link href="#{DictName}.css?dict_id=#{DictID}" rel="stylesheet" />\n<script async src="#{DictName}.js?dict_id=#{DictID}"><\/script>`;
const dictDir = ref('');
const allowSymlinks = ref('true');
const groupLimit = ref('10');
const presetContent = ref(defaultPreset);
const ecdictStatus = ref<ECDICTStatus | null>(null);
const installingECDICT = ref(false);
const message = useMessage();
const booleanOptions = [
  { label: '允许', value: 'true' },
  { label: '不允许', value: 'false' },
];

const ecdictStatusText = computed(() => {
  if (!ecdictStatus.value) return '正在读取词库状态…';
  const count = new Intl.NumberFormat('zh-CN').format(ecdictStatus.value.entryCount);
  return ecdictStatus.value.edition === 'full' ? `完整版本 · ${count} 词条` : `精简版本 · ${count} 词条`;
});

function saveValue(key: string) {
  return (value: string) => savePreferences({ [key]: value });
}

async function openDictDir() {
  if (dictDir.value) await OpenDirOrFile(dictDir.value);
}

async function installECDICT() {
  installingECDICT.value = true;
  try {
    ecdictStatus.value = await installFullECDICT();
    message.success(`完整 ECDICT 已安装，共 ${ecdictStatus.value.entryCount.toLocaleString('zh-CN')} 个词条`);
  } catch (cause) {
    message.error(cause instanceof Error ? cause.message : '完整词库安装失败');
  } finally {
    installingECDICT.value = false;
  }
}

onMounted(async () => {
  const [dir, preferences, status] = await Promise.all([BaseDictDirectory(), getPreferences(), getECDICTStatus()]);
  dictDir.value = dir;
  ecdictStatus.value = status;
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
