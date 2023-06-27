package encode

import (
	"context"
	"image"
	"image/gif"
	"io"
)

// GIFEncoder is a struct that implements the `Encoder` interface for
// encoding GIF images.
type GIFEncoder struct {
	Encoder
	options *gif.Options
}

func init() {

	ctx := context.Background()
	RegisterEncoder(ctx, NewGIFEncoder, "gif")
}

// NewGIFEncoder returns a new `GIFEncoder` instance.
// 'uri' in the form of:
//
//	/path/to/image.gif
func NewGIFEncoder(ctx context.Context, uri string) (Encoder, error) {

	opts := &gif.Options{}

	e := &GIFEncoder{
		options: opts,
	}

	return e, nil
}

// Encode will encode 'im' using the `image/gif` package and write the results to 'wr'
func (e *GIFEncoder) Encode(ctx context.Context, wr io.Writer, im image.Image) error {
	return gif.Encode(wr, im, e.options)
}
