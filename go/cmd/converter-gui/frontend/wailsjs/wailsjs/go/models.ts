export namespace main {
	
	export class ConversionResult {
	    SourcePath: string;
	    DestPath: string;
	    Success: boolean;
	    Error: string;
	
	    static createFrom(source: any = {}) {
	        return new ConversionResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.SourcePath = source["SourcePath"];
	        this.DestPath = source["DestPath"];
	        this.Success = source["Success"];
	        this.Error = source["Error"];
	    }
	}

}

