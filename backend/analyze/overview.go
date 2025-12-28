//go:build darwin

package analyze

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// createOverviewEntries generates the standard overview entries for disk analysis
func createOverviewEntries() []dirEntry {
	home := os.Getenv("HOME")
	entries := []dirEntry{}

	// Separate Home and ~/Library for better visibility and performance
	// Home excludes Library to avoid duplicate scanning
	if home != "" {
		entries = append(entries, dirEntry{Name: "Home", Path: home, IsDir: true, Size: -1})

		// Add ~/Library separately so users can see app data usage
		userLibrary := filepath.Join(home, "Library")
		if _, err := os.Stat(userLibrary); err == nil {
			entries = append(entries, dirEntry{Name: "App Library", Path: userLibrary, IsDir: true, Size: -1})
		}
	}

	entries = append(entries,
		dirEntry{Name: "Applications", Path: "/Applications", IsDir: true, Size: -1},
		dirEntry{Name: "System Library", Path: "/Library", IsDir: true, Size: -1},
	)

	// Add Volumes shortcut only when it contains real mounted folders (e.g., external disks)
	if hasUsefulVolumeMounts("/Volumes") {
		entries = append(entries, dirEntry{Name: "Volumes", Path: "/Volumes", IsDir: true, Size: -1})
	}

	return entries
}

// hasUsefulVolumeMounts checks if /Volumes contains useful mounted volumes
func hasUsefulVolumeMounts(path string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		name := entry.Name()
		// Skip hidden control entries for Spotlight/TimeMachine etc.
		if strings.HasPrefix(name, ".") {
			continue
		}

		info, err := os.Lstat(filepath.Join(path, name))
		if err != nil {
			continue
		}
		if info.Mode()&fs.ModeSymlink != 0 {
			continue // Ignore the synthetic MacintoshHD link
		}
		if info.IsDir() {
			return true
		}
	}
	return false
}
