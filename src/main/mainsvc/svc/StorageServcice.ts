import { Low, JSONFile } from 'lowdb';
import lodash from 'lodash';

import { LowSync, JSONFileSync } from 'lowdb';
import {Config} from '../../../model/Config';


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
  }
}
