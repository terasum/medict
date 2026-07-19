import { describe, expect, it } from 'vitest';
import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';

describe('bookmarks layout', () => {
  it('keeps the shared application footer below the main content', () => {
    const source = readFileSync(resolve(process.cwd(), 'src/view/bookmarks/index.vue'), 'utf8');

    expect(source).toContain('<div class="n-layout-footer">');
    expect(source).toContain('<AppFooter />');
    expect(source).toContain('flex: 0 0 $layout-footer-height');
  });
});
