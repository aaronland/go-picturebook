package util

import (
	"github.com/aaronland/go-image-tools/imaging"
	"image"
)

func RotateWithOrientation(im image.Image, orientation string) (image.Image, error) {

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
