<style lang="scss">
@import '@/style/variables.scss';
    .header-search-box {
      display: flex;
      align-items: center;
      height: 100%;
      width: 100%;
      padding: 0;
      margin: 0;

      .header-navigate-btns {
        height: 100%;
        display: flex;
        align-items: center;
        padding: 0;
        margin: 0 4px 0 0;
        gap: 2px;

        .btn-nav {
          height: 28px;
          width: 28px;
          padding: 0;
          text-align: center;
          font-size: 13px;
          color: #555;
          outline: none;
          border: none;
          background-color: transparent;
          border-radius: 6px;
          cursor: pointer;
          display: flex;
          align-items: center;
          justify-content: center;

          &:hover {
            background-color: rgba(0, 0, 0, 0.06);
          }
          &:active {
            background-color: rgba(0, 0, 0, 0.1);
          }
        }
      }
      .header-search-input {
        flex: 1;
        height: 32px;
        display: flex;
        align-items: center;

        .n-input {
          height: 32px;
          padding: 0 10px;
          margin: 0;
          box-shadow: none;
          font-size: 14px;
          border: 1px solid #e0e0e0;
          border-radius: 8px;
          background-color: #fff;

          &:active {
            outline: none;
            border-color: #bbb;
          }
          &:focus {
            outline: none;
            border-color: #999;
          }
        }
      }
    }

</style>
<template>
   <div class="header-search-box">
        <div class="header-navigate-btns">
          <button
            type="button"
            class="button btn btn-light btn-nav btn-nav-left"
            @click="backHistory()"
          >
            <n-icon><AngleLeft /></n-icon>
          </button>

          <button
            type="button"
            class="button btn btn-light btn-nav btn-nav-right"
            @click="forwardHistory()"
          >
            <n-icon><AngleRight /></n-icon>
          </button>

          <button
            type="button"
            class="button btn btn-light btn-nav"
            @click="toggleBookmark()"
            title="收藏"
          >
            <n-icon><Star /></n-icon>
          </button>
        </div>
        <div class="header-search-input">
          <n-input
            type="text"
            size="small"
            placeholder="搜索"
            @keydown.enter="handleChange"
            v-model:value="inputWord"
          >
            <template #suffix>
              <n-icon :component="Search" />
            </template>
          </n-input>
        </div>
      </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { NIcon } from 'naive-ui';
import { Search, AngleLeft, AngleRight, Star } from '@vicons/fa';
import AppFunctions from '@/components/layout/AppFunctions.vue';

import { useDictQueryStore } from '@/store/dict';
import { addBookmark } from '@/apis/bookmark-api';
import { useUIStore } from '@/store/ui';
import { useRouter } from "vue-router";

const dictQueryStore = useDictQueryStore();
const uiStore = useUIStore();
const router = useRouter();

let inputWord = ref('');
let inputActive = ref(false);


function backHistory() {
  dictQueryStore.backHistory();
}

function forwardHistory() {
  dictQueryStore.forwardHistory();
}

async function toggleBookmark() {
  const word = dictQueryStore.inputSearchWord;
  const dict = dictQueryStore.selectDict;
  if (!word || !dict.id) return;
  await addBookmark(word, dict.id, dict.name || '');
}

let storeChangeUnscribe = null;
function listenInputWordUpdate() {
  storeChangeUnscribe = dictQueryStore.$onAction(({name, store, after}) => {
      after((result) => {
        switch (name) {
          case 'updateInputSearchWord': {
            // inputWord.value = dictQueryStore.inputSearchWord;
            break;
          }
          case 'forwardHistory': {
            inputWord.value = dictQueryStore.inputSearchWord;
            break;
          }
          case 'backHistory': {
            inputWord.value = dictQueryStore.inputSearchWord;
            break;
          }
        }
      });
    }
  );
}

onMounted(() => {
  if (storeChangeUnscribe) {
    storeChangeUnscribe();
    storeChangeUnscribe = null;
  }
  listenInputWordUpdate();
})

///----------------------------
// event listener function
///----------------------------

function handleChange(v) {
  console.info('[app-event](keydown.enter), args:' + inputWord.value);
  if(!uiStore.isSearchInputActive()){
    console.log("[app-event](keydown.enter), input disabled, skipped")
    return;
  }
  let word = inputWord.value.trim();

  dictQueryStore.updateInputSearchWord(word);
  dictQueryStore.searchWord(word);
}
</script>
