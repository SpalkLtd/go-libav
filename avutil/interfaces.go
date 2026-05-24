package avutil

// IDictionary is the interface implemented by *Dictionary.
type IDictionary interface {
	Set(key, value string) error
	Free()
}
