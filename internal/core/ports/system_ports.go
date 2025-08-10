// Package ports represents the interface for how the core interfaces with the rest of the world
package ports

import (
	"net/http"

	"github.com/ctfrancia/maple/internal/core/domain"
)

type SystemHandler interface {
	HealthHandler(w http.ResponseWriter, r *http.Request)
	LoginHandler(w http.ResponseWriter, r *http.Request)
	NewConsumerHandler(w http.ResponseWriter, r *http.Request)
}

type SystemServicer interface {
	ProcessSystemHealthRequest() domain.System
	Login(username, password string) (any, error)
	CreateNewConsumer(consumer domain.NewAPIConsumer) (domain.NewAPIConsumer, error)
	NewAPIConsumer(consumer domain.NewAPIConsumer) (domain.NewAPIConsumer, error)
}

type SystemAdapter interface {
	GetSystemInfo() domain.System
}

type SystemRepository interface {
	SelectByEmail(consumer domain.NewAPIConsumer) error
	CreateNewConsumer(consumer domain.NewAPIConsumer) error
}
