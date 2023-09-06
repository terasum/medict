<template>
  <div id="app">
    <div class="fake-title-bar" data-wails-drag></div>
    <router-view></router-view>
    <h1>APP</h1>
  </div>
</template>

<script>
export default {
  name: "App"
};
</script>