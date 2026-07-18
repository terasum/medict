/**
 * Copyright (C) 2023 Quan Chen <chenquan_act@163.com>
 *
 * GPL-3.0
 */

import { describe, expect, it } from 'vitest';
import { dictionaryIconKind, isDictionarySelected } from './dictionary-toolbar';

describe('dictionaryIconKind', () => {
  it.each([
    ['ECDICT', 'database'],
    ['online', 'online'],
    ['MDICT', 'language'],
    ['STARDICT', 'book'],
    ['', 'book'],
  ])('maps %s dictionaries to the %s fallback icon', (dictType, expected) => {
    expect(dictionaryIconKind(dictType)).toBe(expected);
  });
});

describe('isDictionarySelected', () => {
  it('normalizes serialized dictionary IDs before comparing them', () => {
    expect(isDictionarySelected(42, '42')).toBe(true);
    expect(isDictionarySelected('ecdict', 'online-bing')).toBe(false);
  });
});
