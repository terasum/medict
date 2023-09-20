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
  <AppFooter>
    <span class="index-process-hint">
      <em>{{ percentage_hint.value }}</em>
    </span>

    <div class="building-index-process">
      <n-progress
        type="line"
        color="#c1c1c1"
        :show-indicator="false"
        :percentage="percnetage.value"
      />
    </div>
  </AppFooter>
</template>

<script lang="ts" setup>
import { NProgress } from 'naive-ui';
import AppFooter from '@/components/layout/AppFooter.vue';
import { ref, onMounted, computed } from 'vue';

let _percentage = ref(0);
let _percentage_hint = ref('索引建立中...');

function updatePercentage() {
  let intv = setInterval(() => {
    let step = 2000 / 200;
    _percentage.value += step;
    if (_percentage.value >= 100) {
      clearInterval(intv);
      _percentage_hint.value = '索引建立完成';
    }
  }, 200);
}

onMounted(()=>{
  updatePercentage();
})

const percnetage = computed(() => {
  return _percentage;
});

const percentage_hint = computed(() =>{
  return _percentage_hint
})

</script>

<style lang="scss" scoped>
.index-process-hint {
  font-size: 12px;
  line-height: 20px;
  margin: 0 5px ;

  font-style: normal;
  color: #666;

}
  
  .building-index-process {
    display: flex;
    justify-content: center;
    flex-direction: column;
    height: 20px;
    width: 120px;
  }
</style>
