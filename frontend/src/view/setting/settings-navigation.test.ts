import { describe, expect, it } from 'vitest';

import { settingsNavigation } from './settings-navigation';

describe('settingsNavigation', () => {
  it('provides a flat, single-level menu with unique routes', () => {
    expect(settingsNavigation.map((item) => item.label)).toEqual([
      '词典设置',
      '通用设置',
      '外观设置',
      '插件管理',
      '调试工具',
      '使用说明',
      '关于软件',
    ]);
    expect(settingsNavigation.every((item) => item.label && item.icon)).toBe(true);
    expect(settingsNavigation.every((item) => [...item.label].length === 4)).toBe(true);
    expect(new Set(settingsNavigation.map((item) => item.to)).size).toBe(settingsNavigation.length);
  });

  it('uses the cog icon for debugging tools', () => {
    expect(settingsNavigation.find((item) => item.label === '调试工具')).toEqual({
      label: '调试工具',
      icon: 'icon-cog',
      to: '/setting/debug',
    });
  });

  it('keeps documentation inside the settings route', () => {
    expect(settingsNavigation.find((item) => item.label === '使用说明')?.to).toBe('/setting/docs');
  });
});
