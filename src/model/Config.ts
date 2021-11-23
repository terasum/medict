import { StorabeDictionary } from './StorableDictionary';

export declare class DictItem {
  dictid: string;
  dictionary: StorabeDictionary;
}

export class Config {
  dictBaseDir: string;
  translateApis: {
    youdao: {
      appid: string;
      appkey: string;
    };
    baidu: {
      appid: string;
      appkey: string;
    };
  };
  dicts: DictItem[];
  constructor() {
    this.dictBaseDir = "";
    this.dicts = [];
    this.translateApis = { youdao: { appid: '', appkey: '' }, baidu: { appid: '', appkey: '' } };
  }
}
