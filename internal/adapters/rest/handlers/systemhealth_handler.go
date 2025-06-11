package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ctfrancia/maple/internal/core/ports"
)

type SystemHealthHandler struct {
	system ports.SystemHealthServicer
}

func NewSystemHealthHandler(shs ports.SystemHealthServicer) *SystemHealthHandler {
	handler := &SystemHealthHandler{
		system: shs,
	}

	return handler
}

func (h *SystemHealthHandler) Handle(w http.ResponseWriter, r *http.Request) {
	sysInfo := h.system.ProcessSystemHealthRequest()
	res, err := json.Marshal(sysInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
