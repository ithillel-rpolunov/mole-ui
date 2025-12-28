package services

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"mole-wails/backend/models"
)

type OptimizeService struct {
	scriptsPath string
	ctx         context.Context
}

func NewOptimizeService(scriptsPath string) *OptimizeService {
	return &OptimizeService{
		scriptsPath: scriptsPath,
	}
}

func (s *OptimizeService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

// GetTasks returns available optimization tasks
func (s *OptimizeService) GetTasks() ([]models.OptimizationTask, error) {
	tasks := []models.OptimizationTask{
		{
			ID:           "rebuild_caches",
			Name:         "Rebuild system databases and clear caches",
			Description:  "Rebuilds system-level caches for improved performance",
			Enabled:      true,
			RequiresSudo: true,
		},
		{
			ID:           "reset_network",
			Name:         "Reset network services",
			Description:  "Resets network configuration",
			Enabled:      true,
			RequiresSudo: true,
		},
		{
			ID:           "refresh_ui",
			Name:         "Refresh Finder and Dock",
			Description:  "Restarts Finder and Dock for a fresh UI",
			Enabled:      true,
			RequiresSudo: false,
		},
		{
			ID:           "clean_logs",
			Name:         "Clean diagnostic and crash logs",
			Description:  "Removes old system logs",
			Enabled:      true,
			RequiresSudo: true,
		},
		{
			ID:           "restart_pager",
			Name:         "Remove swap files and restart dynamic pager",
			Description:  "Clears swap files and restarts paging",
			Enabled:      true,
			RequiresSudo: true,
		},
		{
			ID:           "rebuild_services",
			Name:         "Rebuild launch services and spotlight index",
			Description:  "Rebuilds system service databases",
			Enabled:      true,
			RequiresSudo: true,
		},
	}

	return tasks, nil
}

// ExecuteOptimizations runs selected optimization tasks
func (s *OptimizeService) ExecuteOptimizations(taskIDs []string) error {
	scriptPath := filepath.Join(s.scriptsPath, "bin", "optimize.sh")

	cmd := exec.Command("/bin/bash", scriptPath)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start optimization: %w", err)
	}

	// Stream progress
	scanner := bufio.NewScanner(stdout)
	totalTasks := len(taskIDs)
	currentTask := 0

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "✓") {
			currentTask++
			percent := (currentTask * 100) / totalTasks

			progress := models.OptimizeProgress{
				Message: line,
				Percent: percent,
			}

			if s.ctx != nil {
				runtime.EventsEmit(s.ctx, "optimize:progress", progress)
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("optimization failed: %w", err)
	}

	result := models.OptimizeResult{
		TasksCompleted: currentTask,
	}

	if s.ctx != nil {
		runtime.EventsEmit(s.ctx, "optimize:complete", result)
	}

	return nil
}

// GetWhitelist returns optimization tasks in whitelist
func (s *OptimizeService) GetWhitelist() ([]string, error) {
	whitelistPath := filepath.Join(os.Getenv("HOME"), ".config", "mole", "optimize_whitelist")

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

// UpdateWhitelist updates optimization whitelist
func (s *OptimizeService) UpdateWhitelist(tasks []string) error {
	whitelistPath := filepath.Join(os.Getenv("HOME"), ".config", "mole", "optimize_whitelist")

	configDir := filepath.Dir(whitelistPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	content := strings.Join(tasks, "\n")
	if err := os.WriteFile(whitelistPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write whitelist: %w", err)
	}

	return nil
}
