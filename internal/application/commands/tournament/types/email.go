package types

type Email string

func (e Email) Validate() error {
	return nil
}
