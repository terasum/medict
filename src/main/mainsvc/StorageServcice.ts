import { LowSync, JSONFileSync } from 'lowdb';
import { Config } from '../../model/Config';

import fs from 'fs';

export class StorageService {
  dbpath: string;
  db: LowSync<Config>;
  constructor(dbpath: string) {
    this.dbpath = dbpath;
    if (!fs.existsSync(dbpath)) {
      throw new Error('file not found: ' + dbpath);
    }
    this.db = new LowSync(new JSONFileSync<Config>(dbpath));
    this.db.read();

    this.db.data ||= new Config();

  }

  getDataByKey(key: string) {
    if (!this.db.data) {
      return undefined;
    }

    if (this.db.data.hasOwnProperty(key)) {
      return this.db.data[key];
    }
  }

  setDataByKey(key: string, data: any) {
    if (!this.db.data) {
      return false;
    }
    if (this.db.data.hasOwnProperty(key)) {
      this.db.data[key] = data;
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
}
