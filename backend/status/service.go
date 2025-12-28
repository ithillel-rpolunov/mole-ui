package status

import (
	"context"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"mole-wails/backend/models"
)

type Service struct {
	collector *Collector
	ctx       context.Context
	stopChan  chan struct{}
	running   bool
}

func NewService() *Service {
	return &Service{
		collector: NewCollector(),
		stopChan:  make(chan struct{}),
		running:   false,
	}
}

func (s *Service) SetContext(ctx context.Context) {
	s.ctx = ctx
}

// GetMetrics collects and returns current system metrics
func (s *Service) GetMetrics() (*models.MetricsSnapshot, error) {
	snapshot, err := s.collector.Collect()
	if err != nil {
		return nil, err
	}

	return s.convertToModelMetrics(&snapshot), nil
}

// StartMonitoring starts real-time monitoring with periodic updates
func (s *Service) StartMonitoring(interval int) error {
	if s.running {
		return nil // Already running
	}

	s.running = true
	ticker := time.NewTicker(time.Duration(interval) * time.Second)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				metrics, err := s.GetMetrics()
				if err == nil && s.ctx != nil {
					runtime.EventsEmit(s.ctx, "status:update", metrics)
				}
			case <-s.stopChan:
				return
			}
		}
	}()

	return nil
}

// StopMonitoring stops the monitoring loop
func (s *Service) StopMonitoring() {
	if !s.running {
		return
	}

	s.running = false
	close(s.stopChan)
	s.stopChan = make(chan struct{})
}

// Helper function to convert internal metrics to model
func (s *Service) convertToModelMetrics(snapshot *MetricsSnapshot) *models.MetricsSnapshot {
	// Use the primary (first/largest) disk instead of summing all disks
	// This prevents double-counting on macOS where APFS volumes share the same container
	var totalDiskUsed, totalDiskTotal uint64
	var diskPercent float64

	if len(snapshot.Disks) > 0 {
		// Disks are already sorted by size (largest first) in collectDisks()
		primaryDisk := snapshot.Disks[0]
		totalDiskUsed = primaryDisk.Used
		totalDiskTotal = primaryDisk.Total
		diskPercent = primaryDisk.UsedPercent
	}

	// Calculate total network traffic
	var totalRx, totalTx float64
	for _, net := range snapshot.Network {
		totalRx += net.RxRateMBs
		totalTx += net.TxRateMBs
	}

	// Get primary battery and thermal info
	var batteryLevel int
	var batteryStatus string
	var batteryHealth string
	var batteryCycles int
	if len(snapshot.Batteries) > 0 {
		batteryLevel = int(snapshot.Batteries[0].Percent)
		batteryStatus = snapshot.Batteries[0].Status
		batteryHealth = snapshot.Batteries[0].Health
		batteryCycles = snapshot.Batteries[0].CycleCount
	}

	// Get GPU info (use first GPU if available)
	var gpuUsage, gpuTemp float64
	var gpuName string
	if len(snapshot.GPU) > 0 {
		gpuUsage = snapshot.GPU[0].Usage
		gpuName = snapshot.GPU[0].Name
		// Temperature might be in thermal data or GPU data
		if snapshot.Thermal.GPUTemp > 0 {
			gpuTemp = snapshot.Thermal.GPUTemp
		}
	}

	return &models.MetricsSnapshot{
		Hardware: models.HardwareInfo{
			Model:     snapshot.Hardware.Model,
			Processor: snapshot.Hardware.CPUModel,
			Memory:    snapshot.Hardware.TotalRAM,
			OS:        snapshot.Platform,
			OSVersion: snapshot.Hardware.OSVersion,
			Uptime:    snapshot.Uptime,
		},
		Health: snapshot.HealthScore,
		CPU: models.CPUMetrics{
			TotalPercent: snapshot.CPU.Usage,
			LoadAvg:      []float64{snapshot.CPU.Load1, snapshot.CPU.Load5, snapshot.CPU.Load15},
			Cores:        snapshot.CPU.LogicalCPU,
			PerCore:      snapshot.CPU.PerCore,
			Temperature:  snapshot.Thermal.CPUTemp,
		},
		GPU: models.GPUMetrics{
			Usage:       gpuUsage,
			Temperature: gpuTemp,
			Name:        gpuName,
		},
		Memory: models.MemoryMetrics{
			Used:      int64(snapshot.Memory.Used),
			Total:     int64(snapshot.Memory.Total),
			Free:      int64(snapshot.Memory.Total - snapshot.Memory.Used),
			Available: int64(snapshot.Memory.Total - snapshot.Memory.Used),
			Percent:   snapshot.Memory.UsedPercent,
		},
		Disk: models.DiskMetrics{
			Used:       int64(totalDiskUsed),
			Total:      int64(totalDiskTotal),
			Free:       int64(totalDiskTotal - totalDiskUsed),
			Percent:    diskPercent,
			ReadBytes:  0, // Not tracked in current implementation
			WriteBytes: 0, // Not tracked in current implementation
			ReadSpeed:  snapshot.DiskIO.ReadRate,
			WriteSpeed: snapshot.DiskIO.WriteRate,
		},
		Network: models.NetworkMetrics{
			Download:    totalRx,
			Upload:      totalTx,
			ProxyHost:   snapshot.Proxy.Host,
			ProxyPort:   "",
			ProxyType:   snapshot.Proxy.Type,
			BluetoothOn: len(snapshot.Bluetooth) > 0,
		},
		Battery: models.BatteryMetrics{
			Level:       batteryLevel,
			Status:      batteryStatus,
			Health:      batteryHealth,
			Cycles:      batteryCycles,
			Temperature: snapshot.Thermal.CPUTemp, // Using CPU temp as proxy
			FanSpeed:    snapshot.Thermal.FanSpeed,
		},
		Processes: convertProcesses(snapshot.TopProcesses),
		Timestamp: snapshot.CollectedAt,
	}
}

func convertProcesses(processes []ProcessInfo) []models.ProcessInfo {
	result := make([]models.ProcessInfo, len(processes))
	for i, proc := range processes {
		result[i] = models.ProcessInfo{
			Name:       proc.Name,
			PID:        0, // ProcessInfo in status package doesn't have PID
			CPUPercent: proc.CPU,
			MemoryMB:   int64(proc.Memory),
		}
	}
	return result
}
