export namespace main {
	
	export class Config {
	    APIKEY?: any;
	    PROXY_URL?: any;
	    HOST?: any;
	    API_TIMEOUT_MS?: any;
	    LOG?: any;
	    Providers?: any;
	    Router?: any;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.APIKEY = source["APIKEY"];
	        this.PROXY_URL = source["PROXY_URL"];
	        this.HOST = source["HOST"];
	        this.API_TIMEOUT_MS = source["API_TIMEOUT_MS"];
	        this.LOG = source["LOG"];
	        this.Providers = source["Providers"];
	        this.Router = source["Router"];
	    }
	}

}

