import { StorabeDictionary } from '../../model/StorableDictionary';
import { ConfigAccessService } from './ConfigAccess.svc';
import { getConfigJsonPath } from '../../config/config';
import { logger } from '../../utils/logger';

const storageService = new ConfigAccessService(getConfigJsonPath());
export class DictionaryAccessor {
  QueryAllDictonary() {
    return storageService.db.data?.dicts;
  }
  AddNewDictionary(dictID: string, dict: StorabeDictionary) {
    storageService.db.data?.dicts.push({ dictid: dictID, dictionary: dict });
    storageService.db.write();
  }
}
