/**
 * Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
 *
 * GPL-3.0
 */

export function isDictionarySelected(itemId: unknown, selectedId: unknown): boolean {
  return String(itemId ?? '') === String(selectedId ?? '');
}
