package decode

import (
	"context"
	"image"
	_ "image/png"
	"io"
)

// PNGDecoder is a struct that implements the `Decoder` interface for
// decoding PNG image.
type PNGDecoder struct {
	Decoder
}

func init() {

	ctx := context.Background()
	RegisterDecoder(ctx, NewPNGDecoder, "png")
}

// NewPNGDecoder returns a new `PNGDecoder` instance.
// 'uri' in the form of:
//
//	/path/to/image.png
func NewPNGDecoder(ctx context.Context, uri string) (Decoder, error) {

	e := &PNGDecoder{}
	return e, nil
}

// Decode will decode the body of 'r' in to an `image.Image` instance using the `image/png` package.
func (e *PNGDecoder) Decode(ctx context.Context, r io.ReadSeeker) (image.Image, string, error) {
	return image.Decode(r)
}
