package process

// update to use go-image-rotate

import (
	"bytes"
	"context"
	"fmt"
	"io"
	_ "log"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/aaronland/go-image-tools/util"
	"github.com/aaronland/go-picturebook/tempfile"
	"github.com/microcosm-cc/exifutil"
	"github.com/rwcarlsen/goexif/exif"
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

	fh, err := source_bucket.NewReader(ctx, path, nil)

	if err != nil {
		return "", fmt.Errorf("Failed to create new reader for %s, %w", path, err)
	}

	defer fh.Close()

	body, err := io.ReadAll(fh)

	if err != nil {
		return "", fmt.Errorf("Failed to read %s, %w", path, err)
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

	orientation, err := tag.Int64(0)

	if err != nil {
		return "", fmt.Errorf("Failed to derive orientation from tag for %s, %w", path, err)
	}

	if orientation == 1 {
		return "", nil
	}

	br.Seek(0, 0)

	im, _, err := util.DecodeImageFromReader(br)

	if err != nil {
		return "", fmt.Errorf("Failed to decode image for %s, %w", path, err)
	}

	angle, _, _ := exifutil.ProcessOrientation(orientation)
	rotated := exifutil.Rotate(im, angle)

	tmpfile, _, err := tempfile.TempFileWithImage(ctx, target_bucket, rotated)

	if err != nil {
		return "", fmt.Errorf("Failed to write temp file (rotate) for %s, %w", path, err)
	}

	return tmpfile, nil
}
