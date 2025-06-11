package systemhealth

import "github.com/ctfrancia/maple/internal/core/domain"

type SystemHealthAdapter struct {
}

func NewSystemHealthAdapter() *SystemHealthAdapter {
	return &SystemHealthAdapter{}
}

func (sha *SystemHealthAdapter) GetSystemInfo() domain.System {
	system := domain.System{
		Version: "1.0.0",
	}

	return system
}
