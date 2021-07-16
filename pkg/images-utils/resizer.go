// Package utils provides images manipulation utilities.
//
// At this moment, it knows how to produce a thumbnail from an input image
// (gif, jpeg or png).
package utils

import (
	"bytes"
	"image"
	"io"

	// registers supported formats
	_ "image/gif"
	_ "image/jpeg"
	"image/png"

	"github.com/nfnt/resize"
)

// Resize an input image (jpeg, png, gif) to width x height, as a PNG thumbnail
func Resize(img io.Reader, w, h uint) ([]byte, error) {
	decoded, _, err := image.Decode(img)
	if err != nil {
		return nil, err
	}

	thumb := resize.Thumbnail(w, h, decoded, resize.NearestNeighbor)

	var buf bytes.Buffer
	err = png.Encode(&buf, thumb)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
