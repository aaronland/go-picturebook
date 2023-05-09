package decode

import (
	"context"
	"image"
	_ "image/gif"
	"io"
)

// GIFDecoder is a struct that implements the `Decoder` interface for
// decoding GIF image.
type GIFDecoder struct {
	Decoder
}

func init() {
	ctx := context.Background()
	RegisterDecoder(ctx, NewGIFDecoder, "gif")
}

// NewGIFDecoder returns a new `GIFDecoder` instance.
// 'uri' in the form of:
//
//	/path/to/image.gif
func NewGIFDecoder(ctx context.Context, uri string) (Decoder, error) {

	e := &GIFDecoder{}
	return e, nil
}

// Decode will decode the body of 'r' in to an `image.Image` instance using the `image/gif` package.
func (e *GIFDecoder) Decode(ctx context.Context, r io.ReadSeeker) (image.Image, string, error) {
	return image.Decode(r)
}
