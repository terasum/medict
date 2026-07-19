import { describe, expect, it } from 'vitest';
import { mount } from '@vue/test-utils';
import DictionaryGroupSidebar from './DictionaryGroupSidebar.vue';

const groups = [
  { id: 'default', name: '默认组', count: 3 },
  { id: 'english', name: '英语词典', count: 1 },
];

describe('DictionaryGroupSidebar', () => {
  it('wires create, delete and favorite actions', async () => {
    const wrapper = mount(DictionaryGroupSidebar, {
      props: {
        groups,
        selectedGroupId: 'english',
        favoriteGroupId: '',
        defaultGroupId: 'default',
      },
    });

    await wrapper.get('[aria-label="新建词典组"]').trigger('click');
    await wrapper.get('[aria-label="删除当前词典组"]').trigger('click');
    await wrapper.get('[aria-label="设为常用词典组"]').trigger('click');

    expect(wrapper.emitted('create')).toHaveLength(1);
    expect(wrapper.emitted('delete')).toHaveLength(1);
    expect(wrapper.emitted('toggle-favorite')).toHaveLength(1);
  });

  it('protects the default group from deletion and selects groups', async () => {
    const wrapper = mount(DictionaryGroupSidebar, {
      props: {
        groups,
        selectedGroupId: 'default',
        favoriteGroupId: 'default',
        defaultGroupId: 'default',
      },
    });

    expect(wrapper.get('[aria-label="删除当前词典组"]').attributes('disabled')).toBeDefined();
    expect(wrapper.get('[aria-label="取消常用词典组"]').classes()).toContain('active');
    await wrapper.findAll('.dictionary-group-item')[1].trigger('click');
    expect(wrapper.emitted('select')).toEqual([['english']]);
  });
});
