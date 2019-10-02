package util

import (
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

	if err != nil {
		return im, format, nil
	}

	o, err := x.Get(exif.Orientation)

	if err != nil {
		return im, format, nil
	}

	orientation = o.String()

	im, err = RotateWithOrientation(im, orientation)

	if err != nil {
		return nil, "", err
	}

	return im, format, nil
}
