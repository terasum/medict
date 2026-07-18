<template>
  <div class="x-space">
    <div class="x-layout">
      <div class="x-layout-header">
        <AppHeader />
      </div>
      <div class="x-layout-main-area">
        <aside class="bm-sidebar">
          <NotebookSidebar />
        </aside>
        <main class="bm-content">
          <div class="bm-content-head">
            <div class="bm-title">
              <h2>{{ currentNotebookName }}</h2>
              <span class="bm-count">{{ store.currentBookmarks.length }} 词</span>
            </div>
            <div class="bm-actions">
              <n-input v-model:value="filter" placeholder="搜索生词..." size="small">
                <template #suffix>
                  <n-icon :component="Search" />
                </template>
              </n-input>
              <button
                type="button"
                class="btn btn-default bm-export"
                :disabled="exporting || store.currentBookmarks.length === 0"
                @click="onExportAnki"
              >
                <n-icon><Download /></n-icon>
                {{ exporting ? '导出中…' : '导出 Anki' }}
              </button>
            </div>
          </div>
          <div class="bm-list" v-if="filtered.length > 0">
            <table class="table-striped bm-table">
              <colgroup>
                <col class="col-word" />
                <col />
                <col class="col-time" />
                <col class="col-action" />
              </colgroup>
              <thead>
                <tr>
                  <th>单词</th>
                  <th>来源词典</th>
                  <th>收藏日期</th>
                  <th><span class="visually-hidden">操作</span></th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="item in filtered"
                  :key="item.word + '@' + item.dict_id"
                  class="bm-item"
                  @click="lookupWord(item)"
                >
                  <td class="bm-word">
                    <button type="button" class="bm-lookup" @click.stop="lookupWord(item)">
                      {{ item.word }}
                    </button>
                  </td>
                  <td class="bm-dict">{{ item.dict_name }}</td>
                  <td class="bm-time">{{ formatTime(item.saved_at) }}</td>
                  <td class="bm-action">
                    <button
                      type="button"
                      class="bm-remove"
                      :aria-label="`移除生词 ${item.word}`"
                      @click.stop="removeItem(item)"
                    >
                      <n-icon><Times /></n-icon>
                    </button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="bm-empty" v-else>
            <p>该生词本暂无生词。在搜索时点击星标按钮收藏单词。</p>
          </div>
        </main>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { NInput, NIcon, useMessage } from 'naive-ui';
import { Search, Times, Download } from '@vicons/fa';
import AppHeader from '@/components/layout/AppHeader.vue';
import NotebookSidebar from '@/components/bookmarks/NotebookSidebar.vue';
import { useBookmarkStore } from '@/store/bookmark';
import { useDictQueryStore } from '@/store/dict';
import { useUIStore } from '@/store/ui';
import { getBookmarkSnapshot, exportAnkiToApkg, type Bookmark } from '@/apis/bookmark-api';

const router = useRouter();
const store = useBookmarkStore();
const dictQueryStore = useDictQueryStore();
const uiStore = useUIStore();
uiStore.updateCurrentTab('bookmarks');

const filter = ref('');
const message = useMessage();
const exporting = ref(false);

const filtered = computed(() => {
  const list = store.currentBookmarks;
  if (!filter.value.trim()) return list;
  const q = filter.value.toLowerCase();
  return list.filter(
    (b) => b.word.toLowerCase().includes(q) || b.dict_name.toLowerCase().includes(q)
  );
});

const currentNotebookName = computed(() => {
  const nb = store.notebooks.find((n) => n.id === store.selectedNotebookId);
  return nb ? nb.name : '生词本';
});

function formatTime(ts: number): string {
  const d = new Date(ts * 1000);
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;
}

async function removeItem(item: Bookmark) {
  await store.removeBookmark(item.word, item.dict_id, item.notebook_id);
}

async function onExportAnki() {
  exporting.value = true;
  try {
    const path = await exportAnkiToApkg(store.selectedNotebookId);
    if (path) {
      message.success(`已导出：${path}`);
    } // 用户取消保存对话框时不提示
  } catch (e) {
    message.error((e as Error)?.message || '导出失败');
  } finally {
    exporting.value = false;
  }
}

async function lookupWord(item: Bookmark) {
  uiStore.updateCurrentTab('search');
  router.push('/');
  dictQueryStore.updateInputSearchWord(item.word);
  // 词典仍在 → 切到对应词典活查
  if (item.dict_id) {
    try {
      const list = (await dictQueryStore.queryDictList()) as Array<{ id: string }>;
      const target = (list || []).find((d) => d.id === item.dict_id);
      if (target) {
        // updateSelectDict 切词典后会自动 searchWord(inputSearchWord)
        dictQueryStore.updateSelectDict(target);
        return;
      }
    } catch (e) {
      console.warn('[bookmarks] lookupWord: resolve dict failed', e);
    }
  }
  // 词典已卸载 → 显示收藏时保存的 HTML 快照（资源已内联，自包含）
  try {
    const html = await getBookmarkSnapshot(item.word, item.dict_id, item.notebook_id);
    if (html) {
      dictQueryStore.updateMainContentURL('');
      // 用 unicode 安全的 base64 写入 mainContent（store.updateMainContent 的 btoa
      // 对非 ASCII 的词典 HTML 会抛错，且 iframe 的 b64DecodeUnicode 正是此编码的逆）
      dictQueryStore.mainContent = btoa(unescape(encodeURIComponent(html)));
      return;
    }
  } catch (e) {
    console.warn('[bookmarks] lookupWord: snapshot load failed', e);
  }
  dictQueryStore.searchWord(item.word); // 兜底：当前词典
}

onMounted(() => {
  store.ensureLoaded();
});
</script>

<style lang="scss" scoped>
@import '@/style/variables.scss';
@import '@/style/photon/photon.scss';

.x-space {
  width: 100%;
  height: 100%;
  padding: 0;
  margin: 0;
}

.x-layout {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;

  .x-layout-header {
    flex: 0 0 auto;
  }

  .x-layout-main-area {
    flex: 1 1 auto;
    display: flex;
    flex-direction: row;
    min-height: 0;

    .bm-sidebar {
      width: $layout-left-sidebar-width;
      flex: 0 0 auto;
      height: 100%;
      border-right: 1px solid var(--c-gray-300);
      overflow: hidden;
    }

    .bm-content {
      flex: 1 1 auto;
      height: 100%;
      min-width: 0;
      display: flex;
      flex-direction: column;
      padding: 0;
      overflow: hidden;
      background-color: #fff;

      .visually-hidden {
        position: absolute;
        width: 1px;
        height: 1px;
        padding: 0;
        margin: -1px;
        overflow: hidden;
        clip: rect(0, 0, 0, 0);
        white-space: nowrap;
        border: 0;
      }

      .bm-content-head {
        min-height: 38px;
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 4px 8px 4px 10px;
        margin: 0;
        flex: 0 0 auto;
        background: var(--c-gray-100);
        box-shadow: inset 0 -1px 0 var(--c-gray-300);

        .bm-title {
          display: flex;
          align-items: baseline;
          gap: 8px;
          overflow: hidden;

          h2 {
            margin: 0;
            font-size: 13px;
            font-weight: 600;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }
          .bm-count {
            color: var(--c-gray-500);
            font-size: 12px;
            flex-shrink: 0;
          }
        }

        .bm-actions {
          display: flex;
          align-items: center;
          gap: 6px;

          .n-input {
            width: 200px;
          }

          .bm-export {
            height: 28px;
            display: flex;
            align-items: center;
            gap: 5px;

            &:disabled {
              opacity: 0.5;
              cursor: not-allowed;
            }
          }
        }
      }

      .bm-list {
        flex: 1 1 auto;
        overflow-y: auto;
        background-color: #fff;

        .bm-table {
          table-layout: fixed;

          .col-word {
            width: 24%;
          }
          .col-time {
            width: 110px;
          }
          .col-action {
            width: 34px;
          }

          thead {
            position: sticky;
            top: 0;
            z-index: 1;
          }

          th,
          td {
            height: 28px;
            padding: 3px 10px;
          }

          th {
            color: var(--c-gray-600);
          }
        }

        .bm-item {
          cursor: pointer;

          .bm-word {
            font-weight: 600;
            color: var(--c-gray-900);

            .bm-lookup {
              width: 100%;
              padding: 0;
              border: 0;
              background: transparent;
              color: inherit;
              font: inherit;
              font-weight: inherit;
              text-align: left;
              cursor: pointer;

              &:focus-visible {
                outline: 1px solid var(--c-primary);
                outline-offset: 1px;
              }
            }
          }
          .bm-dict {
            color: var(--c-gray-600);
          }
          .bm-time {
            color: var(--c-gray-500);
          }
          .bm-action {
            padding: 0 5px;
            text-align: center;
          }
          .bm-remove {
            width: 22px;
            height: 22px;
            padding: 0;
            border: 0;
            background: transparent;
            color: var(--c-gray-500);
            display: inline-flex;
            align-items: center;
            justify-content: center;
            opacity: 0;

            &:hover {
              color: var(--c-gray-800);
              background-color: var(--c-gray-300);
            }
          }
          &:hover .bm-remove,
          &:focus-within .bm-remove {
            opacity: 1;
          }
        }
      }

      .bm-empty {
        flex: 1 1 auto;
        display: flex;
        align-items: center;
        justify-content: center;
        color: var(--c-gray-500);
        font-size: 12px;
      }
    }
  }
}
</style>
