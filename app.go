package main

import (
	"context"
	"path/filepath"
	"runtime"

	"mole-wails/backend/analyze"
	"mole-wails/backend/models"
	"mole-wails/backend/services"
	"mole-wails/backend/status"
)

// App struct
type App struct {
	ctx context.Context

	// Services
	Clean     *services.CleanService
	Uninstall *services.UninstallService
	Optimize  *services.OptimizeService
	Analyze   *analyze.Service
	Status    *status.Service
	TouchID   *services.TouchIDService
}

// NewApp creates a new App application struct
func NewApp() *App {
	// Determine scripts path
	scriptsPath := getScriptsPath()

	return &App{
		Clean:     services.NewCleanService(scriptsPath),
		Uninstall: services.NewUninstallService(scriptsPath),
		Optimize:  services.NewOptimizeService(scriptsPath),
		Analyze:   analyze.NewService(),
		Status:    status.NewService(),
		TouchID:   services.NewTouchIDService(scriptsPath),
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Set context for all services
	a.Clean.SetContext(ctx)
	a.Uninstall.SetContext(ctx)
	a.Optimize.SetContext(ctx)
	a.Analyze.SetContext(ctx)
	a.Status.SetContext(ctx)
	a.TouchID.SetContext(ctx)
}

// shutdown is called when the app shuts down
func (a *App) shutdown(ctx context.Context) {
	// Cleanup
	a.Status.StopMonitoring()
}

// Helper function to determine scripts path
func getScriptsPath() string {
	if runtime.GOOS == "darwin" {
		scriptsPath, err := filepath.Abs("scripts")
		if err == nil {
			return scriptsPath
		}
	}

	return "/Applications/Mole.app/Contents/Resources/scripts"
}

// ===========================
// Clean Service Methods
// ===========================

func (a *App) CleanScanTargets() ([]models.CleanCategory, error) {
	return a.Clean.ScanTargets()
}

func (a *App) CleanExecute(categories []string, dryRun bool) error {
	return a.Clean.ExecuteClean(categories, dryRun)
}

func (a *App) CleanGetWhitelist() ([]string, error) {
	return a.Clean.GetWhitelist()
}

func (a *App) CleanUpdateWhitelist(paths []string) error {
	return a.Clean.UpdateWhitelist(paths)
}

// ===========================
// Uninstall Service Methods
// ===========================

func (a *App) UninstallScanApps(forceRescan bool) ([]models.Application, error) {
	return a.Uninstall.ScanApplications(forceRescan)
}

func (a *App) UninstallApps(bundleIDs []string) error {
	return a.Uninstall.UninstallApps(bundleIDs)
}

func (a *App) UninstallGetRelatedFiles(bundleID string) ([]string, error) {
	return a.Uninstall.GetRelatedFiles(bundleID)
}

// ===========================
// Optimize Service Methods
// ===========================

func (a *App) OptimizeGetTasks() ([]models.OptimizationTask, error) {
	return a.Optimize.GetTasks()
}

func (a *App) OptimizeExecute(taskIDs []string) error {
	return a.Optimize.ExecuteOptimizations(taskIDs)
}

func (a *App) OptimizeGetWhitelist() ([]string, error) {
	return a.Optimize.GetWhitelist()
}

func (a *App) OptimizeUpdateWhitelist(tasks []string) error {
	return a.Optimize.UpdateWhitelist(tasks)
}

// ===========================
// Analyze Service Methods
// ===========================

func (a *App) AnalyzeScanDirectory(path string) (*models.ScanResult, error) {
	return a.Analyze.ScanDirectory(path)
}

func (a *App) AnalyzeGetLargeFiles(path string, limit int) ([]models.FileEntry, error) {
	return a.Analyze.GetLargeFiles(path, limit)
}

func (a *App) AnalyzeDeletePath(path string) error {
	return a.Analyze.DeletePath(path)
}

func (a *App) AnalyzeOpenInFinder(path string) error {
	return a.Analyze.OpenInFinder(path)
}

// ===========================
// Status Service Methods
// ===========================

func (a *App) StatusGetMetrics() (*models.MetricsSnapshot, error) {
	return a.Status.GetMetrics()
}

func (a *App) StatusStartMonitoring(interval int) error {
	return a.Status.StartMonitoring(interval)
}

func (a *App) StatusStopMonitoring() {
	a.Status.StopMonitoring()
}

// ===========================
// TouchID Service Methods
// ===========================

func (a *App) TouchIDGetStatus() (*models.TouchIDStatus, error) {
	return a.TouchID.GetStatus()
}

func (a *App) TouchIDEnable() error {
	return a.TouchID.Enable()
}

func (a *App) TouchIDDisable() error {
	return a.TouchID.Disable()
}
