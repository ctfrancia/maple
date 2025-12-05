// Package fide contains types related to FIDE.
package fide

type Fide struct {
	eventID string // has to be a valid FIDE event ID
	isValid bool   // is valid for FIDE ELO?
	remarks string
}

func newFide(eventID string, isValid bool, remarks string) Fide {
	return Fide{eventID, isValid, remarks}
}

func (f Fide) String() string {
	return f.eventID
}

func (f Fide) Valid() bool {
	return f.isValid
}

func (f Fide) Remarks() string {
	return f.remarks
}
