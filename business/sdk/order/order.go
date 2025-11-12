package order

import (
	"fmt"
	"strings"
)

const (
	ASC  = "ASC"
	DESC = "DESC"
)

var directions = map[string]string{
	ASC:  "ASC",
	DESC: "DESC",
}

// By represents a field used to order and by which direction.
type By struct {
	Field     string
	Direction string
}

// NewBy creates a new By field.
func NewBy(field string, direction string) By {
	// default if the direction is no found.
	if _, ok := directions[direction]; !ok {
		return By{
			Field:     field,
			Direction: ASC,
		}
	}

	return By{
		Field:     field,
		Direction: direction,
	}
}

// Parse contructs a By value by parsing a string in the form of
// "field,direction" ie "user_id,ASC"
func Parse(fieldMappings map[string]string, orderBy string, defaultOrder By) (By, error) {
	if orderBy == "" {
		return defaultOrder, nil
	}

	orderParts := strings.Split(orderBy, ",")

	orgFieldName := strings.TrimSpace(orderParts[0])
	fieldName, ok := fieldMappings[orgFieldName]
	if !ok {
		return By{}, fmt.Errorf("unknown order: %s", orgFieldName)
	}

	switch len(orderParts) {
	case 1:
		return NewBy(fieldName, ASC), nil

	case 2:
		direction := strings.TrimSpace(orderParts[1])
		if _, ok := directions[direction]; !ok {
			return By{}, fmt.Errorf("unknown direction: %s", direction)
		}

		return NewBy(fieldName, direction), nil

	default:
		return By{}, fmt.Errorf("invalid order: %s", orderBy)
	}
}
