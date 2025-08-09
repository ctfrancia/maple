package ports

// security.PasswordGeneratorDefaultLength

var PasswordGeneratorDefaultLength = 32

type SecurityAdapter interface {
	CreateSecretKey(length int) (string, error)
	Hash(password string) (string, error)
	CompareHashAndPassword(encodedHash, password string) (bool, error)
}
