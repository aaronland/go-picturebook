package decode

import (
	"context"
	"image"
	_ "image/png"
	"io"
)

type PNGDecoder struct {
	Decoder
}

func init() {

	ctx := context.Background()
	err := RegisterDecoder(ctx, NewPNGDecoder, "png")

	if err != nil {
		panic(err)
	}
}

func NewPNGDecoder(ctx context.Context, uri string) (Decoder, error) {

	e := &PNGDecoder{}
	return e, nil
}

func (e *PNGDecoder) Decode(ctx context.Context, r io.ReadSeeker) (image.Image, string, error) {
	return image.Decode(r)
}
