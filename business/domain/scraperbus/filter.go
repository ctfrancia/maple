package scraperbus

import "time"

// QueryFilter holds the available fields a query can be filtered on.
type QueryFilter struct {
	ID         *string
	Federation *string
	Status     *string
	StartTime  *time.Time
	EndTime    *time.Time
}
