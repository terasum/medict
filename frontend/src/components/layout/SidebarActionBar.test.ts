import { describe, expect, it } from 'vitest';
import { mount } from '@vue/test-utils';
import SidebarActionBar from './SidebarActionBar.vue';

describe('SidebarActionBar', () => {
  it('provides one shared toolbar wrapper for sidebar actions', () => {
    const wrapper = mount(SidebarActionBar, {
      slots: { default: '<button aria-label="操作"><span class="icon icon-plus" /></button>' },
    });

    expect(wrapper.classes()).toContain('sidebar-action-bar');
    expect(wrapper.get('button').attributes('aria-label')).toBe('操作');
  });
});
