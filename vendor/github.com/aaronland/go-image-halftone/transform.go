package halftone

import (
	"context"
	"fmt"
	"image"
	"net/url"
	"strconv"

	"github.com/aaronland/go-image/transform"
)

type HalftoneTransformation struct {
	transform.Transformation
	options *HalftoneOptions
}

func init() {
	ctx := context.Background()
	transform.RegisterTransformation(ctx, "halftone", NewHalftoneTransformation)
}

func NewHalftoneTransformation(ctx context.Context, str_url string) (transform.Transformation, error) {

	u, err := url.Parse(str_url)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	opts := NewDefaultHalftoneOptions()

	q := u.Query()

	q_process := q.Get("process")
	q_scale := q.Get("scale-factor")

	if q_process != "" {
		opts.Process = q_process
	}

	if q_scale != "" {

		v, err := strconv.ParseFloat(q_scale, 64)

		if err != nil {
			return nil, err
		}

		opts.ScaleFactor = v
	}

	tr := &HalftoneTransformation{
		options: opts,
	}

	return tr, nil
}

func (tr *HalftoneTransformation) Transform(ctx context.Context, im image.Image) (image.Image, error) {
	return HalftoneImage(ctx, im, tr.options)
}
