package decode

import (
	"context"
	"image"
	_ "image/jpeg"
	"io"
)

// JPEGDecoder is a struct that implements the `Decoder` interface for
// decoding JPEG image.
type JPEGDecoder struct {
	Decoder
}

func init() {

	ctx := context.Background()
	RegisterDecoder(ctx, NewJPEGDecoder, "jpg", "jpeg")
}

// NewJPEGDecoder returns a new `JPEGDecoder` instance.
// 'uri' in the form of:
//
//	/path/to/image.jpg
func NewJPEGDecoder(ctx context.Context, uri string) (Decoder, error) {

	e := &JPEGDecoder{}
	return e, nil
}

// Decode will decode the body of 'r' in to an `image.Image` instance using the `image/jpeg` package.
func (e *JPEGDecoder) Decode(ctx context.Context, r io.ReadSeeker) (image.Image, string, error) {
	return image.Decode(r)
}
