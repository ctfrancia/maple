// Package arbiter contains types related to arbiters.
// https://arbiters.fide.com/wp-content/uploads/Database/List_of_Arbiters.pdf
/*
Arbiters' Titles
-- IA - International Arbiter
-- FA - FIDE Arbiter
-- NA - National Arbiter
Arbiters' Categories
-- A - A Category Arbiter (only IA)
-- B - B Category Arbiter (only IA)
-- C - C Category Arbiter (IA and FA)
-- D - D Category Arbiter (IA and FA)
Arbiters' License levels
-- IA-A - International Arbiter - A License
-- IA-B - International Arbiter - B License
-- IA-C - International Arbiter - C License
-- IA-D - International Arbiter - D License
-- FA-C - FIDE Arbiter - C License
-- FA-D - FIDE Arbiter - D License
-- NA - National Arbiter License
-- No - No License
Flag
-- i - Inactive Arbiter
*/
package arbiter

type Arbiter struct {
	value string
}

// String returns the value of the arbiter.
func (a Arbiter) String() string {
	return a.value
}

// MarshalText implements the encoding.TextMarshaler interface.
func (a Arbiter) MarshalText() ([]byte, error) {
	return []byte(a.value), nil
}
