package services

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"mole-wails/backend/models"
)

// cleanCategory represents a category of files to clean
type cleanCategory struct {
	id          string
	name        string
	description string
	paths       []string
	getSizeFn   func(string) (int64, error) // Custom size calculation if needed
}

type CleanService struct {
	ctx        context.Context
	categories []cleanCategory
	whitelist  map[string]bool
}

func NewCleanService(scriptsPath string) *CleanService {
	homeDir := os.Getenv("HOME")

	return &CleanService{
		whitelist:  make(map[string]bool),
		categories: []cleanCategory{
			{
				id:          "system-caches",
				name:        "System Caches",
				description: "Cached data from system and applications",
				paths: []string{
					filepath.Join(homeDir, "Library", "Caches"),
				},
			},
			{
				id:          "user-logs",
				name:        "User Logs",
				description: "Log files from applications and system",
				paths: []string{
					filepath.Join(homeDir, "Library", "Logs"),
				},
			},
			{
				id:          "temp-files",
				name:        "Temporary Files",
				description: "System temporary files",
				paths: []string{
					"/tmp",
					"/private/var/tmp",
				},
			},
			{
				id:          "browser-caches",
				name:        "Browser Caches",
				description: "Cache files from web browsers",
				paths: []string{
					filepath.Join(homeDir, "Library", "Caches", "com.apple.Safari"),
					filepath.Join(homeDir, "Library", "Caches", "Google", "Chrome"),
					filepath.Join(homeDir, "Library", "Caches", "Firefox"),
					filepath.Join(homeDir, "Library", "Caches", "Chromium"),
				},
			},
			{
				id:          "app-caches",
				name:        "Application Caches",
				description: "Cache files from installed applications",
				paths: []string{
					filepath.Join(homeDir, "Library", "Application Support", "CrashReporter"),
					filepath.Join(homeDir, "Library", "Application Support", "Code", "Cache"),
					filepath.Join(homeDir, "Library", "Application Support", "Code", "CachedData"),
					filepath.Join(homeDir, "Library", "Application Support", "Slack", "Cache"),
					filepath.Join(homeDir, "Library", "Application Support", "Spotify", "PersistentCache"),
				},
			},
			{
				id:          "trash",
				name:        "Trash",
				description: "Files in the Trash",
				paths: []string{
					filepath.Join(homeDir, ".Trash"),
				},
			},
			{
				id:          "download-cache",
				name:        "Download Cache",
				description: "Cached download data",
				paths: []string{
					filepath.Join(homeDir, "Library", "Caches", "com.apple.akd"),
					filepath.Join(homeDir, "Library", "Caches", "com.apple.appstore"),
				},
			},
			{
				id:          "mail-cache",
				name:        "Mail Cache",
				description: "Apple Mail cached data",
				paths: []string{
					filepath.Join(homeDir, "Library", "Mail", "V10", "MailData", "Envelope Index-wal"),
					filepath.Join(homeDir, "Library", "Caches", "com.apple.mail"),
				},
			},
		},
	}
}

// SetContext sets the Wails context for event emission
func (s *CleanService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

// ScanTargets performs a scan to get estimated space savings per category
func (s *CleanService) ScanTargets() ([]models.CleanCategory, error) {
	// Load whitelist
	if err := s.loadWhitelist(); err != nil {
		return nil, fmt.Errorf("failed to load whitelist: %w", err)
	}

	var results []models.CleanCategory

	for _, cat := range s.categories {
		estimatedSize := int64(0)

		for _, path := range cat.paths {
			// Check if path exists
			if _, err := os.Stat(path); os.IsNotExist(err) {
				continue
			}

			// Skip whitelisted paths
			if s.isWhitelisted(path) {
				continue
			}

			// Calculate size
			var size int64
			if cat.getSizeFn != nil {
				var err error
				size, err = cat.getSizeFn(path)
				if err != nil {
					// Skip paths we can't access
					continue
				}
			} else {
				var err error
				size, err = s.calculateDirSize(path)
				if err != nil {
					// Skip paths we can't access
					continue
				}
			}

			estimatedSize += size
		}

		// Convert bytes to MB
		estimatedMB := estimatedSize / (1024 * 1024)

		results = append(results, models.CleanCategory{
			ID:          cat.id,
			Name:        cat.name,
			Description: cat.description,
			Enabled:     true,
			EstimatedMB: estimatedMB,
		})
	}

	return results, nil
}

// ExecuteClean performs the actual cleanup
func (s *CleanService) ExecuteClean(categoryIDs []string, dryRun bool) error {
	// Load whitelist
	if err := s.loadWhitelist(); err != nil {
		return fmt.Errorf("failed to load whitelist: %w", err)
	}

	// Create a map of selected categories
	selectedCats := make(map[string]bool)
	for _, id := range categoryIDs {
		selectedCats[id] = true
	}

	totalSpaceFreed := int64(0)
	totalFilesRemoved := 0
	var errors []string
	cleanedCategories := []string{}

	totalCategories := len(categoryIDs)
	currentCategory := 0

	for _, cat := range s.categories {
		// Skip if not selected
		if !selectedCats[cat.id] {
			continue
		}

		currentCategory++
		categorySpaceFreed := int64(0)
		categoryFilesRemoved := 0

		// Emit progress
		if s.ctx != nil {
			progress := models.CleanProgress{
				Category:   cat.name,
				Message:    fmt.Sprintf("Cleaning %s...", cat.name),
				Percent:    (currentCategory * 100) / totalCategories,
				TotalFiles: totalCategories,
				FilesClean: currentCategory - 1,
			}
			runtime.EventsEmit(s.ctx, "clean:progress", progress)
		}

		// Clean each path in category
		for _, path := range cat.paths {
			// Check if path exists
			if _, err := os.Stat(path); os.IsNotExist(err) {
				continue
			}

			// Skip whitelisted paths
			if s.isWhitelisted(path) {
				continue
			}

			// Clean the path
			spaceFreed, filesRemoved, err := s.cleanPath(path, dryRun)
			if err != nil {
				errors = append(errors, fmt.Sprintf("%s: %v", cat.name, err))
				continue
			}

			categorySpaceFreed += spaceFreed
			categoryFilesRemoved += filesRemoved
		}

		if categoryFilesRemoved > 0 {
			totalSpaceFreed += categorySpaceFreed
			totalFilesRemoved += categoryFilesRemoved
			cleanedCategories = append(cleanedCategories, cat.name)
		}
	}

	// Emit final progress
	if s.ctx != nil {
		progress := models.CleanProgress{
			Category:   "Complete",
			Message:    "Cleaning complete",
			Percent:    100,
			TotalFiles: totalCategories,
			FilesClean: totalCategories,
		}
		runtime.EventsEmit(s.ctx, "clean:progress", progress)
	}

	// Emit complete event
	result := models.CleanResult{
		SpaceFreed:   totalSpaceFreed,
		FilesRemoved: totalFilesRemoved,
		Categories:   cleanedCategories,
		Errors:       errors,
	}

	if s.ctx != nil {
		runtime.EventsEmit(s.ctx, "clean:complete", result)
	}

	return nil
}

// GetWhitelist returns the current whitelist
func (s *CleanService) GetWhitelist() ([]string, error) {
	whitelistPath := filepath.Join(os.Getenv("HOME"), ".config", "mole", "whitelist")

	data, err := os.ReadFile(whitelistPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read whitelist: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	var whitelist []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			whitelist = append(whitelist, line)
		}
	}

	return whitelist, nil
}

// UpdateWhitelist updates the whitelist file
func (s *CleanService) UpdateWhitelist(paths []string) error {
	whitelistPath := filepath.Join(os.Getenv("HOME"), ".config", "mole", "whitelist")

	// Create directory if it doesn't exist
	configDir := filepath.Dir(whitelistPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	content := strings.Join(paths, "\n")
	if err := os.WriteFile(whitelistPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write whitelist: %w", err)
	}

	// Reload whitelist
	s.loadWhitelist()

	return nil
}

// Helper functions

// loadWhitelist loads the whitelist into memory
func (s *CleanService) loadWhitelist() error {
	paths, err := s.GetWhitelist()
	if err != nil {
		return err
	}

	s.whitelist = make(map[string]bool)
	for _, path := range paths {
		s.whitelist[path] = true
	}

	return nil
}

// isWhitelisted checks if a path is in the whitelist
func (s *CleanService) isWhitelisted(path string) bool {
	// Check exact match
	if s.whitelist[path] {
		return true
	}

	// Check if any parent directory is whitelisted
	for whitelistedPath := range s.whitelist {
		if strings.HasPrefix(path, whitelistedPath) {
			return true
		}
	}

	return false
}

// calculateDirSize calculates the total size of a directory
func (s *CleanService) calculateDirSize(path string) (int64, error) {
	var size int64

	err := filepath.WalkDir(path, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			// Skip permission errors and continue
			return nil
		}

		// Skip whitelisted paths
		if s.isWhitelisted(filePath) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return nil // Skip files we can't stat
			}

			// Skip files we wouldn't be able to clean
			if s.shouldSkipFile(filePath, info) {
				return nil
			}

			size += info.Size()
		}

		return nil
	})

	return size, err
}

// cleanPath removes files from a path
func (s *CleanService) cleanPath(path string, dryRun bool) (int64, int, error) {
	var spaceFreed int64
	var filesRemoved int

	info, err := os.Stat(path)
	if err != nil {
		return 0, 0, err
	}

	// If it's a file, remove it directly
	if !info.IsDir() {
		// Check if we should skip this file (e.g., not user-owned in /tmp)
		if s.shouldSkipFile(path, info) {
			return 0, 0, nil
		}

		spaceFreed = info.Size()
		filesRemoved = 1

		if !dryRun {
			if err := os.Remove(path); err != nil {
				return 0, 0, err
			}
		}

		return spaceFreed, filesRemoved, nil
	}

	// If it's a directory, walk and remove contents
	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, 0, err
	}

	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())

		// Skip whitelisted paths
		if s.isWhitelisted(entryPath) {
			continue
		}

		// Get entry info
		entryInfo, err := entry.Info()
		if err != nil {
			continue // Skip files we can't access
		}

		if entry.IsDir() {
			// Recursively clean subdirectories
			dirSize, dirFiles, err := s.cleanPath(entryPath, dryRun)
			if err != nil {
				continue // Skip directories we can't clean
			}
			spaceFreed += dirSize
			filesRemoved += dirFiles

			// Try to remove empty directory
			if !dryRun {
				os.Remove(entryPath) // Ignore error if not empty
			}
		} else {
			// Check if we should skip this file
			if s.shouldSkipFile(entryPath, entryInfo) {
				continue
			}

			// Remove file
			spaceFreed += entryInfo.Size()
			filesRemoved++

			if !dryRun {
				// Check if file is in use before deleting
				if !s.isFileInUse(entryPath) {
					os.Remove(entryPath) // Ignore errors
				}
			}
		}
	}

	return spaceFreed, filesRemoved, nil
}

// shouldSkipFile checks if a file should be skipped during cleaning
func (s *CleanService) shouldSkipFile(path string, info os.FileInfo) bool {
	// For files in /tmp or /private/var/tmp, only clean user-owned files
	if strings.HasPrefix(path, "/tmp/") || strings.HasPrefix(path, "/private/var/tmp/") {
		// Get file's UID
		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			currentUID := uint32(os.Getuid())
			if stat.Uid != currentUID {
				return true // Skip files not owned by current user
			}
		}
	}
	return false
}

// isFileInUse checks if a file is currently in use (basic check)
func (s *CleanService) isFileInUse(path string) bool {
	// Try to open the file exclusively
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_EXCL, 0666)
	if err != nil {
		// If we get a permission error, file might be in use
		if pathErr, ok := err.(*os.PathError); ok {
			if pathErr.Err == syscall.EACCES || pathErr.Err == syscall.EPERM {
				return true
			}
		}
		return false
	}
	file.Close()
	return false
}
