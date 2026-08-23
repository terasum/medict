import { describe, expect, it } from 'vitest';
import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';

describe('app-root platform layout', () => {
  it('shortens the fake title bar on platforms that keep the native title bar', () => {
    const source = readFileSync(resolve(process.cwd(), 'src/App.vue'), 'utf8');

    // the platform class is bound on #app-root
    expect(source).toContain(':class="platformClass"');

    // macOS keeps the full strip; windows/linux shrink it by 12px
    expect(source).toContain('&.os-windows,');
    expect(source).toContain('&.os-linux');
    expect(source).toContain('#{$fake-title-bar-height - 12px}');

    // both the bar and the provider below it follow the CSS variable
    expect(source).toContain('height: var(--fake-title-bar-height)');
    expect(source).toContain('height: calc(100% - var(--fake-title-bar-height))');
  });
});
