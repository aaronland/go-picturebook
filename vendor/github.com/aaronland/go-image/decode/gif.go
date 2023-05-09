package decode

import (
	"context"
	"image"
	_ "image/gif"
	"io"
)

type GIFDecoder struct {
	Decoder
}

func init() {

	ctx := context.Background()
	err := RegisterDecoder(ctx, NewGIFDecoder, "gif")

	if err != nil {
		panic(err)
	}
}

func NewGIFDecoder(ctx context.Context, uri string) (Decoder, error) {

	e := &GIFDecoder{}
	return e, nil
}

func (e *GIFDecoder) Decode(ctx context.Context, r io.ReadSeeker) (image.Image, string, error) {
	return image.Decode(r)
}
