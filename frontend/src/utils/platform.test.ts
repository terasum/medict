import { afterEach, describe, expect, it, vi } from 'vitest';
import { detectPlatform, platformClass, sniffPlatformSync } from './platform';

const withUserAgent = (ua: string, fn: () => void | Promise<void>) => {
  Object.defineProperty(window.navigator, 'userAgent', {
    value: ua,
    configurable: true,
  });
  return fn();
};

const goBridge = (goos?: string) => {
  const platform = goos === undefined ? undefined : vi.fn(() => Promise.resolve(goos));
  (window as any).go = platform ? { main: { App: { Platform: platform } } } : undefined;
  return platform;
};

afterEach(() => {
  delete (window as any).go;
});

describe('sniffPlatformSync', () => {
  it('classifies windows, mac and linux user agents', () => {
    withUserAgent(
      'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
      () => expect(sniffPlatformSync()).toBe('windows')
    );
    withUserAgent(
      'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15',
      () => expect(sniffPlatformSync()).toBe('darwin')
    );
    withUserAgent(
      'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36',
      () => expect(sniffPlatformSync()).toBe('linux')
    );
    withUserAgent('curl/8.0', () => expect(sniffPlatformSync()).toBe('unknown'));
  });
});

describe('detectPlatform', () => {
  it('prefers the Go runtime.GOOS answer over the user agent', async () => {
    withUserAgent('Mozilla/5.0 (Windows NT 10.0; Win64; x64)', async () => {
      const platform = goBridge('linux');
      await expect(detectPlatform()).resolves.toBe('linux');
      expect(platform).toHaveBeenCalled();
    });
  });

  it('rejects unknown values from the bridge and falls back to the user agent', async () => {
    withUserAgent('Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)', async () => {
      goBridge('plan9');
      await expect(detectPlatform()).resolves.toBe('darwin');
    });
  });

  it('falls back to the user agent when the bridge is missing or throws', async () => {
    withUserAgent('Mozilla/5.0 (X11; Linux x86_64)', async () => {
      goBridge(undefined);
      await expect(detectPlatform()).resolves.toBe('linux');

      goBridge('windows');
      (window as any).go.main.App.Platform = vi.fn(() => Promise.reject(new Error('ipc down')));
      await expect(detectPlatform()).resolves.toBe('linux');
    });
  });
});

describe('platformClass', () => {
  it('prefixes the platform for the app-root class binding', () => {
    expect(platformClass('windows')).toBe('os-windows');
    expect(platformClass('linux')).toBe('os-linux');
    expect(platformClass('darwin')).toBe('os-darwin');
    expect(platformClass('unknown')).toBe('os-unknown');
  });
});
