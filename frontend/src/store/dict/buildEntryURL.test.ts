import { describe, expect, it, vi } from 'vitest';

// Mock the wailsjs-backed API layer + static-server helper so importing the
// store module does not pull in the Wails runtime.
vi.mock('@/apis/dicts-api', () => ({
  InitDicts: vi.fn(),
  GetAllDicts: vi.fn(),
  SearchWord: vi.fn(),
  BuildIndex: vi.fn(),
}));
vi.mock('@/apis/apis', () => ({ StaticDictServerURL: vi.fn() }));
vi.mock('@/apis/config', () => ({ getPreferences: vi.fn(), savePreferences: vi.fn() }));

import { buildEntryURL } from './index';

const sampleEntry = {
  keyword: 'apple',
  record_start_offset: 11,
  record_end_offset: 22,
  record_block_data_start_offset: 33,
  record_block_data_compress_size: 44,
  record_block_data_decompress_size: 55,
  keyword_data_start_offset: 66,
  keyword_data_end_offset: 77,
};

describe('buildEntryURL', () => {
  it('encodes dict_id, keyword, entry_id and all offsets', () => {
    const url = buildEntryURL('http://h:1234', 'dictA', sampleEntry, 9);
    expect(url).toContain('/__tcidem_query?dict_id=dictA');
    expect(url).toContain('keyword=apple');
    expect(url).toContain('entry_id=9');
    expect(url).toContain('record_start_offset=11');
    expect(url).toContain('record_end_offset=22');
    expect(url).toContain('record_block_data_start_offset=33');
    expect(url).toContain('record_block_data_decompress_size=55');
    expect(url).toContain('keyword_data_start_offset=66');
  });

  it('uses the given baseurl', () => {
    expect(buildEntryURL('http://x:1', 'd', sampleEntry, 0)).toMatch(/^http:\/\/x:1\/__tcidem_query/);
  });
});
