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
  <div class="parent-window" style="height: 100%">
    <!-- content goes inside .window-content -->
    <div class="window-content">settings</div>
  </div>
</template>

<script lang="ts">
import "../assets/css/photon.min.css"
import { defineComponent } from 'vue';
export default defineComponent({


  components: {},
  methods: {},
});
//  <style scoped lang="css" src="../assets/css/photon.min.css"></style>
</script>

<style lang="scss" scoped>
.parent-window {
  height: 100%;
  .window-content {
    height: 100%;
  }
}
</style>
