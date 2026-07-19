import { describe, expect, it } from 'vitest';
import { DEFAULT_CONTENT_ZOOM, MAX_CONTENT_ZOOM, MIN_CONTENT_ZOOM, adjustContentZoom, normalizeContentZoom } from './content-zoom';

describe('content zoom', () => {
  it('changes zoom in ten-percent steps', () => {
    expect(adjustContentZoom(100, 1)).toBe(110);
    expect(adjustContentZoom(100, -1)).toBe(90);
  });

  it('clamps zoom to supported boundaries', () => {
    expect(adjustContentZoom(MAX_CONTENT_ZOOM, 1)).toBe(MAX_CONTENT_ZOOM);
    expect(adjustContentZoom(MIN_CONTENT_ZOOM, -1)).toBe(MIN_CONTENT_ZOOM);
  });

  it('normalizes persisted values', () => {
    expect(normalizeContentZoom('130')).toBe(130);
    expect(normalizeContentZoom(999)).toBe(MAX_CONTENT_ZOOM);
    expect(normalizeContentZoom('invalid')).toBe(DEFAULT_CONTENT_ZOOM);
  });
});
