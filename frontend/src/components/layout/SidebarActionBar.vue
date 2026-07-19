<template>
  <div class="sidebar-action-bar">
    <slot />
  </div>
</template>

<style lang="scss" scoped>
.sidebar-action-bar {
  display: flex;
  align-items: center;
  flex: 0 0 32px;
  height: 32px;
  padding: 3px 6px;
  box-sizing: border-box;
  border-top: 1px solid var(--c-gray-300);
  background: var(--c-gray-200);
}

.sidebar-action-bar :deep(button) {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 24px;
  padding: 0;
  color: var(--c-gray-700);
  background: transparent;
  background-image: none;
  border: 0;
  border-radius: 4px;
  box-shadow: none;
  font: inherit;
  cursor: pointer;
}

.sidebar-action-bar :deep(button:hover:not(:disabled)) {
  background: var(--c-gray-300);
}

.sidebar-action-bar :deep(button.active) {
  color: var(--c-primary);
  background: var(--c-gray-300);
}

.sidebar-action-bar :deep(button:disabled) {
  color: var(--c-gray-400);
  cursor: default;
}

.sidebar-action-bar :deep(button:focus-visible) {
  outline: 1px solid var(--c-primary);
  outline-offset: -1px;
}

.sidebar-action-bar :deep(.icon) {
  margin: 0;
  line-height: 1;
}
</style>
