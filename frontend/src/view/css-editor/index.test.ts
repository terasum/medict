import { mount } from '@vue/test-utils';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const api = vi.hoisted(() => ({
  get: vi.fn().mockResolvedValue('body { color: red; }'),
  save: vi.fn().mockResolvedValue(undefined),
}));

vi.mock('@/apis/dicts-api', () => ({
  getDictEditorCSS: api.get,
  saveDictUserCSS: api.save,
}));
const runtime = vi.hoisted(() => ({ closeHandler: null as null | (() => void) }));
vi.mock('../../../wailsjs/runtime/runtime', () => ({
  EventsOn: vi.fn((_name: string, handler: () => void) => {
    runtime.closeHandler = handler;
    return vi.fn();
  }),
}));
vi.mock('naive-ui', async () => {
  const { defineComponent, h } = await import('vue');
  return {
    NButton: defineComponent({ setup(_, { slots }) { return () => h('button', slots.default?.()); } }),
    useMessage: () => ({ success: vi.fn(), error: vi.fn() }),
  };
});
vi.mock('vue-codemirror', async () => {
  const { defineComponent, h } = await import('vue');
  return {
    Codemirror: defineComponent({
      props: ['modelValue'], emits: ['update:modelValue', 'change'],
      setup(props, { emit }) {
        return () => h('textarea', {
          value: props.modelValue,
          onInput: (event: Event) => {
            emit('update:modelValue', (event.target as HTMLTextAreaElement).value);
            emit('change');
          },
        });
      },
    }),
  };
});

import CSSWindow from './index.vue';

describe('CSS editor window', () => {
  beforeEach(() => {
    vi.useFakeTimers();
    api.get.mockReset().mockResolvedValue('body { color: red; }');
    api.save.mockReset().mockResolvedValue(undefined);
    (window as any).go = { main: { App: {
      CSSWindowDictionary: vi.fn().mockResolvedValue({ id: 'dict-1', name: 'Test Dict' }),
      CloseCSSWindow: vi.fn().mockResolvedValue(undefined),
    } } };
  });

  it('flushes a pending edit before accepting a native close request', async () => {
    const wrapper = mount(CSSWindow);
    await vi.waitFor(() => expect(api.get).toHaveBeenCalledWith('dict-1'));
    await wrapper.find('textarea').setValue('body { color: green; }');
    runtime.closeHandler?.();
    await vi.waitFor(() => expect(api.save).toHaveBeenCalledWith('dict-1', 'body { color: green; }'));
    await vi.waitFor(() => expect((window as any).go.main.App.CloseCSSWindow).toHaveBeenCalledOnce());
    wrapper.unmount();
  });

  it('loads the selected dictionary and autosaves edits', async () => {
    const wrapper = mount(CSSWindow);
    await vi.waitFor(() => expect(api.get).toHaveBeenCalledWith('dict-1'));
    await wrapper.find('textarea').setValue('body { color: blue; }');
    await vi.advanceTimersByTimeAsync(300);
    expect(api.save).toHaveBeenCalledWith('dict-1', 'body { color: blue; }');
    wrapper.unmount();
  });

  it('keeps the window open when flushing before close fails', async () => {
    api.save.mockRejectedValueOnce(new Error('disk full'));
    const wrapper = mount(CSSWindow);
    await vi.waitFor(() => expect(api.get).toHaveBeenCalledWith('dict-1'));
    await wrapper.find('textarea').setValue('invalid close state');
    runtime.closeHandler?.();
    await vi.waitFor(() => expect(api.save).toHaveBeenCalled());
    expect((window as any).go.main.App.CloseCSSWindow).not.toHaveBeenCalled();
    wrapper.unmount();
  });

  it('drains the latest edit before closing during an in-flight save', async () => {
    let resolveFirst!: () => void;
    let resolveSecond!: () => void;
    const first = new Promise<void>((resolve) => { resolveFirst = resolve; });
    const second = new Promise<void>((resolve) => { resolveSecond = resolve; });
    api.save.mockImplementationOnce(() => first).mockImplementationOnce(() => second);
    const wrapper = mount(CSSWindow);
    await vi.waitFor(() => expect(api.get).toHaveBeenCalledWith('dict-1'));
    await wrapper.find('textarea').setValue('first edit');
    await vi.advanceTimersByTimeAsync(300);
    await vi.waitFor(() => expect(api.save).toHaveBeenCalledWith('dict-1', 'first edit'));
    await wrapper.find('textarea').setValue('latest edit');
    runtime.closeHandler?.();
    expect((window as any).go.main.App.CloseCSSWindow).not.toHaveBeenCalled();
    resolveFirst();
    await vi.waitFor(() => expect(api.save).toHaveBeenCalledWith('dict-1', 'latest edit'));
    expect((window as any).go.main.App.CloseCSSWindow).not.toHaveBeenCalled();
    resolveSecond();
    await vi.waitFor(() => expect((window as any).go.main.App.CloseCSSWindow).toHaveBeenCalledOnce());
    wrapper.unmount();
  });

  it('allows closing without overwriting CSS when the initial load fails', async () => {
    api.get.mockRejectedValueOnce(new Error('read failed'));
    const wrapper = mount(CSSWindow);
    await vi.waitFor(() => expect(api.get).toHaveBeenCalledWith('dict-1'));
    runtime.closeHandler?.();
    await vi.waitFor(() => expect((window as any).go.main.App.CloseCSSWindow).toHaveBeenCalledOnce());
    expect(api.save).not.toHaveBeenCalled();
    wrapper.unmount();
  });

  it('allows closing when the editor context cannot be loaded', async () => {
    (window as any).go.main.App.CSSWindowDictionary.mockRejectedValueOnce(new Error('context failed'));
    const wrapper = mount(CSSWindow);
    await Promise.resolve();
    runtime.closeHandler?.();
    await vi.waitFor(() => expect((window as any).go.main.App.CloseCSSWindow).toHaveBeenCalledOnce());
    expect(api.save).not.toHaveBeenCalled();
    wrapper.unmount();
  });

  it('cannot save before the initial CSS load finishes', async () => {
    let resolveLoad!: (value: string) => void;
    api.get.mockImplementationOnce(() => new Promise<string>((resolve) => { resolveLoad = resolve; }));
    const wrapper = mount(CSSWindow);
    await vi.waitFor(() => expect(api.get).toHaveBeenCalledWith('dict-1'));
    await wrapper.findAll('button')[1].trigger('click');
    window.dispatchEvent(new KeyboardEvent('keydown', { key: 's', metaKey: true }));
    expect(api.save).not.toHaveBeenCalled();
    resolveLoad('body {}');
    wrapper.unmount();
  });

  it('enables saving after an empty CSS override finishes loading', async () => {
    api.get.mockResolvedValueOnce('');
    const wrapper = mount(CSSWindow);
    await vi.waitFor(() => expect(api.get).toHaveBeenCalledWith('dict-1'));
    await wrapper.findAll('button')[1].trigger('click');
    await vi.waitFor(() => expect(api.save).toHaveBeenCalledWith('dict-1', ''));
    wrapper.unmount();
  });
});
