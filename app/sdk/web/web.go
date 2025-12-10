package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ctfrancia/maple/app/sdk/errs"
	"github.com/ctfrancia/maple/foundation/logger"
)

// Response wraps common response patterns
type Response struct {
	Data   any    `json:"data,omitempty"`
	Error  string `json:"error,omitempty"`
	Errors any    `json:"errors,omitempty"`
	Status int    `json:"-"`
}

// HandlerFunc is the custom handler that returns a Response and error
type HandlerFunc func(w http.ResponseWriter, r *http.Request) (*Response, error)

func Wrap(log *logger.Logger, h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		res, err := h(w, r)

		if err != nil {
			log.Error(ctx, "handler error", "path", r.URL.Path, "error", err)
			status := http.StatusInternalServerError
			if res != nil && res.Status != 0 {
				status = res.Status
			}

			// Check if it's a validation error
			var validationErr *errs.ValidationError
			if errors.As(err, &validationErr) {
				status = http.StatusBadRequest
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(status)
				json.NewEncoder(w).Encode(Response{
					Error:  "validation failed",
					Errors: validationErr.Fields,
				})
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)

			json.NewEncoder(w).Encode(Response{Error: err.Error()})
			return
		}

		if res != nil {
			log.Info(ctx, "handler success", "path", r.URL.Path, "status", res.Status)
			status := res.Status
			if status == 0 {
				status = http.StatusOK
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(status)
			if res.Data != nil {
				json.NewEncoder(w).Encode(res.Data)
			}
		}
	}
}
