package transform

import (
	"context"
	"image"
)

func init() {
	ctx := context.Background()
	RegisterTransformation(ctx, "null", NewNullTransformation)
}

type NullTransformation struct {
	Transformation
}

func NewNullTransformation(ctx context.Context, uri string) (Transformation, error) {

	tr := &NullTransformation{}
	return tr, nil
}

func (tr *NullTransformation) Transform(ctx context.Context, im image.Image) (image.Image, error) {
	return im, nil
}
