package encode

import (
	"context"
	"image"
	"image/png"
	"io"
)

// PNGEncoder is a struct that implements the `Encoder` interface for
// encoding PNG images.
type PNGEncoder struct {
	Encoder
}

func init() {

	ctx := context.Background()
	RegisterEncoder(ctx, NewPNGEncoder, "png")
}

// NewPNGEncoder returns a new `PNGEncoder` instance.
// 'uri' in the form of:
//
//	/path/to/image.jpg
func NewPNGEncoder(ctx context.Context, uri string) (Encoder, error) {

	e := &PNGEncoder{}
	return e, nil
}

// Encode will encode 'im' using the `image/png` package and write the results to 'wr'
func (e *PNGEncoder) Encode(ctx context.Context, wr io.Writer, im image.Image) error {
	return png.Encode(wr, im)
}
