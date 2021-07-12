import { StorabeDictionary } from './StorableDictionary';

export declare class DictItem {
  dictid: string;
  dictionary: StorabeDictionary;
}

export class Config {
  dicts: DictItem[];
  constructor() {
    this.dicts = [];
  }
}
