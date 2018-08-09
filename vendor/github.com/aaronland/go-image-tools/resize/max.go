package resize

import (
	"fmt"
	"github.com/aaronland/go-image-tools/util"
	"github.com/facebookgo/atomicfile"
	nfnt_resize "github.com/nfnt/resize"
	"image"
	"math"
	"path/filepath"
	"strings"
)

func ResizeMax(path string, max int) (string, error) {

	abs_path, err := filepath.Abs(path)

	if err != nil {
		return "", err
	}

	im, err := ResizeMaxFromPath(path, max)

	if err != nil {
		return "", err
	}

	root := filepath.Dir(abs_path)
	fname := filepath.Base(abs_path)

	ext := filepath.Ext(abs_path)

	fname = strings.Replace(fname, ext, "", -1)
	fname = fmt.Sprintf("%s-%d%s", fname, max, ext)

	new_path := filepath.Join(root, fname)

	fh, err := atomicfile.New(new_path, 0644)

	if err != nil {
		return "", err
	}

	format := strings.TrimLeft(ext, ".")

	err = util.EncodeImage(im, format, fh)

	if err != nil {
		fh.Abort()
		return "", err
	}

	err = fh.Close()

	if err != nil {
		return "", err
	}

	return new_path, nil
}

func ResizeMaxFromPath(path string, max int) (image.Image, error) {

	abs_path, err := filepath.Abs(path)

	if err != nil {
		return nil, err
	}

	im, _, err := util.DecodeImage(abs_path)

	if err != nil {
		return nil, err
	}

	return ResizeMaxFromReader(im, max)
}

func ResizeMaxFromReader(im image.Image, max int) (image.Image, error) {

	// calculating w,h is probably unnecessary since we're
	// calling resize.Thumbnail but it will do for now...
	// (20180708/thisisaaronland)

	bounds := im.Bounds()
	dims := bounds.Max

	width := dims.X
	height := dims.Y

	ratio_w := float64(max) / float64(width)
	ratio_h := float64(max) / float64(height)

	ratio := math.Min(ratio_w, ratio_h)

	w := uint(float64(width) * ratio)
	h := uint(float64(height) * ratio)

	sm := nfnt_resize.Thumbnail(w, h, im, nfnt_resize.Lanczos3)

	return sm, nil
}
