import { beforeEach, describe, expect, it, vi } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';

// Mock the wailsjs-backed API layer so the store can be exercised in node.
vi.mock('@/apis/bookmark-api', () => ({
  getBookmarks: vi.fn(async () => []),
  addBookmark: vi.fn(async () => {}),
  removeBookmark: vi.fn(async () => {}),
  getNotebooks: vi.fn(async () => []),
  createNotebook: vi.fn(async () => ({ id: 'nb-new', name: 'x', is_default: false, created_at: 0 })),
  renameNotebook: vi.fn(async () => {}),
  deleteNotebook: vi.fn(async () => {}),
  setDefaultNotebook: vi.fn(async () => {}),
}));

import { useBookmarkStore } from './index';

describe('useBookmarkStore (getters)', () => {
  beforeEach(() => setActivePinia(createPinia()));

  it('counts bookmarks per notebook', () => {
    const s = useBookmarkStore();
    s.notebooks = [
      { id: 'a', name: 'A', is_default: false, created_at: 0 },
      { id: 'b', name: 'B', is_default: true, created_at: 0 },
    ];
    s.bookmarks = [
      { word: 'x', dict_id: 'd1', dict_name: 'D', notebook_id: 'a', saved_at: 1 },
      { word: 'y', dict_id: 'd1', dict_name: 'D', notebook_id: 'b', saved_at: 2 },
      { word: 'z', dict_id: 'd2', dict_name: 'D', notebook_id: 'a', saved_at: 3 },
    ];
    expect(s.counts).toEqual({ a: 2, b: 1 });
  });

  it('currentBookmarks filters by selectedNotebookId, preserving order', () => {
    const s = useBookmarkStore();
    s.bookmarks = [
      { word: 'x', dict_id: 'd', dict_name: 'D', notebook_id: 'a', saved_at: 1 },
      { word: 'y', dict_id: 'd', dict_name: 'D', notebook_id: 'b', saved_at: 2 },
      { word: 'z', dict_id: 'd', dict_name: 'D', notebook_id: 'a', saved_at: 3 },
    ];
    s.selectedNotebookId = 'a';
    expect(s.currentBookmarks.map((b) => b.word)).toEqual(['x', 'z']);
  });

  it('defaultNotebook returns the is_default one', () => {
    const s = useBookmarkStore();
    s.notebooks = [
      { id: 'a', name: 'A', is_default: false, created_at: 0 },
      { id: 'b', name: 'B', is_default: true, created_at: 0 },
    ];
    expect(s.defaultNotebook?.id).toBe('b');
    s.notebooks = [];
    expect(s.defaultNotebook).toBeNull();
  });
});
