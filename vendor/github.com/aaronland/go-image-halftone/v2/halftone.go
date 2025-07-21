package halftone

// https://maxhalford.github.io/blog/halftoning-1/

import (
	"context"
	"fmt"
	"image"

	"github.com/MaxHalford/halfgone"
	"github.com/nfnt/resize"
)

type HalftoneOptions struct {
	Process     string
	ScaleFactor float64
}

func NewDefaultHalftoneOptions() *HalftoneOptions {

	opts := &HalftoneOptions{
		Process:     "atkinson",
		ScaleFactor: 2.0,
	}

	return opts
}

func HalftoneImage(ctx context.Context, im image.Image, opts *HalftoneOptions) (image.Image, error) {

	dims := im.Bounds()
	w := uint(dims.Max.X)
	h := uint(dims.Max.Y)

	scale_w := uint(float64(w) / opts.ScaleFactor)
	scale_h := uint(float64(h) / opts.ScaleFactor)

	thumb := resize.Thumbnail(scale_w, scale_h, im, resize.Lanczos3)
	grey := halfgone.ImageToGray(thumb)

	switch opts.Process {
	case "atkinson":
		grey = halfgone.AtkinsonDitherer{}.Apply(grey)
	case "threshold":
		grey = halfgone.ThresholdDitherer{Threshold: 127}.Apply(grey)
	default:
		return nil, fmt.Errorf("Invalid or unsupported process, %s", opts.Process)
	}

	dither := resize.Resize(w, h, grey, resize.NearestNeighbor)
	grey = halfgone.ImageToGray(dither)

	return grey, nil
}
