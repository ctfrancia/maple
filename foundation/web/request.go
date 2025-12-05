package web

import (
	"fmt"
	"io"
	"net/http"
)

// Param returns the web call parameters from the request.
func Param(r *http.Request, key string) string {
	return r.PathValue(key)
}

// Decoder represents data that can be decoded
type Decoder interface {
	Decode(data []byte) error
}

type validator interface {
	Validate() error
}

func Decode(r *http.Request, d Decoder) error {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("request: unable to read payload: %w", err)
	}

	if err := d.Decode(data); err != nil {
		return fmt.Errorf("request: decode: %w", err)
	}

	if v, ok := d.(validator); ok {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	return nil
}
