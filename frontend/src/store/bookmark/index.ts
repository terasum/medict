/**
 *
 * Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

import { defineStore } from 'pinia';
import {
  getBookmarks,
  addBookmark,
  removeBookmark,
  getNotebooks,
  createNotebook,
  renameNotebook,
  deleteNotebook,
  setDefaultNotebook,
  type Bookmark,
  type Notebook,
} from '@/apis/bookmark-api';

// In-flight load promise so concurrent callers (bookmarks page + star popover)
// share one fetch instead of racing. Mirrors the module-level request id used
// in store/dict.
let loadPromise: Promise<void> | null = null;

/**
 * 生词本 store：笔记本列表 + 当前选中本 + 全部生词。被生词页与星标浮层共用，
 * 保证笔记本列表单源。增删改后按需重载以保持「默认本置顶」的排序与词数。
 */
export const useBookmarkStore = defineStore('bookmark', {
  state: () => ({
    notebooks: [] as Notebook[],
    selectedNotebookId: '',
    bookmarks: [] as Bookmark[],
    loaded: false,
  }),
  getters: {
    defaultNotebook: (state): Notebook | null =>
      state.notebooks.find((n) => n.is_default) || null,
    // notebook_id → 词数
    counts: (state): Record<string, number> => {
      const m: Record<string, number> = {};
      for (const b of state.bookmarks) {
        m[b.notebook_id] = (m[b.notebook_id] || 0) + 1;
      }
      return m;
    },
    // 当前选中本的生词（getBookmarks 已按时间倒序，filter 保序）
    currentBookmarks: (state): Bookmark[] =>
      state.bookmarks.filter((b) => b.notebook_id === state.selectedNotebookId),
  },
  actions: {
    async loadNotebooks() {
      this.notebooks = await getNotebooks();
      // 选中本失效时回落到默认本
      const exists = this.notebooks.some((n) => n.id === this.selectedNotebookId);
      if (!exists) {
        const def = this.notebooks.find((n) => n.is_default);
        this.selectedNotebookId = def ? def.id : this.notebooks[0]?.id ?? '';
      }
    },
    async loadBookmarks() {
      this.bookmarks = await getBookmarks();
    },
    // 首次懒加载：并发拉笔记本与生词，重复调用共享同一个 promise。
    ensureLoaded(): Promise<void> {
      if (this.loaded) return Promise.resolve();
      if (loadPromise) return loadPromise;
      loadPromise = (async () => {
        try {
          await Promise.all([this.loadNotebooks(), this.loadBookmarks()]);
          this.loaded = true;
        } finally {
          loadPromise = null;
        }
      })();
      return loadPromise;
    },
    async createNotebook(name: string) {
      await createNotebook(name);
      await this.loadNotebooks();
    },
    async renameNotebook(id: string, name: string) {
      await renameNotebook(id, name);
      const n = this.notebooks.find((x) => x.id === id);
      if (n) n.name = name;
    },
    async setDefault(id: string) {
      await setDefaultNotebook(id);
      await this.loadNotebooks();
    },
    async deleteNotebook(id: string) {
      await deleteNotebook(id);
      // 后端把单词并入默认本 → 笔记本与生词都需重载
      await Promise.all([this.loadNotebooks(), this.loadBookmarks()]);
    },
    async addBookmark(word: string, dictId: string, notebookId: string) {
      await addBookmark(word, dictId, notebookId);
      await this.loadBookmarks();
    },
    async removeBookmark(word: string, dictId: string, notebookId: string) {
      await removeBookmark(word, dictId, notebookId);
      await this.loadBookmarks();
    },
  },
});
