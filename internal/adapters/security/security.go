package security

type SecurityAdapter struct{}

func NewSecurityAdapter() *SecurityAdapter {
	return &SecurityAdapter{}
}

func (sa *SecurityAdapter) CreateSecretKey(length int) (string, error) {
	return "", nil
}

func (sa *SecurityAdapter) Hash(password string) (string, error) {
	return "", nil
}
