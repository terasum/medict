<template>
  <section class="settings-page">
    <header class="settings-page-header">
      <button v-if="back" class="settings-back-button" type="button" aria-label="返回上一页" @click="goBack">
        <span class="icon icon-left-open" aria-hidden="true" />
        返回
      </button>
      <div>
        <h1>{{ title }}</h1>
        <p v-if="description">{{ description }}</p>
      </div>
    </header>
    <div class="settings-page-content">
      <slot />
    </div>
  </section>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router';

withDefaults(defineProps<{ title: string; description?: string; back?: boolean }>(), {
  description: '',
  back: false,
});

const router = useRouter();

function goBack() {
  if (window.history.length > 1) router.back();
  else router.replace('/setting');
}
</script>
