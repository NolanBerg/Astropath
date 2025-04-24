export namespace doctor {
	
	export class DoctorReport {
	    SystemHasFFMPEG: boolean;
	    SystemHasARWConversion: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DoctorReport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.SystemHasFFMPEG = source["SystemHasFFMPEG"];
	        this.SystemHasARWConversion = source["SystemHasARWConversion"];
	    }
	}

}

export namespace image {
	
	export class Point {
	    X: number;
	    Y: number;
	
	    static createFrom(source: any = {}) {
	        return new Point(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.X = source["X"];
	        this.Y = source["Y"];
	    }
	}
	export class Rectangle {
	    Min: Point;
	    Max: Point;
	
	    static createFrom(source: any = {}) {
	        return new Rectangle(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Min = this.convertValues(source["Min"], Point);
	        this.Max = this.convertValues(source["Max"], Point);
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

export namespace main {
	
	export class ImportResult {
	    ErrorMessage: string;
	    FilePaths: string[];
	    Bounds: image.Rectangle;
	
	    static createFrom(source: any = {}) {
	        return new ImportResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ErrorMessage = source["ErrorMessage"];
	        this.FilePaths = source["FilePaths"];
	        this.Bounds = this.convertValues(source["Bounds"], image.Rectangle);
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

export namespace processing {
	
	export class ImageBatch {
	    FilePaths: string[];
	    ImageBounds: image.Rectangle;
	
	    static createFrom(source: any = {}) {
	        return new ImageBatch(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.FilePaths = source["FilePaths"];
	        this.ImageBounds = this.convertValues(source["ImageBounds"], image.Rectangle);
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

export namespace settings {
	
	export class AppSettings {
	    UserFirstTime: boolean;
	    Workflows: workflow.Workflow[];
	    EnableTimelapseGeneration: boolean;
	    EnableARWConversion: boolean;
	    ARWTempFilePath: string;
	
	    static createFrom(source: any = {}) {
	        return new AppSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.UserFirstTime = source["UserFirstTime"];
	        this.Workflows = this.convertValues(source["Workflows"], workflow.Workflow);
	        this.EnableTimelapseGeneration = source["EnableTimelapseGeneration"];
	        this.EnableARWConversion = source["EnableARWConversion"];
	        this.ARWTempFilePath = source["ARWTempFilePath"];
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

export namespace workflow {
	
	export class Workflow {
	    UID: number;
	    Name: string;
	    BlendingMode: number;
	    OutputLocation: string;
	    OutputFileName: string;
	    OutputFormat: number;
	    CreateTimelapseVideo: boolean;
	    TimelapseLocation: string;
	    TimelapseFramesLocation: string;
	    DeleteFramesAfterProcessing: boolean;
	    TimelapseDuration: number;
	
	    static createFrom(source: any = {}) {
	        return new Workflow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.UID = source["UID"];
	        this.Name = source["Name"];
	        this.BlendingMode = source["BlendingMode"];
	        this.OutputLocation = source["OutputLocation"];
	        this.OutputFileName = source["OutputFileName"];
	        this.OutputFormat = source["OutputFormat"];
	        this.CreateTimelapseVideo = source["CreateTimelapseVideo"];
	        this.TimelapseLocation = source["TimelapseLocation"];
	        this.TimelapseFramesLocation = source["TimelapseFramesLocation"];
	        this.DeleteFramesAfterProcessing = source["DeleteFramesAfterProcessing"];
	        this.TimelapseDuration = source["TimelapseDuration"];
	    }
	}

}

