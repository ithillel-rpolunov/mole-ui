package services

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"mole-wails/backend/models"
)

type UninstallService struct {
	scriptsPath string
	ctx         context.Context
}

func NewUninstallService(scriptsPath string) *UninstallService {
	return &UninstallService{
		scriptsPath: scriptsPath,
	}
}

func (s *UninstallService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

// ScanApplications scans installed applications
func (s *UninstallService) ScanApplications(forceRescan bool) ([]models.Application, error) {
	scriptPath := filepath.Join(s.scriptsPath, "bin", "uninstall.sh")

	args := []string{scriptPath}
	if forceRescan {
		args = append(args, "--force-rescan")
	}

	// We need to run the script and parse its output
	// For now, we'll scan applications directories directly
	var apps []models.Application

	appDirs := []string{
		"/Applications",
		filepath.Join(os.Getenv("HOME"), "Applications"),
	}

	for _, dir := range appDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !strings.HasSuffix(entry.Name(), ".app") {
				continue
			}

			appPath := filepath.Join(dir, entry.Name())
			info, err := entry.Info()
			if err != nil {
				continue
			}

			// Get app size
			size, _ := s.getDirSize(appPath)

			// Calculate age
			age := s.calculateAge(info.ModTime())

			apps = append(apps, models.Application{
				Name:         strings.TrimSuffix(entry.Name(), ".app"),
				BundleID:     s.getBundleID(appPath),
				Path:         appPath,
				Size:         size,
				LastModified: info.ModTime(),
				Age:          age,
			})
		}
	}

	return apps, nil
}

// UninstallApps uninstalls selected applications
func (s *UninstallService) UninstallApps(apps []string) error {
	scriptPath := filepath.Join(s.scriptsPath, "bin", "uninstall.sh")

	for i, app := range apps {
		// Execute uninstall script for each app
		cmd := exec.Command("/bin/bash", scriptPath, app)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return fmt.Errorf("failed to create stdout pipe: %w", err)
		}

		if err := cmd.Start(); err != nil {
			return fmt.Errorf("failed to start uninstall: %w", err)
		}

		// Stream progress
		scanner := bufio.NewScanner(stdout)
		filesRemoved := 0

		for scanner.Scan() {
			line := scanner.Text()

			if strings.Contains(line, "✓ Removed") {
				filesRemoved++
			}

			percent := ((i + 1) * 100) / len(apps)

			progress := models.UninstallProgress{
				App:          app,
				Message:      line,
				Percent:      percent,
				FilesRemoved: filesRemoved,
			}

			if s.ctx != nil {
				runtime.EventsEmit(s.ctx, "uninstall:progress", progress)
			}
		}

		if err := cmd.Wait(); err != nil {
			return fmt.Errorf("uninstall failed for %s: %w", app, err)
		}
	}

	result := models.UninstallResult{
		AppsRemoved: len(apps),
	}

	if s.ctx != nil {
		runtime.EventsEmit(s.ctx, "uninstall:complete", result)
	}

	return nil
}

// GetRelatedFiles finds all files related to an application
func (s *UninstallService) GetRelatedFiles(bundleID string) ([]string, error) {
	// Search common locations for app-related files
	searchPaths := []string{
		filepath.Join(os.Getenv("HOME"), "Library", "Application Support"),
		filepath.Join(os.Getenv("HOME"), "Library", "Caches"),
		filepath.Join(os.Getenv("HOME"), "Library", "Preferences"),
		filepath.Join(os.Getenv("HOME"), "Library", "Logs"),
		filepath.Join(os.Getenv("HOME"), "Library", "Cookies"),
	}

	var relatedFiles []string

	for _, searchPath := range searchPaths {
		entries, err := os.ReadDir(searchPath)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if strings.Contains(entry.Name(), bundleID) {
				relatedFiles = append(relatedFiles, filepath.Join(searchPath, entry.Name()))
			}
		}
	}

	return relatedFiles, nil
}

// Helper functions

func (s *UninstallService) getDirSize(path string) (int64, error) {
	var size int64

	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size, err
}

func (s *UninstallService) getBundleID(appPath string) string {
	plistPath := filepath.Join(appPath, "Contents", "Info.plist")

	// Use plutil to read bundle ID
	cmd := exec.Command("plutil", "-extract", "CFBundleIdentifier", "raw", plistPath)
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	return strings.TrimSpace(string(output))
}

func (s *UninstallService) calculateAge(modTime time.Time) string {
	duration := time.Since(modTime)
	days := int(duration.Hours() / 24)

	if days < 7 {
		return "Recent"
	} else if days < 30 {
		return "< 1 month"
	} else if days < 90 {
		return "< 3 months"
	} else if days < 180 {
		return "< 6 months"
	} else {
		return "Old"
	}
}
