package decode

import (
	"context"
	"image"
	_ "image/jpeg"
	"io"
)

type JPEGDecoder struct {
	Decoder
}

func init() {

	ctx := context.Background()
	err := RegisterDecoder(ctx, NewJPEGDecoder, "jpg", "jpeg")

	if err != nil {
		panic(err)
	}
}

func NewJPEGDecoder(ctx context.Context, uri string) (Decoder, error) {

	e := &JPEGDecoder{}
	return e, nil
}

func (e *JPEGDecoder) Decode(ctx context.Context, r io.ReadSeeker) (image.Image, string, error) {
	return image.Decode(r)
}
