import { describe, expect, it } from 'vitest';

import { settingsNavigation } from './settings-navigation';

describe('settingsNavigation', () => {
  it('organises settings into clear, non-empty groups with unique routes', () => {
    expect(settingsNavigation.map((group) => group.title)).toEqual([
      '应用设置',
      '扩展与开发',
      '帮助与关于',
    ]);

    const items = settingsNavigation.flatMap((group) => group.items);
    expect(items.every((item) => item.label && item.description)).toBe(true);
    expect(new Set(items.map((item) => item.to)).size).toBe(items.length);
  });
});
