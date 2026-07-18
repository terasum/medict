<!--

 Copyright (C) 2023 Quan Chen <chenquan_act@163.com>

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.

 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.
-->

<template>
  <div class="notebook-sidebar">
    <nav class="nb-list nav-group">
      <h5 class="nav-group-title">生词本</h5>
      <div
        v-for="nb in store.notebooks"
        :key="nb.id"
        class="nb-item nav-group-item"
        :class="{ active: nb.id === store.selectedNotebookId }"
      >
        <button
          type="button"
          class="nb-select"
          :aria-pressed="nb.id === store.selectedNotebookId"
          @click="store.selectedNotebookId = nb.id"
        >
          <n-icon class="nb-icon">
            <Star v-if="nb.is_default" />
            <Book v-else />
          </n-icon>
          <span class="nb-name" :title="nb.name">{{ nb.name }}</span>
          <span class="nb-count">{{ store.counts[nb.id] || 0 }}</span>
        </button>
        <n-dropdown
          trigger="click"
          placement="bottom-end"
          :options="menuOptions(nb)"
          @select="(key: string) => onMenuSelect(key, nb)"
        >
          <button type="button" class="nb-more" :aria-label="`${nb.name}的更多操作`" @click.stop>
            <n-icon><EllipsisH /></n-icon>
          </button>
        </n-dropdown>
      </div>
      <div v-if="store.notebooks.length === 0" class="nb-empty">暂无生词本</div>
    </nav>

    <div class="nb-footer">
      <button type="button" class="btn btn-default nb-create" @click="openCreate">
        <n-icon><Plus /></n-icon>
        新建生词本
      </button>
    </div>

    <!-- 新建 / 重命名对话框 -->
    <n-modal v-model:show="formShow" preset="dialog" :title="formTitle" :show-icon="false">
      <n-input
        ref="formInputRef"
        v-model:value="formValue"
        placeholder="生词本名称"
        @keydown.enter="confirmForm"
      />
      <template #action>
        <n-button size="small" @click="formShow = false">取消</n-button>
        <n-button size="small" type="primary" :loading="formLoading" @click="confirmForm">
          确定
        </n-button>
      </template>
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, nextTick } from 'vue';
import { NIcon, NButton, NDropdown, NInput, NModal, useMessage, useDialog } from 'naive-ui';
import { Book, Star, EllipsisH, Plus } from '@vicons/fa';
import { useBookmarkStore } from '@/store/bookmark';
import type { Notebook } from '@/apis/bookmark-api';

const store = useBookmarkStore();
const message = useMessage();
const dialog = useDialog();

// —— 新建 / 重命名共用一个对话框 ——
const formShow = ref(false);
const formMode = ref<'create' | 'rename'>('create');
const formValue = ref('');
const formTarget = ref<Notebook | null>(null);
const formLoading = ref(false);
const formInputRef = ref<InstanceType<typeof NInput> | null>(null);

const formTitle = computed(() => (formMode.value === 'create' ? '新建生词本' : '重命名生词本'));

function openCreate() {
  formMode.value = 'create';
  formValue.value = '';
  formTarget.value = null;
  formShow.value = true;
  nextTick(() => formInputRef.value?.focus());
}

function openRename(nb: Notebook) {
  formMode.value = 'rename';
  formValue.value = nb.name;
  formTarget.value = nb;
  formShow.value = true;
  nextTick(() => formInputRef.value?.focus());
}

async function confirmForm() {
  const name = formValue.value.trim();
  if (!name) {
    message.warning('请输入生词本名称');
    return;
  }
  formLoading.value = true;
  try {
    if (formMode.value === 'create') {
      await store.createNotebook(name);
      message.success('已创建生词本');
    } else if (formTarget.value) {
      await store.renameNotebook(formTarget.value.id, name);
      message.success('已重命名');
    }
    formShow.value = false;
  } catch (e: unknown) {
    message.error((e as Error)?.message || '操作失败');
  } finally {
    formLoading.value = false;
  }
}

function menuOptions(nb: Notebook) {
  return [
    { label: '重命名', key: 'rename' },
    { label: '设为默认', key: 'default', disabled: nb.is_default },
    { label: '删除', key: 'delete', disabled: nb.is_default },
  ];
}

function onMenuSelect(key: string, nb: Notebook) {
  if (key === 'rename') {
    openRename(nb);
  } else if (key === 'default') {
    doSetDefault(nb);
  } else if (key === 'delete') {
    doDelete(nb);
  }
}

async function doSetDefault(nb: Notebook) {
  try {
    await store.setDefault(nb.id);
    message.success(`「${nb.name}」已设为默认`);
  } catch (e: unknown) {
    message.error((e as Error)?.message || '操作失败');
  }
}

function doDelete(nb: Notebook) {
  dialog.warning({
    title: '删除生词本',
    content: `确定删除「${nb.name}」？其中的生词将移动到默认生词本。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await store.deleteNotebook(nb.id);
        message.success('已删除');
      } catch (e: unknown) {
        message.error((e as Error)?.message || '删除失败');
      }
    },
  });
}
</script>

<style lang="scss" scoped>
@use '@/style/variables.scss' as *;


.notebook-sidebar {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: var(--c-gray-100);

  .nb-list {
    flex: 1;
    overflow-y: auto;
    padding: 0;
  }

  .nb-item {
    display: flex;
    align-items: center;
    gap: 5px;
    min-height: 26px;
    padding: 3px 6px 3px 10px;
    border-radius: 0;
    cursor: pointer;
    font-size: 12px;
    color: var(--c-gray-700);
    user-select: none;

    &:hover {
      background-color: var(--c-gray-200);
    }
    &.active {
      background-color: var(--c-gray-300);
      color: var(--c-gray-900);
    }

    .nb-icon {
      width: 16px;
      font-size: 13px;
      flex-shrink: 0;
      color: var(--c-gray-600);
    }
    .nb-select {
      min-width: 0;
      flex: 1;
      display: flex;
      align-items: center;
      gap: 5px;
      padding: 0;
      border: 0;
      background: transparent;
      color: inherit;
      font: inherit;
      text-align: left;
      cursor: pointer;
    }
    .nb-name {
      flex: 1;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
    .nb-count {
      font-size: 11px;
      color: var(--c-gray-500);
      flex-shrink: 0;
    }
    .nb-more {
      opacity: 0;
      width: 20px;
      height: 20px;
      padding: 0;
      border: 0;
      background: transparent;
      color: var(--c-gray-600);
      display: flex;
      align-items: center;
      justify-content: center;
      flex-shrink: 0;
      cursor: default;

      &:hover {
        background-color: var(--c-gray-300);
      }
    }
    &:hover .nb-more,
    &.active .nb-more,
    .nb-more:focus-visible {
      opacity: 1;
    }

    .nb-select:focus-visible {
      outline: 1px solid var(--c-primary);
      outline-offset: 1px;
    }
  }

  .nb-empty {
    text-align: center;
    color: var(--c-gray-500);
    font-size: 12px;
    padding: 20px 0;
  }

  .nb-footer {
    height: 32px;
    padding: 3px 6px;
    border-top: 1px solid var(--c-gray-300);
    background-color: var(--c-gray-200);

    .nb-create {
      width: 100%;
      height: 25px;
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 5px;
    }
  }
}
</style>
