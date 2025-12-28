package services

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"mole-wails/backend/models"
)

type TouchIDService struct {
	scriptsPath string
	ctx         context.Context
}

func NewTouchIDService(scriptsPath string) *TouchIDService {
	return &TouchIDService{
		scriptsPath: scriptsPath,
	}
}

func (s *TouchIDService) SetContext(ctx context.Context) {
	s.ctx = ctx
}

// GetStatus checks if Touch ID for sudo is enabled
func (s *TouchIDService) GetStatus() (*models.TouchIDStatus, error) {
	configPath := "/etc/pam.d/sudo"
	pamModulePath := "/usr/lib/pam/pam_tid.so.2"

	// Check if pam_tid module exists
	_, err := os.Stat(pamModulePath)
	available := err == nil

	// Check if Touch ID is enabled in config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read sudoers file: %w", err)
	}

	content := string(data)
	enabled := strings.Contains(content, "pam_tid.so")

	status := &models.TouchIDStatus{
		Enabled:       enabled,
		Available:     available,
		Status:        "Disabled",
		PamModulePath: pamModulePath,
		ConfigPath:    configPath,
	}

	if enabled {
		status.Status = "Enabled"
	}

	return status, nil
}

// Enable enables Touch ID for sudo
func (s *TouchIDService) Enable() error {
	scriptPath := filepath.Join(s.scriptsPath, "bin", "touchid.sh")

	cmd := exec.Command("/bin/bash", scriptPath, "enable")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to enable Touch ID: %w, output: %s", err, string(output))
	}

	return nil
}

// Disable disables Touch ID for sudo
func (s *TouchIDService) Disable() error {
	scriptPath := filepath.Join(s.scriptsPath, "bin", "touchid.sh")

	cmd := exec.Command("/bin/bash", scriptPath, "disable")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to disable Touch ID: %w, output: %s", err, string(output))
	}

	return nil
}
