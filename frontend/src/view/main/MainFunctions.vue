<style lang="scss">
@import '@/style/variables.scss';
@import '@/style/photon/photon.scss';
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
          height: 44px;
          width: 28px;
          padding: 0;
          text-align: center;
          font-size: 13px;
          color: var(--c-gray-700);
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
          &:disabled {
            color: var(--c-gray-400);
            cursor: not-allowed;
            opacity: 0.5;
            &:hover {
              background-color: transparent;
            }
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
          border: 1px solid var(--c-gray-300);
          border-radius: 8px;
          background-color: #fff;

          &:active {
            outline: none;
            border-color: var(--c-gray-400);
          }
          &:focus {
            outline: none;
            border-color: var(--c-gray-500);
          }
        }
      }

      // 非搜索页时输入框置灰（naive-ui 的 disabled 已处理文字与交互）
      &.is-disabled {
        .header-search-input {
          .n-input {
            background-color: var(--c-gray-100);
            border-color: var(--c-gray-200);
            cursor: not-allowed;
          }
        }
      }
    }

    // 星标弹出的生词本选择浮层（n-popover 会 teleport 到 body，需放在非 scoped 样式里）
    .nb-picker {
      .nb-picker-title {
        font-size: 12px;
        color: var(--c-gray-500);
        padding: 2px 4px 6px;
      }
      .nb-picker-item {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 6px 8px;
        border-radius: 6px;
        cursor: pointer;
        font-size: 13px;
        color: var(--c-gray-800);

        &:hover {
          background-color: rgba(0, 0, 0, 0.06);
        }
        .nb-picker-icon {
          font-size: 14px;
          flex-shrink: 0;
        }
        .nb-picker-name {
          flex: 1;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
        .nb-picker-tag {
          font-size: 10px;
          color: var(--c-gray-500);
          border: 1px solid var(--c-gray-300);
          border-radius: 3px;
          padding: 0 4px;
          flex-shrink: 0;
        }
      }
      .nb-picker-empty {
        text-align: center;
        color: var(--c-gray-500);
        font-size: 12px;
        padding: 12px 0;
      }
    }

</style>
<template>
   <div class="header-search-box" :class="{ 'is-disabled': searchDisabled }">
        <div class="header-navigate-btns">
          <button
            type="button"
            class="button btn btn-light btn-nav btn-nav-left"
            :disabled="searchDisabled"
            @click="backHistory()"
          >
            <n-icon><AngleLeft /></n-icon>
          </button>

          <button
            type="button"
            class="button btn btn-light btn-nav btn-nav-right"
            :disabled="searchDisabled"
            @click="forwardHistory()"
          >
            <n-icon><AngleRight /></n-icon>
          </button>

          <n-popover
            trigger="click"
            v-model:show="bookmarkPopoverShow"
            placement="bottom-start"
            :width="200"
            @update:show="onBookmarkPopoverShow"
          >
            <template #trigger>
              <button
                type="button"
                class="button btn btn-light btn-nav"
                :disabled="searchDisabled"
                title="收藏"
              >
                <n-icon><Star /></n-icon>
              </button>
            </template>
            <div class="nb-picker">
              <div class="nb-picker-title">加入生词本</div>
              <div
                v-for="nb in bookmarkStore.notebooks"
                :key="nb.id"
                class="nb-picker-item"
                @click="addToNotebook(nb)"
              >
                <n-icon class="nb-picker-icon">
                  <Star v-if="nb.is_default" />
                  <Book v-else />
                </n-icon>
                <span class="nb-picker-name">{{ nb.name }}</span>
                <span v-if="nb.is_default" class="nb-picker-tag">默认</span>
              </div>
              <div v-if="bookmarkStore.notebooks.length === 0" class="nb-picker-empty">
                暂无生词本
              </div>
            </div>
          </n-popover>
        </div>
        <div class="header-search-input">
          <n-input
            type="text"
            size="small"
            placeholder="搜索"
            :disabled="searchDisabled"
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
import { ref, computed, onMounted } from 'vue';
import { NIcon, NPopover, useMessage } from 'naive-ui';
import { Search, AngleLeft, AngleRight, Star, Book } from '@vicons/fa';
import AppFunctions from '@/components/layout/AppFunctions.vue';

import { useDictQueryStore } from '@/store/dict';
import { useBookmarkStore } from '@/store/bookmark';
import { useUIStore } from '@/store/ui';
import { useRouter } from "vue-router";

const dictQueryStore = useDictQueryStore();
const bookmarkStore = useBookmarkStore();
const uiStore = useUIStore();
const router = useRouter();
const message = useMessage();

let inputWord = ref('');
let inputActive = ref(false);

// 非搜索页（切换到词典/生词/设置等 FunctionTab）时，搜索框与导航按钮禁用而非隐藏
const searchDisabled = computed(() => !uiStore.isSearchInputActive());


function backHistory() {
  dictQueryStore.backHistory();
}

function forwardHistory() {
  dictQueryStore.forwardHistory();
}

// 星标 → 弹出生词本选择浮层（默认本置顶并标记）
const bookmarkPopoverShow = ref(false);

function onBookmarkPopoverShow(show) {
  if (show) {
    bookmarkStore.ensureLoaded();
  }
}

async function addToNotebook(nb) {
  const word = dictQueryStore.inputSearchWord;
  const dict = dictQueryStore.selectDict;
  if (!word || !dict.id) {
    message.warning('请先输入单词');
    return;
  }
  try {
    await bookmarkStore.addBookmark(word, dict.id, nb.id);
    message.success(`已加入「${nb.name}」`);
    bookmarkPopoverShow.value = false;
  } catch (e) {
    message.error((e && e.message) || '收藏失败');
  }
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
