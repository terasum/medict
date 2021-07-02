import { Low, JSONFile } from 'lowdb';

import fs from 'fs';

export class StorageService {
  dbpath: string;
  db: Low;
  constructor(dbpath: string) {
    this.dbpath = dbpath;
    if (!fs.existsSync(dbpath)) {
      throw new Error('file not found: ' + dbpath);
    }
    const adapter = new JSONFile(dbpath);
    this.db = new Low(adapter);
  }

  async load() {
    // Read data from JSON file, this will set db.data content
    await this.db.read();
  }

  async loadDataByKey(key: string) {
    await this.load();
    return (this.db as any).data[key];
  }

  async saveDataByKey(key: string, data: any) {
    await this.load();
    (this.db as any).data[key] = data;
    await this.write();
  }

  async write() {
    await this.db.write();
  }
}
