/**
 * Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
 *
 * GPL-3.0
 */

export type DictionaryIconKind = 'database' | 'online' | 'language' | 'book';

export function dictionaryIconKind(dictType: unknown): DictionaryIconKind {
  switch (String(dictType || '').toUpperCase()) {
    case 'ECDICT':
      return 'database';
    case 'ONLINE':
      return 'online';
    case 'MDICT':
      return 'language';
    default:
      return 'book';
  }
}

export function isDictionarySelected(itemId: unknown, selectedId: unknown): boolean {
  return String(itemId ?? '') === String(selectedId ?? '');
}
