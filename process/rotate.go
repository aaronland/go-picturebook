package process

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/aaronland/go-image/decode"
	"github.com/aaronland/go-image/rotate"
	"github.com/aaronland/go-picturebook/tempfile"
	"gocloud.dev/blob"
)

func init() {

	ctx := context.Background()
	err := RegisterProcess(ctx, "rotate", NewRotateProcess)

	if err != nil {
		panic(err)
	}
}

// type RotateProcess implements the `Process` interface and rotates and image based on its EXIF `Orientation` property.
type RotateProcess struct {
	Process
}

// NewRotateProcess returns a new instance of `RotateProcess` for 'uri' which must be parsable as a valid `net/url` URL instance.
func NewRotateProcess(ctx context.Context, uri string) (Process, error) {

	_, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI for NewRotateProcess, %w", err)
	}

	f := &RotateProcess{}

	return f, nil
}

// Tranform rotates the image 'path' in 'source_bucket' and writes the results to 'target_bucket' returning
// a new relative path on success. If an image is not a JPEG file the method return an empty string.
func (f *RotateProcess) Transform(ctx context.Context, source_bucket *blob.Bucket, target_bucket *blob.Bucket, path string) (string, error) {

	ext := filepath.Ext(path)
	ext = strings.ToLower(ext)

	if ext != ".jpg" && ext != ".jpeg" {
		return "", nil
	}

	r, err := source_bucket.NewReader(ctx, path, nil)

	if err != nil {
		return "", fmt.Errorf("Failed to create new reader for %s, %w", path, err)
	}

	defer r.Close()

	o, err := rotate.GetImageOrientation(ctx, r)

	if err != nil {
		return "", fmt.Errorf("Failed to derive orientation for %s, %w", path, err)
	}

	_, err = r.Seek(0, 0)

	if err != nil {
		return "", fmt.Errorf("Failed to rewind %s, %w", path, err)
	}

	dec, err := decode.NewDecoder(ctx, path)

	if err != nil {
		return "", fmt.Errorf("Failed to create new decoder, %w", err)
	}

	im, _, err := dec.Decode(ctx, r)

	if err != nil {
		return "", fmt.Errorf("Failed to decode image for %s, %w", path, err)
	}

	rotated, err := rotate.RotateImageWithOrientation(ctx, im, o)

	if err != nil {
		return "", fmt.Errorf("Failed to rotate %s, %w", path, err)
	}

	tmpfile, _, err := tempfile.TempFileWithImage(ctx, target_bucket, rotated)

	if err != nil {
		return "", fmt.Errorf("Failed to write temp file (rotate) for %s, %w", path, err)
	}

	return tmpfile, nil
}
