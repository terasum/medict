import { describe, expect, it } from 'vitest';
import {
  DEFAULT_DICTIONARY_GROUP_ID,
  createDictionaryGroup,
  deleteDictionaryGroup,
  filterDictionaryGroupMembers,
  parseDictionaryGroups,
  updateDictionaryGroupMembers,
} from './dictionary-groups';

describe('dictionary groups', () => {
  it('restores valid custom groups and removes missing dictionary ids', () => {
    const groups = parseDictionaryGroups(JSON.stringify([
      { id: 'english', name: '英语词典', dictIds: ['d1', 'missing', 'd1'] },
      { id: '', name: 'broken', dictIds: [] },
    ]), ['d1', 'd2']);

    expect(groups).toEqual([{ id: 'english', name: '英语词典', dictIds: ['d1'] }]);
  });

  it('creates named groups and rejects empty or duplicate names', () => {
    const groups = createDictionaryGroup([], ' 英语词典 ', 'english');
    expect(groups).toEqual([{ id: 'english', name: '英语词典', dictIds: [] }]);
    expect(() => createDictionaryGroup(groups, '英语词典', 'duplicate')).toThrow('词典组名称已存在');
    expect(() => createDictionaryGroup(groups, '   ', 'empty')).toThrow('请输入词典组名称');
  });

  it('protects the default group and removes custom groups', () => {
    const groups = [{ id: 'english', name: '英语词典', dictIds: ['d1'] }];
    expect(() => deleteDictionaryGroup(groups, DEFAULT_DICTIONARY_GROUP_ID)).toThrow('默认组不能删除');
    expect(deleteDictionaryGroup(groups, 'english')).toEqual([]);
  });

  it('updates membership and filters the main dictionary list', () => {
    const groups = [{ id: 'english', name: '英语词典', dictIds: ['d1'] }];
    const updated = updateDictionaryGroupMembers(groups, 'english', ['d2', 'd2'], ['d1', 'd2']);
    const dictionaries = [{ id: 'd1' }, { id: 'd2' }, { id: 'd3' }];

    expect(updated[0].dictIds).toEqual(['d2']);
    expect(filterDictionaryGroupMembers(dictionaries, updated[0])).toEqual([{ id: 'd2' }]);
    expect(filterDictionaryGroupMembers(dictionaries, null)).toEqual(dictionaries);
  });
});
