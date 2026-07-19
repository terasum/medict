import { describe, expect, it } from 'vitest';
import { composeResourceRequestURL } from './resource-query';

describe('composeResourceRequestURL', () => {
  const server = 'http://127.0.0.1:9081/__mdict';

  it('composes a relative dictionary resource URL', () => {
    expect(composeResourceRequestURL(server, 'images/logo.png', 'dict-1')).toBe(
      'http://127.0.0.1:9081/__mdict/images/logo.png?dict_id=dict-1&d=0',
    );
  });

  it('resolves a service-root path from the server origin', () => {
    expect(composeResourceRequestURL(server, '/__mdict/style.css', 'dict-2')).toBe(
      'http://127.0.0.1:9081/__mdict/style.css?dict_id=dict-2&d=0',
    );
  });

  it('keeps a complete HTTP URL unchanged', () => {
    const url = 'https://example.com/asset.css?debug=1';
    expect(composeResourceRequestURL(server, url, 'dict-3')).toBe(url);
  });

  it('returns an empty preview until server and resource input are available', () => {
    expect(composeResourceRequestURL('', 'image.png', 'dict-1')).toBe('');
    expect(composeResourceRequestURL(server, '   ', 'dict-1')).toBe('');
  });
});
