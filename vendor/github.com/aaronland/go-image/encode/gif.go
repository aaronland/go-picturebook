package encode

import (
	"context"
	"image"
	"image/gif"
	"io"
)

type GIFEncoder struct {
	Encoder
	options *gif.Options
}

func init() {

	ctx := context.Background()
	err := RegisterEncoder(ctx, NewGIFEncoder, "gif")

	if err != nil {
		panic(err)
	}
}

func NewGIFEncoder(ctx context.Context, uri string) (Encoder, error) {

	opts := &gif.Options{}

	e := &GIFEncoder{
		options: opts,
	}

	return e, nil
}

func (e *GIFEncoder) Encode(ctx context.Context, wr io.Writer, im image.Image) error {
	return gif.Encode(wr, im, e.options)
}
