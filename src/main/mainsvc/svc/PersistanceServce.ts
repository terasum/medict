import { StorabeDictionary } from '../../../model/StorableDictionary';
import { StorageService } from './StorageServcice';
import { getConfigJsonPath } from '../../../config/config';

const storageService = new StorageService(getConfigJsonPath());
export class DictionaryAccessor {
  QueryAllDictonary() {
    return storageService.db.data?.dicts;
  }
  AddNewDictionary(dictID: string, dict: StorabeDictionary) {
    storageService.db.data?.dicts.push({ dictid: dictID, dictionary: dict });
    storageService.db.write();
  }
}
