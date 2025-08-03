package ports

type ValidatorServicer interface {
	Valid() bool
	AddError(key, message string)
	Check(ok bool, key, message string)
	In(key string, permittedValues ...string) bool
	ReturnErrors() map[string]string
}
