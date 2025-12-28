package analyze

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"mole-wails/backend/models"
)

// cacheManager provides in-memory caching for scan results
type cacheManager struct {
	mu    sync.RWMutex
	cache map[string]*scanResult
}

func newCacheManager() *cacheManager {
	return &cacheManager{
		cache: make(map[string]*scanResult),
	}
}

func (cm *cacheManager) get(path string) (*scanResult, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	result, found := cm.cache[path]
	return result, found
}

func (cm *cacheManager) set(path string, result *scanResult) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.cache[path] = result
}

func (cm *cacheManager) invalidate(path string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.cache, path)
}

type Service struct {
	cache *cacheManager
	ctx   context.Context
}

func NewService() *Service {
	return &Service{
		cache: newCacheManager(),
	}
}

func (s *Service) SetContext(ctx context.Context) {
	s.ctx = ctx
}

// ScanDirectory scans a directory and returns results
func (s *Service) ScanDirectory(path string) (*models.ScanResult, error) {
	// Check cache first
	if cached, found := s.cache.get(path); found {
		return s.convertToModelScanResult(cached), nil
	}

	// Perform scan
	result, err := scanDirectoryInternal(path, s.ctx)
	if err != nil {
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	// Cache result
	s.cache.set(path, result)

	return s.convertToModelScanResult(result), nil
}

// GetLargeFiles returns the largest files in a directory
func (s *Service) GetLargeFiles(path string, limit int) ([]models.FileEntry, error) {
	result, err := s.ScanDirectory(path)
	if err != nil {
		return nil, err
	}

	// Return top N large files
	if limit > len(result.LargeFiles) {
		limit = len(result.LargeFiles)
	}

	return result.LargeFiles[:limit], nil
}

// DeletePath deletes a file or directory with safety checks
func (s *Service) DeletePath(path string) error {
	// Safety checks
	if err := validatePathForDeletion(path); err != nil {
		return fmt.Errorf("path validation failed: %w", err)
	}

	// Use safe deletion
	if err := safeDelete(path); err != nil {
		return fmt.Errorf("deletion failed: %w", err)
	}

	return nil
}

// OpenInFinder opens a path in Finder
func (s *Service) OpenInFinder(path string) error {
	cmd := exec.Command("open", "-R", path)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open in Finder: %w", err)
	}
	return nil
}

// Helper functions

func (s *Service) convertToModelScanResult(internal *scanResult) *models.ScanResult {
	result := &models.ScanResult{
		Entries:    make([]models.DirEntry, len(internal.Entries)),
		LargeFiles: make([]models.FileEntry, len(internal.LargeFiles)),
		TotalSize:  internal.TotalSize,
		Path:       "", // Will be set by caller
	}

	for i, entry := range internal.Entries {
		result.Entries[i] = models.DirEntry{
			Name:       entry.Name,
			Path:       entry.Path,
			Size:       entry.Size,
			IsDir:      entry.IsDir,
			LastAccess: entry.LastAccess,
			Percent:    float64(entry.Size) / float64(internal.TotalSize) * 100,
		}
	}

	for i, file := range internal.LargeFiles {
		result.LargeFiles[i] = models.FileEntry{
			Name: file.Name,
			Path: file.Path,
			Size: file.Size,
		}
	}

	return result
}

func validatePathForDeletion(path string) error {
	// Protected paths
	protectedPaths := []string{
		"/",
		"/System",
		"/bin",
		"/sbin",
		"/usr",
		"/etc",
		"/var",
		"/private",
		"/Applications",
		os.Getenv("HOME"),
	}

	for _, protected := range protectedPaths {
		if path == protected {
			return fmt.Errorf("cannot delete protected path: %s", path)
		}
	}

	// Must be absolute path
	if !filepath.IsAbs(path) {
		return fmt.Errorf("path must be absolute: %s", path)
	}

	return nil
}

func safeDelete(path string) error {
	// Move to trash instead of permanent deletion
	cmd := exec.Command("osascript", "-e", fmt.Sprintf(
		`tell application "Finder" to delete POSIX file "%s"`, path))

	if err := cmd.Run(); err != nil {
		// Fallback to rm if osascript fails
		cmd = exec.Command("rm", "-rf", path)
		return cmd.Run()
	}

	return nil
}

// scanDirectoryInternal performs the actual directory scan with progress reporting
func scanDirectoryInternal(path string, ctx context.Context) (*scanResult, error) {
	// Check disk cache first
	if cached, err := loadCacheFromDisk(path); err == nil {
		result := &scanResult{
			Entries:    cached.Entries,
			LargeFiles: cached.LargeFiles,
			TotalSize:  cached.TotalSize,
		}
		return result, nil
	}

	// Initialize progress counters
	var filesScanned, dirsScanned, bytesScanned int64
	var currentPath string

	// Progress callback for emitting events to frontend
	progressCallback := func() {
		if ctx != nil {
			runtime.EventsEmit(ctx, "analyze:progress", models.ScanProgress{
				Path:         path,
				ItemsScanned: int(atomic.LoadInt64(&filesScanned) + atomic.LoadInt64(&dirsScanned)),
				TotalSize:    atomic.LoadInt64(&bytesScanned),
			})
		}
	}

	// Start periodic progress updates
	stopProgress := make(chan struct{})
	defer close(stopProgress)

	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				progressCallback()
			case <-stopProgress:
				return
			}
		}
	}()

	// Perform the scan using the concurrent scanner
	result, err := scanPathConcurrent(path, &filesScanned, &dirsScanned, &bytesScanned, &currentPath)
	if err != nil {
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	// Final progress update
	progressCallback()

	// Cache the result to disk
	_ = saveCacheToDisk(path, result)

	return &result, nil
}
