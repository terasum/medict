<template>
  <div class="x-space">
    <div class="x-layout">
      <div class="x-layout-header"><AppHeader /></div>
      <div class="x-layout-main-area">
        <aside class="x-layout-sidebar">
          <AppSidebar>
            <DictionaryGroupSidebar
              :groups="sidebarGroups"
              :selected-group-id="selectedGroupId"
              :favorite-group-id="favoriteGroupId"
              :default-group-id="DEFAULT_DICTIONARY_GROUP_ID"
              @select="selectedGroupId = $event"
              @create="openCreateGroup"
              @delete="confirmDeleteGroup"
              @toggle-favorite="toggleFavoriteGroup"
            />
          </AppSidebar>
        </aside>

        <main class="x-layout-content">
          <AppMainContent>
            <div class="dictionary-content-toolbar">
              <div>
                <strong>{{ selectedGroupName }}</strong>
                <span>{{ visibleDicts.length }} 个词典</span>
              </div>
              <NButton v-if="selectedGroupId !== DEFAULT_DICTIONARY_GROUP_ID" size="small" secondary @click="openMembers">
                管理成员
              </NButton>
            </div>

            <div class="dict-main-area">
              <table v-if="visibleDicts.length" class="table-striped">
                <thead>
                  <tr>
                    <th>词典名称</th>
                    <th>标题（内置）</th>
                    <th>词典类型</th>
                    <th>引擎版本</th>
                    <th>创建时间</th>
                    <th>id</th>
                    <th>所在目录</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="item in visibleDicts" :key="item.id">
                    <td>{{ item.name }}</td>
                    <td>{{ item.title }}</td>
                    <td>{{ item.dictType }}</td>
                    <td>{{ item.generateEngineVersion }}</td>
                    <td>{{ item.createDate }}</td>
                    <td>{{ item.id }}</td>
                    <td>{{ item.baseDir }}</td>
                  </tr>
                </tbody>
              </table>
              <div v-else class="dictionary-empty">
                <span class="icon icon-book" aria-hidden="true" />
                <strong>{{ selectedGroupId === DEFAULT_DICTIONARY_GROUP_ID ? '尚未加载到词典' : '当前词典组为空' }}</strong>
                <span>{{ selectedGroupId === DEFAULT_DICTIONARY_GROUP_ID ? '请检查词典目录或重新加载应用。' : '点击“管理成员”向该组添加词典。' }}</span>
                <NButton v-if="selectedGroupId !== DEFAULT_DICTIONARY_GROUP_ID" size="small" secondary @click="openMembers">
                  管理成员
                </NButton>
              </div>
            </div>
          </AppMainContent>
        </main>
      </div>

      <div class="n-layout-footer">
        <AppFooter><span class="dicts-hint">当前组拥有词典数: {{ visibleDicts.length }}</span></AppFooter>
      </div>
    </div>

    <NModal v-model:show="createGroupVisible" preset="dialog" title="新建词典组" :show-icon="false">
      <NInput
        ref="createGroupInput"
        v-model:value="newGroupName"
        placeholder="词典组名称"
        maxlength="30"
        @keydown.enter="createGroup"
      />
      <template #action>
        <NButton size="small" @click="createGroupVisible = false">取消</NButton>
        <NButton size="small" type="primary" :loading="saving" @click="createGroup">创建</NButton>
      </template>
    </NModal>

    <NModal v-model:show="membersVisible" preset="dialog" :title="`管理「${selectedGroupName}」成员`" :show-icon="false">
      <div class="member-dialog-hint">选择需要显示在该词典组中的词典。</div>
      <NCheckboxGroup v-model:value="memberDraft" class="member-checkbox-list">
        <NCheckbox v-for="dict in dictsList" :key="dict.id" :value="dict.id" :label="dict.name" />
      </NCheckboxGroup>
      <template #action>
        <NButton size="small" @click="membersVisible = false">取消</NButton>
        <NButton size="small" type="primary" :loading="saving" @click="saveMembers">保存</NButton>
      </template>
    </NModal>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, ref } from 'vue';
import { NButton, NCheckbox, NCheckboxGroup, NInput, NModal, useDialog, useMessage } from 'naive-ui';
import AppHeader from '@/components/layout/AppHeader.vue';
import AppSidebar from '@/components/layout/AppSidebar.vue';
import AppFooter from '@/components/layout/AppFooter.vue';
import AppMainContent from '@/components/layout/AppMainContent.vue';
import { useUIStore } from '@/store/ui';
import { GetAllDicts } from '@/apis/dicts-api';
import { getPreferences, savePreferences } from '@/apis/config';
import DictionaryGroupSidebar from './DictionaryGroupSidebar.vue';
import {
  DEFAULT_DICTIONARY_GROUP_ID,
  type DictionaryGroup,
  createDictionaryGroup,
  deleteDictionaryGroup,
  filterDictionaryGroupMembers,
  parseDictionaryGroups,
  updateDictionaryGroupMembers,
} from './dictionary-groups';

interface DictionaryRow {
  id: string;
  name: string;
  dictType: string;
  title: string;
  generateEngineVersion: string;
  createDate: string;
  baseDir: string;
}

const GROUPS_PREFERENCE_KEY = 'dictionarygroupsjson';
const FAVORITE_GROUP_PREFERENCE_KEY = 'dictionaryfavoritegroupid';

const uiStore = useUIStore();
const message = useMessage();
const dialog = useDialog();
const dictsList = ref<DictionaryRow[]>([]);
const customGroups = ref<DictionaryGroup[]>([]);
const selectedGroupId = ref(DEFAULT_DICTIONARY_GROUP_ID);
const favoriteGroupId = ref('');
const createGroupVisible = ref(false);
const newGroupName = ref('');
const createGroupInput = ref<InstanceType<typeof NInput> | null>(null);
const membersVisible = ref(false);
const memberDraft = ref<string[]>([]);
const saving = ref(false);

uiStore.updateCurrentTab('dict');

const selectedCustomGroup = computed(() => customGroups.value.find((group) => group.id === selectedGroupId.value) || null);
const selectedGroupName = computed(() => selectedCustomGroup.value?.name || '默认组');
const visibleDicts = computed(() => filterDictionaryGroupMembers(dictsList.value, selectedCustomGroup.value));
const sidebarGroups = computed(() => [
  { id: DEFAULT_DICTIONARY_GROUP_ID, name: '默认组', count: dictsList.value.length },
  ...customGroups.value.map((group) => ({ id: group.id, name: group.name, count: group.dictIds.length })),
]);

function groupID() {
  return globalThis.crypto?.randomUUID?.() || `group-${Date.now()}-${Math.random().toString(16).slice(2)}`;
}

async function persistGroups() {
  await savePreferences({
    [GROUPS_PREFERENCE_KEY]: JSON.stringify(customGroups.value),
    [FAVORITE_GROUP_PREFERENCE_KEY]: favoriteGroupId.value,
  });
}

function openCreateGroup() {
  newGroupName.value = '';
  createGroupVisible.value = true;
  nextTick(() => createGroupInput.value?.focus());
}

async function createGroup() {
  const previousGroups = customGroups.value;
  const previousSelected = selectedGroupId.value;
  try {
    const next = createDictionaryGroup(customGroups.value, newGroupName.value, groupID());
    saving.value = true;
    customGroups.value = next;
    selectedGroupId.value = next.at(-1)?.id || DEFAULT_DICTIONARY_GROUP_ID;
    await persistGroups();
    createGroupVisible.value = false;
    message.success('词典组已创建');
  } catch (cause) {
    customGroups.value = previousGroups;
    selectedGroupId.value = previousSelected;
    message.warning(cause instanceof Error ? cause.message : '创建词典组失败');
  } finally {
    saving.value = false;
  }
}

function confirmDeleteGroup() {
  if (!selectedCustomGroup.value) return;
  const target = selectedCustomGroup.value;
  dialog.warning({
    title: '删除词典组',
    content: `确定删除「${target.name}」？词典文件不会被删除。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      const previousGroups = customGroups.value;
      const previousFavorite = favoriteGroupId.value;
      const previousSelected = selectedGroupId.value;
      try {
        customGroups.value = deleteDictionaryGroup(customGroups.value, target.id);
        if (favoriteGroupId.value === target.id) favoriteGroupId.value = '';
        selectedGroupId.value = DEFAULT_DICTIONARY_GROUP_ID;
        await persistGroups();
        message.success('词典组已删除');
      } catch (cause) {
        customGroups.value = previousGroups;
        favoriteGroupId.value = previousFavorite;
        selectedGroupId.value = previousSelected;
        message.error(cause instanceof Error ? cause.message : '删除词典组失败');
      }
    },
  });
}

async function toggleFavoriteGroup() {
  const previous = favoriteGroupId.value;
  favoriteGroupId.value = previous === selectedGroupId.value ? '' : selectedGroupId.value;
  try {
    await persistGroups();
    message.success(favoriteGroupId.value ? `「${selectedGroupName.value}」已设为常用组` : '已取消常用组');
  } catch (cause) {
    favoriteGroupId.value = previous;
    message.error(cause instanceof Error ? cause.message : '保存常用组失败');
  }
}

function openMembers() {
  if (!selectedCustomGroup.value) return;
  memberDraft.value = [...selectedCustomGroup.value.dictIds];
  membersVisible.value = true;
}

async function saveMembers() {
  if (!selectedCustomGroup.value) return;
  const previous = customGroups.value;
  customGroups.value = updateDictionaryGroupMembers(
    customGroups.value,
    selectedCustomGroup.value.id,
    memberDraft.value,
    dictsList.value.map((dict) => dict.id),
  );
  saving.value = true;
  try {
    await persistGroups();
    membersVisible.value = false;
    message.success('词典组成员已更新');
  } catch (cause) {
    customGroups.value = previous;
    message.error(cause instanceof Error ? cause.message : '保存词典组成员失败');
  } finally {
    saving.value = false;
  }
}

onMounted(async () => {
  try {
    const [dictionaries, preferences] = await Promise.all([GetAllDicts(), getPreferences()]);
    dictsList.value = dictionaries.map((dict: any) => ({
      id: dict.id,
      name: dict.name,
      dictType: dict.type || '',
      title: dict.description?.title || '',
      generateEngineVersion: dict.description?.generateEngineVersion || '',
      createDate: dict.description?.createDate || '',
      baseDir: dict.base_dir || '',
    }));
    customGroups.value = parseDictionaryGroups(preferences[GROUPS_PREFERENCE_KEY], dictsList.value.map((dict) => dict.id));
    const storedFavorite = String(preferences[FAVORITE_GROUP_PREFERENCE_KEY] || '');
    favoriteGroupId.value = storedFavorite === DEFAULT_DICTIONARY_GROUP_ID || customGroups.value.some((group) => group.id === storedFavorite)
      ? storedFavorite
      : '';
    selectedGroupId.value = favoriteGroupId.value || DEFAULT_DICTIONARY_GROUP_ID;
  } catch (cause) {
    message.error(cause instanceof Error ? cause.message : '加载词典列表失败');
  }
});
</script>

<style lang="scss" scoped>
@use '@/style/variables.scss' as *;

.x-space,
.x-layout { width: 100%; height: 100%; padding: 0; margin: 0; }

.x-layout-main-area {
  display: flex;
  width: 100%;
  height: calc(100% - $layout-footer-height);
}

.x-layout-sidebar { width: $layout-left-sidebar-width; height: 100%; background: var(--c-gray-100); }
.x-layout-content { width: calc(100% - $layout-left-sidebar-width); height: 100%; min-width: 0; }
.n-layout-footer { width: 100%; height: $layout-footer-height; }

.dictionary-content-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-height: 42px;
  padding: 6px 12px;
  box-sizing: border-box;
  border-bottom: 1px solid var(--c-gray-300);
  background: var(--c-gray-50);
}

.dictionary-content-toolbar > div { display: flex; align-items: baseline; gap: 8px; }
.dictionary-content-toolbar strong { color: var(--c-gray-900); font-size: 13px; }
.dictionary-content-toolbar span { color: var(--c-gray-500); font-size: 11px; }
.dictionary-content-toolbar :deep(.n-button), :deep(.n-button) { background-image: none; box-shadow: none; font-size: 12px; }

.dict-main-area { width: 100%; height: calc(100% - 42px); overflow: auto; }
.dict-main-area table { width: 100%; }
.dicts-hint { height: 20px; font-size: 12px; line-height: 20px; }

.dictionary-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 240px;
  gap: 7px;
  color: var(--c-gray-500);
  font-size: 12px;
}

.dictionary-empty > .icon { font-size: 24px; color: var(--c-gray-400); }
.dictionary-empty strong { color: var(--c-gray-800); font-size: 13px; }
.member-dialog-hint { margin-bottom: 12px; color: var(--c-gray-600); font-size: 12px; }
.member-checkbox-list { display: grid; max-height: 320px; gap: 8px; overflow-y: auto; }
</style>
