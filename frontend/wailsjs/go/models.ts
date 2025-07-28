export namespace main {
	
	export class Router {
	    default?: string;
	    background?: string;
	    think?: string;
	    longContext?: string;
	    longContextThreshold?: number;
	    webSearch?: string;
	
	    static createFrom(source: any = {}) {
	        return new Router(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.default = source["default"];
	        this.background = source["background"];
	        this.think = source["think"];
	        this.longContext = source["longContext"];
	        this.longContextThreshold = source["longContextThreshold"];
	        this.webSearch = source["webSearch"];
	    }
	}
	export class Provider {
	    name: string;
	    api_base_url: string;
	    api_key: string;
	    models: string[];
	    transformer?: any;
	
	    static createFrom(source: any = {}) {
	        return new Provider(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.api_base_url = source["api_base_url"];
	        this.api_key = source["api_key"];
	        this.models = source["models"];
	        this.transformer = source["transformer"];
	    }
	}
	export class Config {
	    APIKEY?: string;
	    PROXY_URL?: string;
	    HOST?: string;
	    API_TIMEOUT_MS?: number;
	    LOG?: boolean;
	    Providers?: Provider[];
	    Router?: Router;
	
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
	        this.Providers = this.convertValues(source["Providers"], Provider);
	        this.Router = this.convertValues(source["Router"], Router);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	

}

