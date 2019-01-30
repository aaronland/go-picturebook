package util

import (
	"github.com/aaronland/go-image-tools/imaging"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"image"
	"io"
	"os"
)

func init() {
	exif.RegisterParsers(mknote.All...)
}

func NewImageWithRotationFromPath(path string) (image.Image, string, error) {

	fh, err := os.Open(path)

	if err != nil {
		return nil, "", err
	}

	defer fh.Close()

	return NewImageWithRotationFromReader(fh)
}

func NewImageWithRotationFromReader(r io.ReadSeeker) (image.Image, string, error) {

	im, format, err := DecodeImageFromReader(r)

	if err != nil {
		return nil, "", err
	}

	if format != "jpeg" {
		return im, format, nil
	}

	orientation := "1"

	_, err = r.Seek(0, 0)

	if err != nil {
		return nil, "", err
	}

	x, err := exif.Decode(r)

	if err == nil {

		o, err := x.Get(exif.Orientation)

		if err == nil {
			orientation = o.String()
		}
	}

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

	return im, format, nil
}
