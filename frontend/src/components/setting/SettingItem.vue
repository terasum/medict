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
<script setup>
import { defineProps } from 'vue';

defineProps(['title', 'value']);
</script>
<style lang="scss" scoped>
.setting-group {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  margin: 8 0;
  border-bottom: 1px solid #ccc;
  &:nth-last-child(1) {
    border-bottom: none;
  }
  padding: 10px 0;

  .setting-item {
    min-width: 320px;
    max-width: calc(100% - 40px);
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    .setting-item-label {
      display: flex;
      flex-direction: row;
      padding-right: 30px;
      margin-bottom: 15px;

      .setting-item-title {
        display: flex;
        flex-direction: column;
        min-width: 80px;
        font-weight: bold;
        justify-content: center;
      }
      .setting-item-desc {
        color: #999;
        font-size: 13px;
        line-break: loose;
        padding-left: 20px;
        font-style: italic;
      }
    }
  }
  .setting-item-action {
    // width: 40px;
  }
}
</style>
