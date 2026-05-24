package avutil

type IDictionary interface {
	Set(key, value string) error
	Free()
}
