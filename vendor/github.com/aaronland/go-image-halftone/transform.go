package halftone

import (
	"context"
	"github.com/aaronland/go-image-transform"
	"image"
	"net/url"
	"strconv"
)

type HalftoneTransformation struct {
	transform.Transformation
	options *HalftoneOptions
}

func init() {

	ctx := context.Background()
	err := transform.RegisterTransformation(ctx, "Halftone", NewHalftoneTransformation)

	if err != nil {
		panic(err)
	}
}

func NewHalftoneTransformation(ctx context.Context, str_url string) (transform.Transformation, error) {

	parsed, err := url.Parse(str_url)

	if err != nil {
		return nil, err
	}

	opts := NewDefaultHalftoneOptions()

	query := parsed.Query()

	mode := query.Get("mode")
	str_scale := query.Get("scale-factor")

	if mode != "" {
		opts.Mode = mode
	}

	if str_scale != "" {

		scale, err := strconv.ParseFloat(str_scale, 64)

		if err != nil {
			return nil, err
		}

		opts.ScaleFactor = scale
	}

	tr := &HalftoneTransformation{
		options: opts,
	}

	return tr, nil
}

func (tr *HalftoneTransformation) Transform(ctx context.Context, im image.Image) (image.Image, error) {
	return HalftoneImage(ctx, im, tr.options)
}
