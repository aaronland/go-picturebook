package encode

import (
	"context"
	"image"
	"image/png"
	"io"
)

type PNGEncoder struct {
	Encoder
}

func init() {

	ctx := context.Background()
	RegisterEncoder(ctx, NewPNGEncoder, "png")
}

func NewPNGEncoder(ctx context.Context, uri string) (Encoder, error) {

	e := &PNGEncoder{}
	return e, nil
}

func (e *PNGEncoder) Encode(ctx context.Context, wr io.Writer, im image.Image) error {
	return png.Encode(wr, im)
}
