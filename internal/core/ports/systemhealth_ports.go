package ports

import (
	"net/http"

	"github.com/ctfrancia/maple/internal/core/domain"
)

type SystemHandler interface {
	HealthHandler(w http.ResponseWriter, r *http.Request)
	LoginHandler(w http.ResponseWriter, r *http.Request)
	CreateNewHandler(w http.ResponseWriter, r *http.Request)
}

type SystemServicer interface {
	ProcessSystemHealthRequest() domain.System
}

type SystemAdapter interface {
	GetSystemInfo() domain.System
}
