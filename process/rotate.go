package process

// update to use go-image-rotate

import (
	"context"
	"github.com/aaronland/go-image-tools/util"
	"github.com/microcosm-cc/exifutil"
	"github.com/rwcarlsen/goexif/exif"
	"net/url"
	"os"
	"path/filepath"
)

func init() {

	ctx := context.Background()
	err := RegisterProcess(ctx, "rotate", NewRotateProcess)

	if err != nil {
		panic(err)
	}
}

type RotateProcess struct {
	Process
}

func NewRotateProcess(ctx context.Context, uri string) (Process, error) {

	_, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	f := &RotateProcess{}

	return f, nil
}

func (f *RotateProcess) Transform(ctx context.Context, path string) (string, error) {

	ext := filepath.Ext(path)

	if ext != ".jpg" && ext != ".jpeg" {
		return "", nil
	}

	fh, err := os.Open(path)

	if err != nil {
		return "", err
	}

	defer fh.Close()

	x, err := exif.Decode(fh)

	if err != nil {
		return "", err
	}

	tag, err := x.Get(exif.Orientation)

	if err != nil {
		return "", nil
	}

	// log.Println(path, tag)

	orientation, err := tag.Int64(0)

	if err != nil {
		return "", nil
	}

	if orientation == 1 {
		return "", nil
	}

	im, format, err := util.DecodeImage(path)

	if err != nil {
		return "", err
	}

	angle, _, _ := exifutil.ProcessOrientation(orientation)
	rotated := exifutil.Rotate(im, angle)

	return util.EncodeTempImage(rotated, format)
}