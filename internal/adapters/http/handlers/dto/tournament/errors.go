package dto

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Errors map[string]string `json:"errors"`
}

func (ve ValidationError) Error() string {
	if len(ve.Errors) == 0 {
		return "validation failed"
	}

	var messages []string
	for field, msg := range ve.Errors {
		messages = append(messages, fmt.Sprintf("%s: %s", field, msg))
	}

	return fmt.Sprintf("validation failed: %s", strings.Join(messages, ", "))
}
