export namespace model {

  export class KeyBlockEntry {
    id: number;
    record_start_offset: number;
    record_end_offset: number;
    keyword: string;
    key_block_idx: number;

    static createFrom(source: any = {}) {
      return new KeyBlockEntry(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) source = JSON.parse(source);
      this.id = source["id"];
      this.record_start_offset = source["record_start_offset"];
      this.record_end_offset = source["record_end_offset"];
      this.keyword = source["keyword"];
      this.key_block_idx = source["key_block_idx"];
    }
  }
  export class PlainDictionaryItem {
    id: string;
    name: string;
    path: string;

    static createFrom(source: any = {}) {
      return new PlainDictionaryItem(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) source = JSON.parse(source);
      this.id = source["id"];
      this.name = source["name"];
      this.path = source["path"];
    }
  }
  export class Resp {
    data: any;
    err: string;
    code: number;

    static createFrom(source: any = {}) {
      return new Resp(source);
    }

    constructor(source: any = {}) {
      if ('string' === typeof source) source = JSON.parse(source);
      this.data = source["data"];
      this.err = source["err"];
      this.code = source["code"];
    }
  }

}
