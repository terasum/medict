import { LowSync, JSONFileSync } from 'lowdb';
import { Config, DictItem } from '../../model/Config';
import { StorabeDictionary } from '../../model/StorableDictionary';

import fs from 'fs';

export class ConfigAccessService {
  dbpath: string;
  db: LowSync<Config>;

  constructor(dbpath: string) {
    this.dbpath = dbpath;
    if (!fs.existsSync(dbpath)) {
      throw new Error('file not found: ' + dbpath);
    }
    this.db = new LowSync(new JSONFileSync<Config>(dbpath));
    this.db.read();
    this.db.data || new Config();
  }

  getDataByKey(key: string) {
    if (!this.db.data) {
      return undefined;
    }

    if (this.db.data.hasOwnProperty(key)) {
      return this.db.data[key];
    }

    return undefined;
  }

  setDataByKey(key: string, data: any) {
    if (!this.db.data) {
      return false;
    }

    this.db.data[key] = data;
    // write every time
    this.db.write();
    return true;
  }

  setDataByKeyFn(key: string, func: (fd: any) => boolean) {
    if (!this.db.data) {
      return false;
    }

    if (this.db.data.hasOwnProperty(key)) {
      return func(this.db.data[key]);
    }
    return false;
  }

  getBaiduTranslateApis() {
    return this.db.data?.translateApis.baidu;
  }

  setBaiduTranslateApis(arg: { appid: string, appkey: string }) {
    this.db.data!.translateApis.baidu.appid = arg.appid;
    this.db.data!.translateApis.baidu.appkey = arg.appkey;
    // write every time
    this.db.write();
  }

  getYoudaoTranslateApis() {
    return this.db.data?.translateApis.youdao;
  }

  setYoudaoTranslateApis(arg: { appid: string, appkey: string }) {

    this.db.data!.translateApis.youdao.appid = arg.appid;
    this.db.data!.translateApis.youdao.appkey = arg.appkey;
    // write every time
    this.db.write();
  }

  getDictBaseDir(){
    return this.db.data!.dictBaseDir;
  }

  setDictBaseDir(dictBaseDir: string) {
    this.db.data!.dictBaseDir = dictBaseDir;
    // write every time
    this.db.write();
  }

  getDictItems() {
    return this.db.data?.dicts;
  }

  setDictItems(dicts: StorabeDictionary[]) {
    let dictItems: DictItem[] = [];

    dicts.forEach((dict) => {
      dictItems.push({ dictid: dict.id, dictionary: dict })
    })
    this.db.data!.dicts = dictItems;
    // write every time
    this.db.write();
    return true;
  }

}
