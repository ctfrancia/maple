package commands

type RegistrationStatus string

const (
	RegistrationStatusOpen   RegistrationStatus = "open"
	RegistrationStatusClosed RegistrationStatus = "closed"
)

type PaymentType string

const (
	PaymentTypeMonetary PaymentType = "monetary" // money
	PaymentTypePhysical PaymentType = "physical" // e.g. book/lesson/etc.
	PaymentTypeOther    PaymentType = "other"
)
