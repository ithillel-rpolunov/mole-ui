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
	// Validate path
	if path == "" {
		return nil, fmt.Errorf("path cannot be empty")
	}

	// Check if path exists
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("path does not exist: %s", path)
		}
		return nil, fmt.Errorf("cannot access path: %w", err)
	}

	// Check cache first
	if cached, found := s.cache.get(path); found {
		result := s.convertToModelScanResult(cached)
		result.Path = path
		return result, nil
	}

	// Perform scan
	result, err := scanDirectoryInternal(path, s.ctx)
	if err != nil {
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	// Cache result
	s.cache.set(path, result)

	modelResult := s.convertToModelScanResult(result)
	modelResult.Path = path
	return modelResult, nil
}

// GetLargeFiles returns the largest files in a directory
func (s *Service) GetLargeFiles(path string, limit int) ([]models.FileEntry, error) {
	fmt.Printf("[analyze] GetLargeFiles called: path=%s, limit=%d\n", path, limit)

	result, err := s.ScanDirectory(path)
	if err != nil {
		return nil, err
	}

	fmt.Printf("[analyze] ScanDirectory returned: LargeFiles count=%d\n", len(result.LargeFiles))

	// Return top N large files
	if limit > len(result.LargeFiles) {
		limit = len(result.LargeFiles)
	}

	fmt.Printf("[analyze] Returning %d large files\n", limit)
	for i := 0; i < limit && i < len(result.LargeFiles); i++ {
		fmt.Printf("[analyze] LargeFile[%d]: Name=%s, Size=%d\n", i, result.LargeFiles[i].Name, result.LargeFiles[i].Size)
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
	if internal == nil {
		return &models.ScanResult{
			Entries:    []models.DirEntry{},
			LargeFiles: []models.FileEntry{},
			TotalSize:  0,
			TotalItems: 0,
			Path:       "",
		}
	}

	result := &models.ScanResult{
		Entries:    make([]models.DirEntry, len(internal.Entries)),
		LargeFiles: make([]models.FileEntry, len(internal.LargeFiles)),
		TotalSize:  internal.TotalSize,
		TotalItems: len(internal.Entries),
		Path:       "", // Will be set by caller
	}

	for i, entry := range internal.Entries {
		percent := 0.0
		if internal.TotalSize > 0 {
			percent = float64(entry.Size) / float64(internal.TotalSize) * 100
		}

		result.Entries[i] = models.DirEntry{
			Name:       entry.Name,
			Path:       entry.Path,
			Size:       entry.Size,
			IsDir:      entry.IsDir,
			LastAccess: entry.LastAccess,
			Percent:    percent,
		}
	}

	for i, file := range internal.LargeFiles {
		result.LargeFiles[i] = models.FileEntry{
			Name: file.Name,
			Path: file.Path,
			Size: file.Size,
		}
		fmt.Printf("[analyze] LargeFile[%d]: Name=%s, Path=%s, Size=%d\n", i, file.Name, file.Path, file.Size)
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
		fmt.Printf("[analyze] Using cached scan result for path: %s (TotalSize: %d)\n", path, result.TotalSize)
		return result, nil
	}

	fmt.Printf("[analyze] Starting fresh scan for path: %s\n", path)

	// Initialize progress counters
	var filesScanned, dirsScanned, bytesScanned int64
	var currentPath string

	// Progress callback for emitting events to frontend
	progressCallback := func() {
		if ctx != nil {
			progress := models.ScanProgress{
				Path:         path,
				ItemsScanned: int(atomic.LoadInt64(&filesScanned) + atomic.LoadInt64(&dirsScanned)),
				TotalSize:    atomic.LoadInt64(&bytesScanned),
			}
			runtime.EventsEmit(ctx, "analyze:progress", progress)
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
		fmt.Printf("[analyze] Scan failed for path %s: %v\n", path, err)
		return nil, fmt.Errorf("scan failed: %w", err)
	}

	fmt.Printf("[analyze] Scan completed for path: %s (TotalSize: %d, Entries: %d, LargeFiles: %d)\n",
		path, result.TotalSize, len(result.Entries), len(result.LargeFiles))

	// Final progress update
	progressCallback()

	// Cache the result to disk
	_ = saveCacheToDisk(path, result)

	return &result, nil
}
