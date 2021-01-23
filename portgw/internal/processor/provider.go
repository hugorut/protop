package processor

import "errors"

// FileProcessor is an interface for a concrete type of processor
// that takes a file location and kicks off the file processing logic.
type FileProcessor interface {
	Process(string) (int, error)
}

var (
	ErrorInvalidProvider = errors.New("provider not found")
)

// Provider is a wrapper around a map, providing a very naive
// lookup container satisfying the Provider interface.
type Provider map[string]FileProcessor


// Get provides
func (receiver Provider) Get(name string) (FileProcessor, error)  {
	if _, ok := receiver[name]; !ok {
		return nil, ErrorInvalidProvider
	}

	return receiver[name], nil
}
