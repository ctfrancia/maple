package systemhealth

import "github.com/ctfrancia/maple/internal/core/domain"

func GetSystemInfo() domain.System {
	system := domain.System{
		Version: "1.0.0",
	}

	return system
}
