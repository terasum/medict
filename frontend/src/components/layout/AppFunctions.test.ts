import { beforeEach, describe, expect, it, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import { createPinia, setActivePinia } from 'pinia';
import { useUIStore } from '@/store/ui';
import AppFunctions from './AppFunctions.vue';

const replace = vi.fn();

vi.mock('vue-router', () => ({
  useRouter: () => ({ replace }),
}));

describe('AppFunctions', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    replace.mockReset();
  });

  it('navigates to settings even when the tab state is already marked active', async () => {
    useUIStore().updateCurrentTab('setting');
    const wrapper = mount(AppFunctions, {
      global: {
        stubs: { Search: true, Book: true, ToggleOn: true, Star: true },
      },
    });

    await wrapper.get('[data-test="function-setting"]').trigger('click');

    expect(replace).toHaveBeenCalledWith({ path: '/setting' });
    expect(useUIStore().currentTab).toBe('setting');
  });
});
