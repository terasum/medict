<template>
  <div class="setting-group">
    <div class="setting-item">
      <div class="setting-item-value">
        <span class="setting-item-label">{{ title }}</span>
        <span class="setting-item-data">{{ value }} </span>
      </div>
      <span class="setting-item-desc">
        <slot></slot>
      </span>
    </div>
    <span class="setting-item-action">
      <slot name="action"></slot>
    </span>
  </div>
</template>
<script setup>
import { defineProps } from 'vue';

defineProps(['title', 'value']);
</script>
<style lang="scss" scoped>
.setting-group {
  display: flex;
  justify-content: space-between;
  margin: 8 0;
  border-bottom: 1px solid #ccc;
  &:nth-last-child(1) {
    border-bottom: none;
  }
  padding: 10px 0;

  .setting-item {
    max-width: calc(100% - 40px);
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    .setting-item-value {
      display: flex;
      flex-direction: row;
      justify-content: space-between;
      padding-right: 30px;

      .setting-item-label {
        display: flex;
        flex-direction: column;
        min-width: 80px;
        font-weight: bold;
        justify-content: center;
      }
      .setting-item-data {
        color: #666;
        line-break: loose;
        padding-left: 20px;;
      }
    }
    .setting-item-desc {
      color: #999;
      font-size: 13px;
    }
  }
  .setting-item-action {
    min-width: 40px;
  }
}
</style>
