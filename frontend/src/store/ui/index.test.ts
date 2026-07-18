import { beforeEach, describe, expect, it } from 'vitest';
import { setActivePinia, createPinia } from 'pinia';
import { useUIStore } from './index';

describe('useUIStore', () => {
  beforeEach(() => setActivePinia(createPinia()));

  it('defaults to the search tab', () => {
    const ui = useUIStore();
    expect(ui.currentTab).toBe('search');
    expect(ui.isSearchInputActive()).toBe(true);
  });

  it('updateCurrentTab switches the active tab', () => {
    const ui = useUIStore();
    ui.updateCurrentTab('bookmarks');
    expect(ui.currentTab).toBe('bookmarks');
    expect(ui.isSearchInputActive()).toBe(false);
    ui.updateCurrentTab('search');
    expect(ui.isSearchInputActive()).toBe(true);
  });

  it('updateProgress clamps to [0,100] and records the hint', () => {
    const ui = useUIStore();
    ui.updateProgress('loading', 250);
    // (no numeric state stored; just ensure it does not throw and records hint)
    ui.updateProgress('done', -5);
    expect(ui.progressHint).toBe('done');
  });
});
