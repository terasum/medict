import { StorabeDictionary } from '../../model/StorableDictionary';
import { StorageService } from './StorageServcice';
import { getConfigJsonPath } from '../../config/config';

export class DictionaryAccessor {
  storage: StorageService;

  constructor() {
    this.storage = new StorageService(getConfigJsonPath());
  }
  async QueryAllDictonary() {
    (await this.storage.loadDataByKey(
      'dicts'
    )) as unknown as StorabeDictionary[];
  }
  async AddNewDictionary(dictID: string, dict: StorabeDictionary) {}
}
