import { describe, expect, it } from 'vitest';
import { docsNavigation } from './docs-navigation';

describe('docsNavigation', () => {
  it('provides unique documentation routes inside settings', () => {
    expect(docsNavigation.map((item) => item.label)).toEqual([
      '界面介绍',
      '词典配置与使用',
      '常见问题',
      '隐私声明',
      '开源许可',
    ]);
    expect(docsNavigation.every((item) => item.to.startsWith('/setting/docs/'))).toBe(true);
    expect(new Set(docsNavigation.map((item) => item.to)).size).toBe(docsNavigation.length);
  });
});
