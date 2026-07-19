export const MIN_CONTENT_ZOOM = 80;
export const MAX_CONTENT_ZOOM = 160;
export const DEFAULT_CONTENT_ZOOM = 100;
export const CONTENT_ZOOM_STEP = 10;

export function normalizeContentZoom(value: unknown) {
  const parsed = Number(value);
  if (!Number.isFinite(parsed)) return DEFAULT_CONTENT_ZOOM;
  return Math.min(MAX_CONTENT_ZOOM, Math.max(MIN_CONTENT_ZOOM, Math.round(parsed / CONTENT_ZOOM_STEP) * CONTENT_ZOOM_STEP));
}

export function adjustContentZoom(current: number, direction: -1 | 1) {
  return normalizeContentZoom(current + direction * CONTENT_ZOOM_STEP);
}
