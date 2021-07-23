import { StorabeDictionary } from './StorableDictionary';

export declare class DictItem {
  dictid: string;
  dictionary: StorabeDictionary;
}

export class Config {
  translateApis: {
    baidu: {
      appid: string;
      appkey: string;
    };
  };
  dicts: DictItem[];
  constructor() {
    this.dicts = [];
    this.translateApis = { baidu: { appid: '', appkey: '' } };
  }
}
