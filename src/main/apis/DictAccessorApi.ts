import { Dictionary } from '../../infra/Dictionary';
import { DictService } from '../mainsvc/DictionaryService';
import { StorabeDictionary } from '../../model/StorableDictionary';

export class DictAccessorApi {
  dictService: DictService
  constructor() {
    this.dictService = new DictService();
  }
  dictFindOne(arg: { dictid: string }) {
    const dict = this.dictService.findOne(arg.dictid);
    if (!dict) {
      return null;
    }
    return StorabeDictionary.clone(dict);
  }

  dictFindAll(arg: any) {
    const list = this.dictService.findAll();
    const newList: StorabeDictionary[] = [];
    list.forEach(item => {
      newList.push(StorabeDictionary.clone(item));
    });
    return newList;
  }

  dictAddOne(arg: { dict: Dictionary }) {
    return this.dictService.addOne(arg.dict);
  }

  dictDeleteOne(arg: { dictid: string }) {
    return this.dictService.deleteOne(arg.dictid);
  }
}

