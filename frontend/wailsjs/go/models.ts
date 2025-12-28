export namespace models {
	
	export class Application {
	    name: string;
	    bundleId: string;
	    path: string;
	    size: number;
	    // Go type: time
	    lastModified: any;
	    age: string;
	    icon?: string;
	
	    static createFrom(source: any = {}) {
	        return new Application(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.bundleId = source["bundleId"];
	        this.path = source["path"];
	        this.size = source["size"];
	        this.lastModified = this.convertValues(source["lastModified"], null);
	        this.age = source["age"];
	        this.icon = source["icon"];
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
	export class BatteryMetrics {
	    level: number;
	    status: string;
	    health: string;
	    cycles: number;
	    temperature: number;
	    fanSpeed: number;
	
	    static createFrom(source: any = {}) {
	        return new BatteryMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.level = source["level"];
	        this.status = source["status"];
	        this.health = source["health"];
	        this.cycles = source["cycles"];
	        this.temperature = source["temperature"];
	        this.fanSpeed = source["fanSpeed"];
	    }
	}
	export class CPUMetrics {
	    totalPercent: number;
	    loadAvg: number[];
	    cores: number;
	    perCore: number[];
	    temperature: number;
	
	    static createFrom(source: any = {}) {
	        return new CPUMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.totalPercent = source["totalPercent"];
	        this.loadAvg = source["loadAvg"];
	        this.cores = source["cores"];
	        this.perCore = source["perCore"];
	        this.temperature = source["temperature"];
	    }
	}
	export class CleanCategory {
	    id: string;
	    name: string;
	    description: string;
	    enabled: boolean;
	    estimatedMB: number;
	
	    static createFrom(source: any = {}) {
	        return new CleanCategory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.enabled = source["enabled"];
	        this.estimatedMB = source["estimatedMB"];
	    }
	}
	export class DirEntry {
	    name: string;
	    path: string;
	    size: number;
	    isDir: boolean;
	    // Go type: time
	    lastAccess: any;
	    percent: number;
	
	    static createFrom(source: any = {}) {
	        return new DirEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.size = source["size"];
	        this.isDir = source["isDir"];
	        this.lastAccess = this.convertValues(source["lastAccess"], null);
	        this.percent = source["percent"];
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
	export class DiskMetrics {
	    used: number;
	    total: number;
	    free: number;
	    percent: number;
	    readBytes: number;
	    writeBytes: number;
	    readSpeed: number;
	    writeSpeed: number;
	
	    static createFrom(source: any = {}) {
	        return new DiskMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.used = source["used"];
	        this.total = source["total"];
	        this.free = source["free"];
	        this.percent = source["percent"];
	        this.readBytes = source["readBytes"];
	        this.writeBytes = source["writeBytes"];
	        this.readSpeed = source["readSpeed"];
	        this.writeSpeed = source["writeSpeed"];
	    }
	}
	export class FileEntry {
	    name: string;
	    path: string;
	    size: number;
	
	    static createFrom(source: any = {}) {
	        return new FileEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.size = source["size"];
	    }
	}
	export class GPUMetrics {
	    usage: number;
	    temperature: number;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new GPUMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.usage = source["usage"];
	        this.temperature = source["temperature"];
	        this.name = source["name"];
	    }
	}
	export class HardwareInfo {
	    model: string;
	    processor: string;
	    memory: string;
	    os: string;
	    osVersion: string;
	    uptime: string;
	
	    static createFrom(source: any = {}) {
	        return new HardwareInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.model = source["model"];
	        this.processor = source["processor"];
	        this.memory = source["memory"];
	        this.os = source["os"];
	        this.osVersion = source["osVersion"];
	        this.uptime = source["uptime"];
	    }
	}
	export class MemoryMetrics {
	    used: number;
	    total: number;
	    free: number;
	    available: number;
	    percent: number;
	
	    static createFrom(source: any = {}) {
	        return new MemoryMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.used = source["used"];
	        this.total = source["total"];
	        this.free = source["free"];
	        this.available = source["available"];
	        this.percent = source["percent"];
	    }
	}
	export class ProcessInfo {
	    name: string;
	    pid: number;
	    cpuPercent: number;
	    memoryMB: number;
	
	    static createFrom(source: any = {}) {
	        return new ProcessInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.pid = source["pid"];
	        this.cpuPercent = source["cpuPercent"];
	        this.memoryMB = source["memoryMB"];
	    }
	}
	export class NetworkMetrics {
	    download: number;
	    upload: number;
	    proxyHost: string;
	    proxyPort: string;
	    proxyType: string;
	    bluetoothOn: boolean;
	
	    static createFrom(source: any = {}) {
	        return new NetworkMetrics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.download = source["download"];
	        this.upload = source["upload"];
	        this.proxyHost = source["proxyHost"];
	        this.proxyPort = source["proxyPort"];
	        this.proxyType = source["proxyType"];
	        this.bluetoothOn = source["bluetoothOn"];
	    }
	}
	export class MetricsSnapshot {
	    hardware: HardwareInfo;
	    health: number;
	    cpu: CPUMetrics;
	    gpu: GPUMetrics;
	    memory: MemoryMetrics;
	    disk: DiskMetrics;
	    network: NetworkMetrics;
	    battery: BatteryMetrics;
	    processes: ProcessInfo[];
	    // Go type: time
	    timestamp: any;
	
	    static createFrom(source: any = {}) {
	        return new MetricsSnapshot(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hardware = this.convertValues(source["hardware"], HardwareInfo);
	        this.health = source["health"];
	        this.cpu = this.convertValues(source["cpu"], CPUMetrics);
	        this.gpu = this.convertValues(source["gpu"], GPUMetrics);
	        this.memory = this.convertValues(source["memory"], MemoryMetrics);
	        this.disk = this.convertValues(source["disk"], DiskMetrics);
	        this.network = this.convertValues(source["network"], NetworkMetrics);
	        this.battery = this.convertValues(source["battery"], BatteryMetrics);
	        this.processes = this.convertValues(source["processes"], ProcessInfo);
	        this.timestamp = this.convertValues(source["timestamp"], null);
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
	
	export class OptimizationTask {
	    id: string;
	    name: string;
	    description: string;
	    enabled: boolean;
	    requiresSudo: boolean;
	
	    static createFrom(source: any = {}) {
	        return new OptimizationTask(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.enabled = source["enabled"];
	        this.requiresSudo = source["requiresSudo"];
	    }
	}
	
	export class ScanResult {
	    entries: DirEntry[];
	    largeFiles: FileEntry[];
	    totalSize: number;
	    totalItems: number;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new ScanResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.entries = this.convertValues(source["entries"], DirEntry);
	        this.largeFiles = this.convertValues(source["largeFiles"], FileEntry);
	        this.totalSize = source["totalSize"];
	        this.totalItems = source["totalItems"];
	        this.path = source["path"];
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
	export class TouchIDStatus {
	    enabled: boolean;
	    available: boolean;
	    status: string;
	    pamModulePath: string;
	    configPath: string;
	
	    static createFrom(source: any = {}) {
	        return new TouchIDStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this.available = source["available"];
	        this.status = source["status"];
	        this.pamModulePath = source["pamModulePath"];
	        this.configPath = source["configPath"];
	    }
	}

}

