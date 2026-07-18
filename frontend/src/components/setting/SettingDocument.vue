<template>
  <SettingPage :title="title" :description="description" back>
    <article class="settings-document markdown-body">
      <component :is="document" />
    </article>
  </SettingPage>
</template>

<script setup lang="ts">
import type { Component } from 'vue';
import SettingPage from './SettingPage.vue';

defineProps<{ title: string; description?: string; document: Component }>();
</script>
