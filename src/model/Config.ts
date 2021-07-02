import { StorabeDictionary } from './StorableDictionary';

declare class DictItem {
  dictid: string;
  dictionary: StorabeDictionary;
}

export class Config {
  dicts: DictItem[];
  constructor() {
    this.dicts = [];
  }
}
