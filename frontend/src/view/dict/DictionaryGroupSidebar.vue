<template>
  <div class="dictionary-group-sidebar">
    <nav class="dictionary-group-list nav-group" aria-label="词典组">
      <h5 class="nav-group-title">词典组</h5>
      <button
        v-for="group in groups"
        :key="group.id"
        type="button"
        class="dictionary-group-item nav-group-item"
        :class="{ active: group.id === selectedGroupId }"
        :aria-pressed="group.id === selectedGroupId"
        @click="$emit('select', group.id)"
      >
        <span class="icon icon-archive" aria-hidden="true" />
        <span class="dictionary-group-name" :title="group.name">{{ group.name }}</span>
        <span v-if="group.id === favoriteGroupId" class="icon icon-star favorite-icon" aria-label="常用组" />
        <span class="dictionary-group-count">{{ group.count }}</span>
      </button>
    </nav>

    <SidebarActionBar class="dictionary-group-actions" aria-label="词典组操作">
      <button type="button" aria-label="新建词典组" title="新建词典组" @click="$emit('create')">
        <span class="icon icon-plus" aria-hidden="true" />
      </button>
      <button
        type="button"
        aria-label="删除当前词典组"
        title="删除当前词典组"
        :disabled="selectedGroupId === defaultGroupId"
        @click="$emit('delete')"
      >
        <span class="icon icon-minus" aria-hidden="true" />
      </button>
      <button
        type="button"
        :aria-label="selectedGroupId === favoriteGroupId ? '取消常用词典组' : '设为常用词典组'"
        :title="selectedGroupId === favoriteGroupId ? '取消常用词典组' : '设为常用词典组'"
        :class="{ active: selectedGroupId === favoriteGroupId }"
        @click="$emit('toggle-favorite')"
      >
        <span :class="['icon', selectedGroupId === favoriteGroupId ? 'icon-star' : 'icon-star-empty']" aria-hidden="true" />
      </button>
    </SidebarActionBar>
  </div>
</template>

<script setup lang="ts">
import SidebarActionBar from '@/components/layout/SidebarActionBar.vue';

interface GroupListItem {
  id: string;
  name: string;
  count: number;
}

defineProps<{
  groups: GroupListItem[];
  selectedGroupId: string;
  favoriteGroupId: string;
  defaultGroupId: string;
}>();

defineEmits<{
  select: [id: string];
  create: [];
  delete: [];
  'toggle-favorite': [];
}>();
</script>

<style lang="scss" scoped>
.dictionary-group-sidebar {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--c-gray-100);
}

.dictionary-group-list {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 0;
}

.dictionary-group-item {
  display: flex;
  align-items: center;
  width: 100%;
  min-height: 27px;
  gap: 5px;
  padding: 3px 8px 3px 10px;
  color: var(--c-gray-700);
  background: transparent;
  border: 0;
  border-radius: 0;
  font: inherit;
  font-size: 12px;
  text-align: left;
  cursor: pointer;
}

.dictionary-group-item:hover { background: var(--c-gray-200); }
.dictionary-group-item.active { color: var(--c-gray-900); background: var(--c-gray-300); }
.dictionary-group-item:focus-visible { outline: 1px solid var(--c-primary); outline-offset: -1px; }
.dictionary-group-item > .icon:first-child { flex: 0 0 16px; width: 16px; text-align: center; }

.dictionary-group-name {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.favorite-icon { flex: 0 0 12px; color: var(--c-gray-600); font-size: 11px; }
.dictionary-group-count { flex: 0 0 auto; color: var(--c-gray-500); font-size: 11px; }

</style>
