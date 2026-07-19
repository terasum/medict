import { describe, expect, it } from 'vitest';
import { mount } from '@vue/test-utils';
import SettingPage from './SettingPage.vue';

describe('SettingPage', () => {
  it('renders setting content without an internal header or back navigation', () => {
    const wrapper = mount(SettingPage, {
      slots: { default: '<div data-test="content">设置内容</div>' },
    });

    expect(wrapper.find('.settings-page-header').exists()).toBe(false);
    expect(wrapper.find('[aria-label="返回上一页"]').exists()).toBe(false);
    expect(wrapper.find('h1').exists()).toBe(false);
    expect(wrapper.get('[data-test="content"]').text()).toBe('设置内容');
  });
});
