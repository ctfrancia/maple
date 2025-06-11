package services

import (
	"github.com/ctfrancia/maple/internal/core/domain"
	"github.com/ctfrancia/maple/internal/core/ports"
)

type SystemHealthServicer struct {
	shAdapter ports.SystemHealthAdapter
}

func NewSystemHealthServicer(aha ports.SystemHealthAdapter) *SystemHealthServicer {
	return &SystemHealthServicer{}
}

func (shs *SystemHealthServicer) ProcessSystemHealthRequest() domain.System {
	return shs.shAdapter.GetSystemInfo()
}
