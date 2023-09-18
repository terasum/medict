export namespace model {
	
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

