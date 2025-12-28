package models

import "time"

// Clean service types

type CleanCategory struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	EstimatedMB int64  `json:"estimatedMB"`
}

type CleanProgress struct {
	Category    string  `json:"category"`
	Message     string  `json:"message"`
	Percent     int     `json:"percent"`
	CurrentFile string  `json:"currentFile"`
	TotalFiles  int     `json:"totalFiles"`
	FilesClean  int     `json:"filesClean"`
}

type CleanResult struct {
	SpaceFreed   int64    `json:"spaceFreed"`
	FilesRemoved int      `json:"filesRemoved"`
	Categories   []string `json:"categories"`
	Errors       []string `json:"errors"`
}

// Uninstall service types

type Application struct {
	Name         string    `json:"name"`
	BundleID     string    `json:"bundleId"`
	Path         string    `json:"path"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
	Age          string    `json:"age"`
	Icon         string    `json:"icon,omitempty"`
}

type UninstallProgress struct {
	App           string `json:"app"`
	Message       string `json:"message"`
	Percent       int    `json:"percent"`
	FilesRemoved  int    `json:"filesRemoved"`
	TotalFiles    int    `json:"totalFiles"`
	SpaceFreed    int64  `json:"spaceFreed"`
}

type UninstallResult struct {
	AppsRemoved  int      `json:"appsRemoved"`
	FilesRemoved int      `json:"filesRemoved"`
	SpaceFreed   int64    `json:"spaceFreed"`
	Errors       []string `json:"errors"`
}

// Optimize service types

type OptimizationTask struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	RequiresSudo bool  `json:"requiresSudo"`
}

type OptimizeProgress struct {
	Task    string `json:"task"`
	Message string `json:"message"`
	Percent int    `json:"percent"`
}

type OptimizeResult struct {
	TasksCompleted int      `json:"tasksCompleted"`
	Errors         []string `json:"errors"`
}

// Analyze service types

type FileEntry struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Size int64  `json:"size"`
}

type DirEntry struct {
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	Size       int64     `json:"size"`
	IsDir      bool      `json:"isDir"`
	LastAccess time.Time `json:"lastAccess"`
	Percent    float64   `json:"percent"`
}

type ScanResult struct {
	Entries    []DirEntry  `json:"entries"`
	LargeFiles []FileEntry `json:"largeFiles"`
	TotalSize  int64       `json:"totalSize"`
	Path       string      `json:"path"`
}

type ScanProgress struct {
	Path         string `json:"path"`
	ItemsScanned int    `json:"itemsScanned"`
	TotalSize    int64  `json:"totalSize"`
}

// Status service types

type MetricsSnapshot struct {
	// Hardware info
	Hardware HardwareInfo `json:"hardware"`

	// Health score (0-100)
	Health int `json:"health"`

	// CPU metrics
	CPU CPUMetrics `json:"cpu"`

	// GPU metrics
	GPU GPUMetrics `json:"gpu"`

	// Memory metrics
	Memory MemoryMetrics `json:"memory"`

	// Disk metrics
	Disk DiskMetrics `json:"disk"`

	// Network metrics
	Network NetworkMetrics `json:"network"`

	// Battery/Power metrics
	Battery BatteryMetrics `json:"battery"`

	// Top processes
	Processes []ProcessInfo `json:"processes"`

	// Timestamp
	Timestamp time.Time `json:"timestamp"`
}

type HardwareInfo struct {
	Model      string `json:"model"`
	Processor  string `json:"processor"`
	Memory     string `json:"memory"`
	OS         string `json:"os"`
	OSVersion  string `json:"osVersion"`
	Uptime     string `json:"uptime"`
}

type CPUMetrics struct {
	TotalPercent float64   `json:"totalPercent"`
	LoadAvg      []float64 `json:"loadAvg"` // 1m, 5m, 15m
	Cores        int       `json:"cores"`
	PerCore      []float64 `json:"perCore"`
	Temperature  float64   `json:"temperature"`
}

type GPUMetrics struct {
	Usage       float64 `json:"usage"`
	Temperature float64 `json:"temperature"`
	Name        string  `json:"name"`
}

type MemoryMetrics struct {
	Used      int64   `json:"used"`
	Total     int64   `json:"total"`
	Free      int64   `json:"free"`
	Available int64   `json:"available"`
	Percent   float64 `json:"percent"`
}

type DiskMetrics struct {
	Used        int64   `json:"used"`
	Total       int64   `json:"total"`
	Free        int64   `json:"free"`
	Percent     float64 `json:"percent"`
	ReadBytes   int64   `json:"readBytes"`
	WriteBytes  int64   `json:"writeBytes"`
	ReadSpeed   float64 `json:"readSpeed"`   // MB/s
	WriteSpeed  float64 `json:"writeSpeed"`  // MB/s
}

type NetworkMetrics struct {
	Download    float64 `json:"download"`    // MB/s
	Upload      float64 `json:"upload"`      // MB/s
	ProxyHost   string  `json:"proxyHost"`
	ProxyPort   string  `json:"proxyPort"`
	ProxyType   string  `json:"proxyType"`
	BluetoothOn bool    `json:"bluetoothOn"`
}

type BatteryMetrics struct {
	Level      int     `json:"level"`      // Percentage
	Status     string  `json:"status"`     // Charging, Charged, Discharging
	Health     string  `json:"health"`     // Normal, Replace Soon, Replace Now
	Cycles     int     `json:"cycles"`
	Temperature float64 `json:"temperature"`
	FanSpeed   int     `json:"fanSpeed"`   // RPM
}

type ProcessInfo struct {
	Name       string  `json:"name"`
	PID        int     `json:"pid"`
	CPUPercent float64 `json:"cpuPercent"`
	MemoryMB   int64   `json:"memoryMB"`
}

// Purge service types

type Project struct {
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	Type         string    `json:"type"` // node_modules, target, venv, etc.
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
	IsRecent     bool      `json:"isRecent"` // < 7 days
	Selected     bool      `json:"selected"`
}

type PurgeProgress struct {
	Project  string `json:"project"`
	Message  string `json:"message"`
	Percent  int    `json:"percent"`
	Size     int64  `json:"size"`
}

type PurgeResult struct {
	ProjectsCleaned int      `json:"projectsCleaned"`
	SpaceFreed      int64    `json:"spaceFreed"`
	Errors          []string `json:"errors"`
}

// TouchID service types

type TouchIDStatus struct {
	Enabled bool   `json:"enabled"`
	Status  string `json:"status"`
}

// Common types

type ErrorResponse struct {
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
