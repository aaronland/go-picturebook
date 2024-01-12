package contour

import (
	"context"
	"fmt"
	"image"
	"net/url"
	"strconv"

	"github.com/aaronland/go-image/transform"
)

// ContourTransformation is a struct that implements the `Transformation` interface for
// replacing the contents of an image with contours derived from its pixel values.
type ContourTransformation struct {
	transform.Transformation
	n     int
	scale float64
}

func init() {
	ctx := context.Background()
	transform.RegisterTransformation(ctx, "contour", NewContourTransformation)
}

// NewContourWriter returns a new `ContourTransformation` instance configure by 'uri'
// in the form of:
//
//	contour://?n={N}&scale={SCALE}
func NewContourTransformation(ctx context.Context, uri string) (transform.Transformation, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()

	n := 12
	scale := 1.0

	q_n := q.Get("n")
	q_scale := q.Get("scale")

	if q_n != "" {

		v, err := strconv.Atoi(q_n)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse ?n= parameter, %w", err)
		}

		n = v
	}

	if q_scale != "" {

		v, err := strconv.ParseFloat(q_scale, 64)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse ?scale= parameter, %w", err)
		}

		scale = v
	}

	tr := &ContourTransformation{
		n:     n,
		scale: scale,
	}
	return tr, nil
}

// Transform will derive the contours of 'im' and draw them to a new `image.Image` instance.
func (tr *ContourTransformation) Transform(ctx context.Context, im image.Image) (image.Image, error) {
	return ContourImage(ctx, im, tr.n, tr.scale)
}
