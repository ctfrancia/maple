package ports

import (
	"net/http"

	"github.com/ctfrancia/maple/internal/core/domain"
)

type SystemHealthHandler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

type SystemHealthServicer interface {
	ProcessSystemHealthRequest() domain.System
}

type SystemHealthAdapter interface {
	GetSystemInfo() domain.System
}
