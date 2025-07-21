package process

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/aaronland/go-image/v2/decode"
	"github.com/aaronland/go-picturebook/bucket"
	"github.com/aaronland/go-picturebook/tempfile"
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
func (f *RotateProcess) Transform(ctx context.Context, source_bucket bucket.Bucket, target_bucket bucket.Bucket, path string) (string, error) {

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

	// decode.DecodeImage assumes Rotate: true by default
	// but being explicit to make the code a little clear
	
	decode_opts := &decode.DecodeImageOptions{
		Rotate: true,
	}
	
	rotated, _, _, err := decode.DecodeImageWithOptions(ctx, r, decode_opts)

	if err != nil {
		return "", fmt.Errorf("Failed to decode image for %s, %w", path, err)
	}

	tmpfile, _, err := tempfile.TempFileWithImage(ctx, target_bucket, rotated)

	if err != nil {
		return "", fmt.Errorf("Failed to write temp file (rotate) for %s, %w", path, err)
	}

	return tmpfile, nil
}
