// Package system provides a system adapter that is used for system information operations such as health checks
package system

import (
	"github.com/ctfrancia/maple/internal/core/domain"
	"github.com/ctfrancia/maple/internal/core/ports"
)

type SystemAdapter struct{}

func NewSystemAdapter() ports.SystemAdapter {
	return &SystemAdapter{}
}

func (sha *SystemAdapter) GetSystemInfo() domain.System {
	system := domain.System{
		Version: "1.0.0",
	}

	return system
}
