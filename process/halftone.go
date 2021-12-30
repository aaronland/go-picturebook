package process

import (
	"context"
	"fmt"
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

// type HalftoneProcess implements the `Process` interface and applies a "halftone" dithering transformation to an image.
type HalftoneProcess struct {
	Process
}

// NewHalftoneProcess returns a new instance of `HalftoneProcess` for 'uri' which must be parsable as a valid `net/url` URL instance.
func NewHalftoneProcess(ctx context.Context, uri string) (Process, error) {

	_, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URL for NewHalftoneProcess, %w", err)
	}

	f := &HalftoneProcess{}

	return f, nil
}

// Tranform applies a "halftone" dithering tranformation to 'path' in 'source_bucket' and writes the results to 'target_bucket' returning
// a new relative path on success.
func (f *HalftoneProcess) Transform(ctx context.Context, source_bucket *blob.Bucket, target_bucket *blob.Bucket, path string) (string, error) {

	fh, err := source_bucket.NewReader(ctx, path, nil)

	if err != nil {
		return "", fmt.Errorf("Failed to create new reader for %s, %w", path, err)
	}

	defer fh.Close()

	im, _, err := util.DecodeImageFromReader(fh)

	if err != nil {
		return "", fmt.Errorf("Failed to decode image for %s, %w", path, err)
	}

	opts := halftone.NewDefaultHalftoneOptions()
	dithered, err := halftone.HalftoneImage(ctx, im, opts)

	if err != nil {
		return "", fmt.Errorf("Failed to halftone image for %s, %w", path, err)
	}

	tmpfile, _, err := tempfile.TempFileWithImage(ctx, target_bucket, dithered)

	if err != nil {
		return "", fmt.Errorf("Failed to write temp file (halftone) for %s, %w", path, err)
	}

	return tmpfile, nil
}
