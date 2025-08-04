export namespace main {
	
	export class Config {
	    APIKEY?: any;
	    PROXY_URL?: any;
	    HOST?: any;
	    PORT?: any;
	    API_TIMEOUT_MS?: any;
	    LOG?: any;
	    NPM_GLOBAL_PREFIX?: any;
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
	        this.PORT = source["PORT"];
	        this.API_TIMEOUT_MS = source["API_TIMEOUT_MS"];
	        this.LOG = source["LOG"];
	        this.NPM_GLOBAL_PREFIX = source["NPM_GLOBAL_PREFIX"];
	        this.Providers = source["Providers"];
	        this.Router = source["Router"];
	    }
	}
	export class ServiceStatus {
	    isRunning: boolean;
	    pid: number;
	
	    static createFrom(source: any = {}) {
	        return new ServiceStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.isRunning = source["isRunning"];
	        this.pid = source["pid"];
	    }
	}

}

