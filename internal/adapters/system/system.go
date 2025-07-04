// Package system provides a system adapter that is used for system information operations such as health checks
package system

import "github.com/ctfrancia/maple/internal/core/domain"

type SystemAdapter struct{}

func NewSystemAdapter() *SystemAdapter {
	return &SystemAdapter{}
}

func (sha *SystemAdapter) GetSystemInfo() domain.System {
	system := domain.System{
		Version: "1.0.0",
	}

	return system
}
