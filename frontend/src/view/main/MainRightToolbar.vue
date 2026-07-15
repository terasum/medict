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

<style lang="scss" scoped>
@import '@/style/variables.scss';

.dictionaries {
  display: flex;
  flex-direction: column;
  width: 60px;
  height: 100%;

  .dictionary-item {
    display: block;
    position: relative;
    width: 32px;
    height: 32px;
    text-align: center;
    line-height: 32px;
    margin: 6px auto;
    border: 1px solid #ccc;
    border-radius: 8px;
    font-size: 13px;
    cursor: pointer;
    user-select: none;
    -webkit-user-select: none;
    // box-shadow: rgba(0, 0, 0, 0.16) 0px 3px 6px, rgba(0, 0, 0, 0.23) 0px 3px 6px;
    box-shadow: rgba(0, 0, 0, 0.1) 0px 4px 6px -1px,
      rgba(0, 0, 0, 0.06) 0px 2px 4px -1px;

    &:hover {
      background-color: #f1f1f1;
    }
  }
  .dictionary-item-active {
    border: 1px solid rgba(17, 168, 255, 0.858);
    width: 36px;
    height: 36px;
    line-height: 36px;
    padding: 2px;
    box-shadow: rgba(0, 0, 0, 0.1) 0px 20px 25px -5px,
      rgba(0, 0, 0, 0.04) 0px 10px 10px -5px;
  }
  // 多词典模式开关
  .multi-toggle {
    width: 32px;
    height: 24px;
    margin: 4px auto 8px;
    border: 1px solid #ccc;
    border-radius: 6px;
    background: #fff;
    color: #666;
    font-size: 12px;
    cursor: pointer;
    user-select: none;
    &:hover {
      background-color: #f1f1f1;
    }
  }
  .multi-toggle-active {
    background: #11a8ff;
    color: #fff;
    border-color: #11a8ff;
  }
  // 多模式下激活词典的角标
  .dict-check {
    position: absolute;
    top: -3px;
    right: -3px;
    width: 14px;
    height: 14px;
    border-radius: 50%;
    background: #11a8ff;
    color: #fff;
    font-size: 10px;
    line-height: 14px;
    text-align: center;
    box-shadow: 0 0 0 1px #fff;
  }
}
</style>
<template>
  <AppRightToolbar>
    <div class="dictionaries">
      <button
        class="multi-toggle"
        :class="{ 'multi-toggle-active': dictQueryStore.multiMode }"
        title="多词典查询(开启后点词典图标多选)"
        @click="dictQueryStore.setMultiMode(!dictQueryStore.multiMode)"
      >多</button>
      <n-popover
        v-for="item in state.dictList"
        :overlap="false"
        placement="left"
        trigger="hover"
      >
        <template #trigger>
          <span
            class="dictionary-item"
            :class="isActive(item) ? 'dictionary-item-active' : ''"
            :key="item.id"
            @click="chooseDict(item)"
            :style="getBackground(item)"
          ><span v-if="dictQueryStore.multiMode && isInMultiSet(item.id)" class="dict-check">✓</span></span>
        </template>
        <div class="large-text">
          <div>
            {{
              item.description && item.description.title
                ? item.description.title
                : item.name
            }}
          </div>
          <!-- <div style="width: 95%; max-width: 340px; max-height: 200px; overflow-y: auto"> -->
            <!-- <div v-html="item.description.description"></div> -->
          <!-- </div> -->
        </div>
      </n-popover>
    </div>
  </AppRightToolbar>
</template>
<script setup>
import { useDictQueryStore } from '@/store/dict';
import { useUIStore } from '@/store/ui';
import { reactive, onMounted } from 'vue';
import { BuildIndex } from '@/apis/dicts-api';
import AppRightToolbar from '@/components/layout/AppRightToolbar.vue';

import { NPopover } from 'naive-ui';

const dictQueryStore = useDictQueryStore();
const uiStore = useUIStore();

const state = reactive({
  dictList: [],
});

function chooseDict(item) {
  // 多词典模式:点图标增/减激活集;否则单词典切换。
  if (dictQueryStore.multiMode) {
    dictQueryStore.toggleMultiDict(item);
  } else {
    dictQueryStore.updateSelectDict(item);
  }
}

// 选中态:多模式下=在激活集里;单模式下=当前 selectDict。
function isActive(item) {
  if (dictQueryStore.multiMode) {
    return dictQueryStore.multiSelectedDicts.some((d) => d.id === item.id);
  }
  return item.id == dictQueryStore.selectDict.id;
}
function isInMultiSet(id) {
  return dictQueryStore.multiSelectedDicts.some((d) => d.id === id);
}

function getBackground(item) {
  if (item.background) {
    let style = `background:url(${item.background});`;
    style += `background-size:cover;`;
    style += `background-repeat:no-repeat;`;
    style += `background-position:center;`;
    style += `color: #fff;`;
    return style;
  }
  return '';
}

function loadDictionaries() {
  dictQueryStore.queryDictList().then((res) => {
    console.log("[app-init] loading dictionaries success, result:", res)
    if (res.length > 0) {
      dictQueryStore.updateSelectDict(res[0]);
    }
    const totalNumber = res.length;
    const updater = updateProgress(totalNumber);

    function sequenceHandle(promiseArr) {
      const pro = promiseArr.shift()
      if(pro && pro.handle) {
        pro.handle().then((resp)=>{
          pro.callback(resp);
          sequenceHandle(promiseArr)
        })
      }
    }

    function buildIndexPromise(i, id, name) {
      return {
        handle:  function() {
        return BuildIndex(id)
      },
      callback: function(resp) {
        let progressHint = `词典 ${name} 加载完成`;
        console.log(`[app-init] building success, index: ${i}`, resp);
        updater(progressHint)
      }
    }
  }

    let promiseArray = [];
    for (let i = 0; i < res.length; i++) {
      state.dictList.push(res[i]);
      promiseArray.push(buildIndexPromise(i, res[i].id, res[i].name))
    }
    sequenceHandle(promiseArray);
  });
}


function updateProgress(totalNumber) {
  let total = totalNumber;
  let count = 0;
  let plist = []
  let progress = 0;

  let intv = setInterval(() => {
    if (plist.length > 0) {
      let item = plist.shift();
      count++;
      progress = (count / total) * 100;
      uiStore.updateProgress(item.hint, progress);
      if (progress >= 100) {
        clearInterval(intv);
        setTimeout(() =>{
          uiStore.updateProgress("全部加载完成", 100);
        },150)
      }
    }
  }, 200);
  
  return function(hint) {
      plist.push({
        hint: hint,
      })
  }
}


onMounted(() => {
  dictQueryStore.initDicts().then((res)=>{
    console.log("[dict-api] init dicts success, result:", res)
    loadDictionaries();
  })
});
</script>
