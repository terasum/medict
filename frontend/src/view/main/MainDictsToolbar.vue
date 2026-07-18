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
@import '@/style/photon/photon.scss';

.dictionaries {
  display: flex;
  flex-direction: row;
  height: 100%;
  .dictionary-item {
    margin: 4px 3px;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 22px;
    height: 22px;
    padding: 0;
    border: 0;
    text-align: center;
    border-radius: 3px;
    cursor: pointer;
    user-select: none;
    -webkit-user-select: none;
    color: var(--c-gray-700);
    background-color: var(--c-gray-50);
    // box-shadow: rgba(0, 0, 0, 0.16) 0px 3px 6px, rgba(0, 0, 0, 0.23) 0px 3px 6px;
    box-shadow: rgba(0, 0, 0, 0.1) 0px 4px 6px -1px,
      rgba(0, 0, 0, 0.06) 0px 2px 4px -1px;

    &:hover {
      background-color: var(--c-gray-200);
    }

    &:focus-visible {
      outline: 2px solid var(--c-primary);
      outline-offset: 1px;
    }

    .dictionary-fallback-icon {
      width: 14px;
      height: 14px;
      font-size: 14px;
    }
  }
  .dictionary-item-active {
    color: var(--c-primary);
    box-shadow:
      inset 0 0 0 2px var(--c-primary),
      rgba(0, 0, 0, 0.1) 0 2px 4px;
  }
}
</style>
<template>
  <div class="dictionaries">
    <button
      v-for="item in state.dictList"
      type="button"
      class="dictionary-item"
      :class="{ 'dictionary-item-active': isSelected(item) }"
      :key="item.id"
      @click="chooseDict(item)"
      :style="getBackground(item)"
      :title="item.name"
      :aria-label="`选择词典 ${item.name}`"
      :aria-pressed="isSelected(item)"
    >
      <n-icon
        v-if="!item.background"
        class="dictionary-fallback-icon"
        :component="getFallbackIcon(item)"
      />
    </button>
  </div>
</template>
<script setup>
import { useDictQueryStore } from '@/store/dict';
import { useUIStore } from '@/store/ui';
import { reactive, onMounted } from 'vue';
import { BuildIndex } from '@/apis/dicts-api';
import { NIcon } from 'naive-ui';
import { BookOpen, Database, Globe, Language } from '@vicons/fa';
import { dictionaryIconKind, isDictionarySelected } from './dictionary-toolbar';

const dictQueryStore = useDictQueryStore();
const uiStore = useUIStore();

const state = reactive({
  dictList: [],
});

function chooseDict(item) {
  dictQueryStore.updateSelectDict(item);
}

function isSelected(item) {
  return isDictionarySelected(item.id, dictQueryStore.selectDict.id);
}

function getBackground(item) {
  if (item.background) {
    return {
      backgroundImage: `url(${item.background})`,
      backgroundSize: 'cover',
      backgroundRepeat: 'no-repeat',
      backgroundPosition: 'center',
    };
  }
  return undefined;
}

const fallbackIcons = {
  database: Database,
  online: Globe,
  language: Language,
  book: BookOpen,
};

function getFallbackIcon(item) {
  return fallbackIcons[dictionaryIconKind(item.dict_type)];
}

function loadDictionaries() {
  dictQueryStore.queryDictList().then((res) => {
    console.log('[app-init] loading dictionaries success, result:', res);
    if (res.length > 0) {
      dictQueryStore.updateSelectDict(res[0]);
    }
    const totalNumber = res.length;
    const updater = updateProgress(totalNumber);

    function sequenceHandle(promiseArr) {
      const pro = promiseArr.shift();
      if (pro && pro.handle) {
        pro.handle().then((resp) => {
          pro.callback(resp);
          sequenceHandle(promiseArr);
        });
      }
    }

    function buildIndexPromise(i, id, name) {
      return {
        handle: function () {
          return BuildIndex(id);
        },
        callback: function (resp) {
          let progressHint = `词典 ${name} 加载完成`;
          console.log(`[app-init] building success, index: ${i}`, resp);
          updater(progressHint);
        },
      };
    }

    let promiseArray = [];
    for (let i = 0; i < res.length; i++) {
      state.dictList.push(res[i]);
      promiseArray.push(buildIndexPromise(i, res[i].id, res[i].name));
    }
    sequenceHandle(promiseArray);
  });
}

function updateProgress(totalNumber) {
  let total = totalNumber;
  let count = 0;
  let plist = [];
  let progress = 0;

  let intv = setInterval(() => {
    if (plist.length > 0) {
      let item = plist.shift();
      count++;
      progress = (count / total) * 100;
      uiStore.updateProgress(item.hint, progress);
      if (progress >= 100) {
        clearInterval(intv);
        setTimeout(() => {
          uiStore.updateProgress('全部加载完成', 100);
        }, 150);
      }
    }
  }, 200);

  return function (hint) {
    plist.push({
      hint: hint,
    });
  };
}

onMounted(() => {
  dictQueryStore.initDicts().then((res) => {
    console.log('[dict-api] init dicts success, result:', res);
    loadDictionaries();
  });
});
</script>
