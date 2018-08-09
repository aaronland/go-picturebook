package util

import (
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"os"
)

func EncodeTempImage(im image.Image, format string) (string, error) {

	fh, err := ioutil.TempFile("", "picturebook-")

	if err != nil {
		return "", err
	}

	defer fh.Close()

	err = EncodeImage(im, format, fh)

	if err != nil {
		return "", err
	}

	// see what's going on here - this (appending the format
	// extension) is necessary because without it fpdf.GetImageInfo
	// gets confused and FREAKS out triggering fatal errors
	// along the way... oh well (20171125/thisisaaronland)

	fname := fh.Name()
	fh.Close()

	fq_fname := fmt.Sprintf("%s.%s", fname, format)

	err = os.Rename(fname, fq_fname)

	if err != nil {
		return "", err
	}

	return fq_fname, nil
}

func EncodeImage(im image.Image, format string, wr io.Writer) error {

	var err error

	if format == "jpg" {
		format = "jpeg"
	}

	switch format {
	case "jpeg":
		opts := jpeg.Options{Quality: 100}
		err = jpeg.Encode(wr, im, &opts)
	case "png":
		err = png.Encode(wr, im)
	case "gif":
		opts := gif.Options{}
		err = gif.Encode(wr, im, &opts)
	default:
		err = errors.New("Invalid or unsupported format")
	}

	return err
}
