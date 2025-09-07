package types

type PaymentType int

const (
	PaymentTypeMonetary PaymentType = iota
	PaymentTypePhysical
	PaymentTypeOther
)

// Payment represents the payment information for the tournament
type Payment struct {
	Place  int         `json:"place"`  // 1st, 2nd, etc.
	Amount int64       `json:"amount"` // if type is monetary
	Type   PaymentType `json:"type"`
}
