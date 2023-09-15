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
  <span class="my-icon-plugins">
    <svg
      xmlns="http://www.w3.org/2000/svg"
      p-id="2475"
      :width="width"
      :height="height"
      fill="currentColor"
      t="1634517990995"
      version="1.1"
      class="bi bi-puzzle-fill"
      viewBox="0 0 1024 1024"
    >
     <path d="M388.704 144.768q22.656 0 38.784 15.68t16.128 38.336v217.824q0 22.656-16.128 38.784t-38.784 16.128H170.88q-22.656 0-38.336-16.128t-15.68-38.784V198.784q0-22.656 15.68-38.336t38.336-15.68h217.824z m0 435.616q22.656 0 38.784 15.68t16.128 38.336v218.688q0 22.656-16.128 38.336t-38.784 15.68H170.88q-22.656 0-38.336-15.68t-15.68-38.336v-218.688q0-22.656 15.68-38.336t38.336-15.68h217.824z m436.512 0q22.656 0 38.336 15.68t15.68 38.336v218.688q0 22.656-15.68 38.336t-38.336 15.68h-217.824q-22.656 0-38.784-15.68t-16.128-38.336v-218.688q0-22.656 16.128-38.336t38.784-15.68h217.824z m121.088-339.776q17.44 17.44 17.44 40.96t-17.44 40.064l-162.048 162.048q-17.44 17.44-40.512 17.44t-40.512-17.44l-162.048-162.048q-16.544-16.544-16.544-40.064t16.544-40.96l162.048-162.048q17.44-16.544 40.512-16.544t40.512 16.544z" p-id="1258"></path>
    </svg>
  </span>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
export default defineComponent({
  props: {
    width: {
      type: Number,
      default: 32,
    },
    height: {
      type: Number,
      default: 32,
    },
  },
});
</script>

<style lang="scss" scoped>
</style>