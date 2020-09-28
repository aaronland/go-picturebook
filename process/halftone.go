package process

import (
	"context"
	"github.com/aaronland/go-image-halftone"
	"github.com/aaronland/go-image-tools/util"
	"github.com/aaronland/go-picturebook/tempfile"
	"gocloud.dev/blob"
	"net/url"
)

func init() {

	ctx := context.Background()
	err := RegisterProcess(ctx, "halftone", NewHalftoneProcess)

	if err != nil {
		panic(err)
	}
}

type HalftoneProcess struct {
	Process
}

func NewHalftoneProcess(ctx context.Context, uri string) (Process, error) {

	_, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	f := &HalftoneProcess{}

	return f, nil
}

func (f *HalftoneProcess) Transform(ctx context.Context, bucket *blob.Bucket, path string) (string, error) {

	fh, err := bucket.NewReader(ctx, path, nil)

	if err != nil {
		return "", err
	}

	defer fh.Close()

	im, _, err := util.DecodeImageFromReader(fh)

	if err != nil {
		return "", err
	}

	opts := halftone.NewDefaultHalftoneOptions()
	dithered, err := halftone.HalftoneImage(ctx, im, opts)

	if err != nil {
		return "", err
	}

	tmpfile, _, err := tempfile.TempFileWithImage(ctx, bucket, dithered)
	return tmpfile, err
}
