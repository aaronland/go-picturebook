package rotate

import (
	"context"
	"fmt"
	"image"
	"net/url"

	"github.com/aaronland/go-image/imaging"
	"github.com/aaronland/go-image/transform"
)

type RotateTransformation struct {
	transform.Transformation
	orientation string
}

func init() {

	ctx := context.Background()
	transform.RegisterTransformation(ctx, "rotate", NewRotateTransformation)
}

func NewRotateTransformation(ctx context.Context, uri string) (transform.Transformation, error) {

	parsed, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	query := parsed.Query()
	orientation := query.Get("orientation")

	if orientation == "" {
		orientation = "1"
	}

	tr := &RotateTransformation{
		orientation: orientation,
	}

	return tr, nil
}

func (tr *RotateTransformation) Transform(ctx context.Context, im image.Image) (image.Image, error) {
	return RotateImageWithOrientation(ctx, im, tr.orientation)
}

func RotateImageWithOrientation(ctx context.Context, im image.Image, orientation string) (image.Image, error) {

	switch orientation {
	case "1":
		// pass
	case "2":
		im = imaging.FlipV(im)
	case "3":
		im = imaging.Rotate180(im)
	case "4":
		im = imaging.Rotate180(imaging.FlipV(im))
	case "5":
		im = imaging.Rotate270(imaging.FlipV(im))
	case "6":
		im = imaging.Rotate270(im)
	case "7":
		im = imaging.Rotate90(imaging.FlipV(im))
	case "8":
		im = imaging.Rotate90(im)
	}

	return im, nil
}

func RotateImageWithDegrees(ctx context.Context, im image.Image, degrees float64) (image.Image, error) {

	// See also: https://github.com/anthonynsimon/bild#rotate
	// The problem is that bild doesn't rotate the "canvas" just the image

	switch degrees {
	case 90.0:
		im = imaging.Rotate90(im)
	case 180.0:
		im = imaging.Rotate180(im)
	case 270.0:
		im = imaging.Rotate270(im)
	default:
		return nil, fmt.Errorf("Unsupported value, %f", degrees)
	}

	return im, nil
}
