// Package repository provides a set of functions for interacting with the database
package repository

import (
	"github.com/ctfrancia/maple/internal/core/domain"
)

type SystemRepository struct {
	db any // TODO: add the database connection it will be a gorm.DB
}

func NewSystemRepository() *SystemRepository {
	return &SystemRepository{}
}

func (sr *SystemRepository) SelectByEmail(consumer domain.NewAPIConsumer) error {
	return nil
}

func (sr *SystemRepository) CreateNewConsumer(consumer domain.NewAPIConsumer) error {
	return nil
}
