import { mount } from '@vue/test-utils';
import { describe, expect, it, vi } from 'vitest';

import EditableSettingValue from './EditableSettingValue.vue';

describe('EditableSettingValue', () => {
  it('opens an inline editor and saves the changed value', async () => {
    const save = vi.fn().mockResolvedValue(undefined);
    const wrapper = mount(EditableSettingValue, {
      props: { modelValue: '10s', onSave: save },
    });

    await wrapper.get('[data-test="edit-setting"]').trigger('click');
    await wrapper.get('input').setValue('20s');
    await wrapper.get('[data-test="save-setting"]').trigger('click');

    expect(save).toHaveBeenCalledWith('20s');
    expect(wrapper.emitted('update:modelValue')).toEqual([['20s']]);
    expect(wrapper.find('input').exists()).toBe(false);
  });

  it('cancels without changing the value', async () => {
    const save = vi.fn();
    const wrapper = mount(EditableSettingValue, {
      props: { modelValue: '10', onSave: save },
    });

    await wrapper.get('[data-test="edit-setting"]').trigger('click');
    await wrapper.get('input').setValue('99');
    await wrapper.get('[data-test="cancel-setting"]').trigger('click');

    expect(save).not.toHaveBeenCalled();
    expect(wrapper.text()).toContain('10');
  });
});
