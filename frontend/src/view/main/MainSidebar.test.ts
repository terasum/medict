import { mount } from '@vue/test-utils';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const store = vi.hoisted(() => ({
  multiMode: false,
  queryPendingList: [
    { id: 0, keyword: 'they' },
    { id: 0, keyword: "they're" },
  ],
  locateWord: vi.fn(),
  locateInMultiPrimary: vi.fn(),
}));

vi.mock('@/store/dict', () => ({
  useDictQueryStore: () => store,
}));

import MainSidebar from './MainSidebar.vue';

describe('MainSidebar', () => {
  beforeEach(() => {
    store.locateWord.mockClear();
  });

  it('locates by list position when dictionary entry IDs are duplicated', async () => {
    const wrapper = mount(MainSidebar, {
      global: {
        stubs: {
          AppSidebar: { template: '<aside><slot /></aside>' },
        },
      },
    });

    await wrapper.findAll('li')[1].trigger('click');

    expect(store.locateWord).toHaveBeenCalledWith(1);
    wrapper.unmount();
  });
});
