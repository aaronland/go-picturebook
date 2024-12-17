package process

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aaronland/go-image-contour"
	"github.com/aaronland/go-image/decode"
	"github.com/aaronland/go-picturebook/bucket"
	"github.com/aaronland/go-picturebook/tempfile"
)

func init() {

	ctx := context.Background()
	err := RegisterProcess(ctx, "contour", NewContourProcess)

	if err != nil {
		panic(err)
	}
}

// type ContourProcess implements the `Process` interface and transforms an image in to a series of black and white "contour" lines.
type ContourProcess struct {
	Process
	iterations int
	scale      float64
}

// NewContourProcess returns a new instance of `ContourProcess` for 'uri' which must be parsable as a valid `net/url` URL instance.
//
//	contour://?{PARAMETERS}
//
// Where valid parameters are:
// * `iterations` The number of iterations to perform during the contour process. Default is 12.
// * `scale` The scale of the final contoured image. Default is 1.0.
func NewContourProcess(ctx context.Context, uri string) (Process, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI for NewContourProcess, %w", err)
	}

	q := u.Query()

	iterations := 12
	scale := 1.0

	str_iterations := q.Get("iterations")
	str_scale := q.Get("scale")

	if str_iterations != "" {
		v, err := strconv.Atoi(str_iterations)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse ?iterations= parameter, %w", err)
		}

		iterations = v
	}

	if str_scale != "" {
		v, err := strconv.ParseFloat(str_scale, 64)

		if err != nil {
			return nil, fmt.Errorf("Failed to parse ?scale= parameter, %w", err)
		}

		scale = v
	}

	f := &ContourProcess{
		iterations: iterations,
		scale:      scale,
	}

	return f, nil
}

// Tranform contours the image 'path' in 'bucket_bucket' and writes the results to 'target_bucket' returning
// a new relative path on success. If an image is not a JPEG file the method return an empty string.
func (f *ContourProcess) Transform(ctx context.Context, source_bucket bucket.Bucket, target_bucket bucket.Bucket, path string) (string, error) {

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

	dec, err := decode.NewDecoder(ctx, path)

	if err != nil {
		return "", fmt.Errorf("Failed to create new decoder, %w", err)
	}

	im, _, err := dec.Decode(ctx, r)

	if err != nil {
		return "", fmt.Errorf("Failed to decode image for %s, %w", path, err)
	}

	contoured_im, err := contour.ContourImage(ctx, im, f.iterations, f.scale)

	if err != nil {
		return "", fmt.Errorf("Failed to contour image for %s, %w", path, err)
	}

	tmpfile, _, err := tempfile.TempFileWithImage(ctx, target_bucket, contoured_im)

	if err != nil {
		return "", fmt.Errorf("Failed to write temp file (contour) for %s, %w", path, err)
	}

	return tmpfile, nil
}
