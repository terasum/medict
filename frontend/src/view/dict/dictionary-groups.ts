export const DEFAULT_DICTIONARY_GROUP_ID = 'default';

export interface DictionaryGroup {
  id: string;
  name: string;
  dictIds: string[];
}

function uniqueValidIDs(ids: unknown, validIDs: Set<string>) {
  if (!Array.isArray(ids)) return [];
  return [...new Set(ids.filter((id): id is string => typeof id === 'string' && validIDs.has(id)))];
}

export function parseDictionaryGroups(raw: unknown, validDictionaryIDs: string[]): DictionaryGroup[] {
  let value = raw;
  if (typeof raw === 'string') {
    try {
      value = JSON.parse(raw);
    } catch {
      return [];
    }
  }
  if (!Array.isArray(value)) return [];

  const validIDs = new Set(validDictionaryIDs);
  const seenGroupIDs = new Set<string>();
  const seenNames = new Set<string>();
  const groups: DictionaryGroup[] = [];

  for (const candidate of value) {
    if (!candidate || typeof candidate !== 'object') continue;
    const record = candidate as Record<string, unknown>;
    const id = typeof record.id === 'string' ? record.id.trim() : '';
    const name = typeof record.name === 'string' ? record.name.trim() : '';
    if (!id || !name || id === DEFAULT_DICTIONARY_GROUP_ID || seenGroupIDs.has(id) || seenNames.has(name)) continue;
    seenGroupIDs.add(id);
    seenNames.add(name);
    groups.push({ id, name, dictIds: uniqueValidIDs(record.dictIds, validIDs) });
  }
  return groups;
}

export function createDictionaryGroup(groups: DictionaryGroup[], rawName: string, id: string): DictionaryGroup[] {
  const name = rawName.trim();
  if (!name) throw new Error('请输入词典组名称');
  if (groups.some((group) => group.name === name) || name === '默认组') throw new Error('词典组名称已存在');
  return [...groups, { id, name, dictIds: [] }];
}

export function deleteDictionaryGroup(groups: DictionaryGroup[], id: string): DictionaryGroup[] {
  if (id === DEFAULT_DICTIONARY_GROUP_ID) throw new Error('默认组不能删除');
  return groups.filter((group) => group.id !== id);
}

export function updateDictionaryGroupMembers(
  groups: DictionaryGroup[],
  groupID: string,
  dictionaryIDs: string[],
  validDictionaryIDs: string[],
): DictionaryGroup[] {
  const validIDs = new Set(validDictionaryIDs);
  const nextIDs = uniqueValidIDs(dictionaryIDs, validIDs);
  return groups.map((group) => group.id === groupID ? { ...group, dictIds: nextIDs } : group);
}

export function filterDictionaryGroupMembers<T extends { id: string }>(dictionaries: T[], group: DictionaryGroup | null) {
  if (!group) return dictionaries;
  const ids = new Set(group.dictIds);
  return dictionaries.filter((dictionary) => ids.has(dictionary.id));
}
