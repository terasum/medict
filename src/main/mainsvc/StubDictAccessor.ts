import { Dictionary } from '../domain/Dictionary';
import { DictService } from './svc/DictionaryService';
import { StorabeDictionary } from '../../model/StorableDictionary';

const dictService = new DictService();
export class StubDictAccessor {
  dictFindOne(arg: { dictid: string }) {
    const dict = dictService.findOne(arg.dictid);
    if (!dict) {
      return null;
    }
    return StorabeDictionary.clone(dict);
  }

  dictFindAll(arg: any) {
    const list = dictService.findAll();
    const newList: StorabeDictionary[] = [];
    list.forEach(item => {
      newList.push(StorabeDictionary.clone(item));
    });
    return newList;
  }
  dictAddOne(arg: { dict: Dictionary }) {
    return dictService.addOne(arg.dict);
  }
  dictDeleteOne(arg: { dictid: string }) {
    return dictService.deleteOne(arg.dictid);
  }
}
