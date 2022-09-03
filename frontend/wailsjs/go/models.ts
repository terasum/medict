export namespace model {
	
	export class PlainDictItem {
	    id: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new PlainDictItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}
	export class WrappedWordItem {
	    dict_id: string;
	    raw_key_word: string;
	    key_word: string;
	    record_start: number;
	
	    static createFrom(source: any = {}) {
	        return new WrappedWordItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dict_id = source["dict_id"];
	        this.raw_key_word = source["raw_key_word"];
	        this.key_word = source["key_word"];
	        this.record_start = source["record_start"];
	    }
	}

}

