package services

import (
	"github.com/ctfrancia/maple/internal/core/domain"
	"github.com/ctfrancia/maple/internal/core/ports"
)

type SystemHealthServicer struct {
	sAdapter ports.SystemAdapter
}

func NewSystemHealthServicer(sa ports.SystemAdapter) *SystemHealthServicer {
	return &SystemHealthServicer{
		sAdapter: sa,
	}
}

func (shs *SystemHealthServicer) ProcessSystemHealthRequest() domain.System {
	return shs.sAdapter.GetSystemInfo()
}
