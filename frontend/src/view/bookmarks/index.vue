<template>
  <div class="bookmarks-page">
    <div class="bookmarks-header">
      <h2>生词本</h2>
      <n-input v-model:value="filter" placeholder="搜索生词..." size="small" style="width: 200px;">
        <template #suffix>
          <n-icon :component="Search" />
        </template>
      </n-input>
    </div>
    <div class="bookmarks-list" v-if="filtered.length > 0">
      <div
        v-for="(item, idx) in filtered"
        :key="idx"
        class="bookmark-item"
        @click="lookupWord(item)"
      >
        <span class="bookmark-word">{{ item.word }}</span>
        <span class="bookmark-dict">{{ item.dict_name }}</span>
        <span class="bookmark-time">{{ formatTime(item.saved_at) }}</span>
        <n-button quaternary size="tiny" @click.stop="removeItem(item)">
          <n-icon><Times /></n-icon>
        </n-button>
      </div>
    </div>
    <div class="bookmarks-empty" v-else>
      <p>暂无生词。在搜索时点击星标按钮收藏单词。</p>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { NInput, NButton, NIcon } from 'naive-ui';
import { Search, Times } from '@vicons/fa';
import { getBookmarks, removeBookmark, type Bookmark } from '@/apis/bookmark-api';
import { useDictQueryStore } from '@/store/dict';
import { useUIStore } from '@/store/ui';

const router = useRouter();
const dictQueryStore = useDictQueryStore();
const uiStore = useUIStore();

const bookmarks = ref<Bookmark[]>([]);
const filter = ref('');

const filtered = computed(() => {
  if (!filter.value.trim()) return bookmarks.value;
  const q = filter.value.toLowerCase();
  return bookmarks.value.filter(
    (b) => b.word.toLowerCase().includes(q) || b.dict_name.toLowerCase().includes(q)
  );
});

function formatTime(ts: number): string {
  const d = new Date(ts * 1000);
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;
}

async function loadBookmarks() {
  bookmarks.value = await getBookmarks();
}

async function removeItem(item: Bookmark) {
  await removeBookmark(item.word, item.dict_id);
  await loadBookmarks();
}

function lookupWord(item: Bookmark) {
  // switch to search tab and look up the word
  uiStore.updateCurrentTab('search');
  router.push('/');
  dictQueryStore.updateInputSearchWord(item.word);
  dictQueryStore.searchWord(item.word);
}

onMounted(() => {
  loadBookmarks();
});
</script>

<style lang="scss" scoped>
.bookmarks-page {
  padding: 16px 20px;
  height: 100%;
  overflow-y: auto;

  .bookmarks-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 16px;

    h2 {
      margin: 0;
      font-size: 18px;
    }
  }

  .bookmarks-list {
    .bookmark-item {
      display: flex;
      align-items: center;
      padding: 8px 12px;
      border-radius: 6px;
      cursor: pointer;
      transition: background 0.15s;
      gap: 12px;

      &:hover {
        background: rgba(128, 128, 128, 0.08);
      }

      .bookmark-word {
        font-weight: 500;
        min-width: 120px;
      }
      .bookmark-dict {
        color: #888;
        font-size: 13px;
        flex: 1;
      }
      .bookmark-time {
        color: #aaa;
        font-size: 12px;
      }
    }
  }

  .bookmarks-empty {
    text-align: center;
    color: #999;
    padding: 40px;
  }
}
</style>
