package process

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aaronland/go-image/colour"
	"github.com/aaronland/go-image/decode"
	"github.com/aaronland/go-picturebook/bucket"
	"github.com/aaronland/go-picturebook/tempfile"
)

func init() {

	ctx := context.Background()
	RegisterProcess(ctx, "colorspace", NewColourSpaceProcess)
	RegisterProcess(ctx, "colourspace", NewColourSpaceProcess)

}

// type ColourSpaceProcess implements the `Process` interface to ensure that all the pixels in an image are mapped to a specific colour space.
type ColourSpaceProcess struct {
	Process
	profile string
}

// NewColourSpaceProcess returns a new instance of `ColourSpaceProcess` for 'uri' which is expected to take the form of:
//
//	colourspace://{PROFILE}
//	colorspace://{PROFILE}
//
// Where {PROFILE} is one of the following:
// * `displayp3` which maps pixels to Apple's Display P3 colour space
// * `adobergb` which maps pixels to Adobe's RGB colour space
func NewColourSpaceProcess(ctx context.Context, uri string) (Process, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URL for NewColourSpaceProcess, %w", err)
	}

	profile := u.Host

	switch profile {
	case "displayp3", "adobergb":
		// pass
	default:
		return nil, fmt.Errorf("Unsupported profile")
	}

	f := &ColourSpaceProcess{
		profile: profile,
	}

	return f, nil
}

// Tranform maps all the pixels in the image located at 'path' in 'bucket_bucket' and writes the results to 'target_bucket' returning
// a new relative path on success.
func (f *ColourSpaceProcess) Transform(ctx context.Context, source_bucket bucket.Bucket, target_bucket bucket.Bucket, path string) (string, error) {

	r, err := source_bucket.NewReader(ctx, path, nil)

	if err != nil {
		return "", fmt.Errorf("Failed to create new reader for %s, %w", path, err)
	}

	defer r.Close()

	dec, err := decode.NewDecoder(ctx, path)

	if err != nil {
		return "", fmt.Errorf("Failed to create new decoder for %s, %w", path, err)
	}

	im, _, err := dec.Decode(ctx, r)

	if err != nil {
		return "", fmt.Errorf("Failed to decode image for %s, %w", path, err)
	}

	switch f.profile {
	case "adobergb":
		im = colour.ToAdobeRGB(im)
	case "displayp3":
		im = colour.ToDisplayP3(im)
	default:
		return "", fmt.Errorf("Failed to adjust %s, unsupported profile", path)
	}

	tmpfile, _, err := tempfile.TempFileWithImage(ctx, target_bucket, im)

	if err != nil {
		return "", fmt.Errorf("Failed to write temp file for %s, %w", path, err)
	}

	return tmpfile, nil
}
