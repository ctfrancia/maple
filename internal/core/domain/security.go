package domain

import "errors"

// PasswordGeneratorDefaultLength is the default length of the password generator
// the default is 8 characters for dev, and will be longer for production
const PasswordGeneratorDefaultLength = 8

type A2params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)
