package process

// update to use go-image-rotate

import (
	"bytes"
	"context"
	"github.com/aaronland/go-image-tools/util"
	"github.com/aaronland/go-picturebook/tempfile"
	"github.com/microcosm-cc/exifutil"
	"github.com/rwcarlsen/goexif/exif"
	"gocloud.dev/blob"
	"io/ioutil"
	_ "log"
	"net/url"
	"path/filepath"
	"strings"
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

func (f *RotateProcess) Transform(ctx context.Context, bucket *blob.Bucket, path string) (string, error) {

	ext := filepath.Ext(path)
	ext = strings.ToLower(ext)

	if ext != ".jpg" && ext != ".jpeg" {
		return "", nil
	}

	fh, err := bucket.NewReader(ctx, path, nil)

	if err != nil {
		return "", err
	}

	defer fh.Close()

	body, err := ioutil.ReadAll(fh)

	if err != nil {
		return "", err
	}

	br := bytes.NewReader(body)

	x, err := exif.Decode(br)

	if err != nil {

		if exif.IsExifError(err) {
			return "", nil
		}

		if exif.IsCriticalError(err) {
			return "", nil
		}

		return "", err
	}

	tag, err := x.Get(exif.Orientation)

	if err != nil {
		return "", nil
	}

	// log.Println(path, tag)

	orientation, err := tag.Int64(0)

	if err != nil {
		return "", err
	}

	if orientation == 1 {
		return "", nil
	}

	br.Seek(0, 0)

	im, _, err := util.DecodeImageFromReader(br)

	if err != nil {
		return "", err
	}

	angle, _, _ := exifutil.ProcessOrientation(orientation)
	rotated := exifutil.Rotate(im, angle)

	tmpfile, _, err := tempfile.TempFileWithImage(ctx, bucket, rotated)
	return tmpfile, err
}
