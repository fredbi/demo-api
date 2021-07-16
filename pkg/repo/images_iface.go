package repo

import (
	"errors"
	"io"
)

// ErrNotFound indicates that the key was not found for the requested item
var ErrNotFound = errors.New("image not found")

// ImagesRepo knows how to manage images on a persistent store.
//
// NOTE: defining this as an interface allows for different implementations
// of the underlying persistent store. In this demo implementation we use a plain
// KV store.
//
// Using an interface also allows for generating a mock, which is most convenient
// to test the API layer independently. We don't add the generated mock here, for the
// sake of brevity.
type ImagesRepo interface {
	// Get an image
	Get(string) (io.ReadCloser, error)

	// List all images
	List() ([]Image, error)

	// Create a new image. Fails if it already exists.
	Create(string, io.ReadCloser) error

	// Update an existing image. Fails if it does not already exist.
	Update(string, io.ReadCloser) error

	// Delete an image
	Delete(string) error
}

// Image holds metadata about an image, including its PNG thumbnail
type Image struct {
	// key is the original file name
	Key string

	// base64 encoded thumbnail data
	Thumb string
}
