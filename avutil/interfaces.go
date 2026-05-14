package avutil

// IDictionary is the interface a libav AVDictionary exposes. *Dictionary
// implements it directly; tests may substitute fakes.
type IDictionary interface {
	Set(key, value string) error
	Free()
}
