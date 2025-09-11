package types

type PaymentType string

const (
	PaymentTypeMonetary PaymentType = "monetary" // with money
	PaymentTypePhysical PaymentType = "physical" // with physical goods (e.g. book/lesson/etc.)
	PaymentTypeOther    PaymentType = "other"    // with other
)

// Payment represents the payment information for the tournament
type Payment struct {
	Place  *int         `json:"place,omitempty"`  // 1st, 2nd, etc.
	Amount *int64       `json:"amount,omitempty"` // if type is monetary (in cents)
	Type   *PaymentType `json:"type,omitempty"`
}

func (p *Payment) Validate() error {
	return nil
}
