<template>
  <div class="setting-group">
    <div class="setting-item">
      <div class="setting-item-label">
        <span class="setting-item-title">{{ title }}</span>
        <span class="setting-item-desc"><slot name="desc"></slot> </span>
      </div>

      <span class="setting-item-action">
        <slot name="action"></slot>
      </span>
    </div>
    <div class="setting-item-content">
      <slot></slot>
    </div>
  </div>
</template>
<script setup lang="ts">
defineProps(['title', 'value']);

// Type the named slots so consumers' #desc / #action are recognized (vue-tsc, #702).
defineSlots<{
	desc(): void;
	action(): void;
	default(): void;
}>();
</script>
<style lang="scss" scoped>
.setting-group {
  display: grid;
  grid-template-columns: minmax(180px, 0.8fr) minmax(220px, 1.2fr);
  gap: 20px;
  border-bottom: 1px solid var(--c-gray-300);
  &:nth-last-child(1) {
    border-bottom: none;
  }
  padding: 14px 16px;

  .setting-item {
    min-width: 0;
    .setting-item-label {
      display: flex;
      flex-direction: column;
      gap: 4px;

      .setting-item-title {
        color: var(--c-gray-900);
        font-size: 13px;
        font-weight: 600;
      }
      .setting-item-desc {
        color: var(--c-gray-600);
        font-size: 12px;
        line-height: 1.45;
      }
    }
  }

  .setting-item-content {
    min-width: 0;
    color: var(--c-gray-800);
    font-size: 13px;
  }

  .setting-item-action {
    display: none;
  }
}

@media (max-width: 720px) {
  .setting-group {
    grid-template-columns: 1fr;
    gap: 8px;
  }
}
</style>
