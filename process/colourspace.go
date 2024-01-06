package process

import (
	"context"
	"fmt"
	"net/url"
	// "io"

	// "github.com/mandykoh/prism/meta/autometa"
	"github.com/aaronland/go-image/colour"
	"github.com/aaronland/go-image/decode"
	"github.com/aaronland/go-picturebook/tempfile"
	"gocloud.dev/blob"
)

func init() {

	ctx := context.Background()
	RegisterProcess(ctx, "colorspace", NewColourSpaceProcess)
	RegisterProcess(ctx, "colourspace", NewColourSpaceProcess)

}

// type ColourSpaceProcess implements the `Process` interface and applies a "ColourSpace" dithering transformation to an image.
type ColourSpaceProcess struct {
	Process
	profile string
}

// NewColourSpaceProcess returns a new instance of `ColourSpaceProcess` for 'uri' which must be parsable as a valid `net/url` URL instance.
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

// Tranform applies a "ColourSpace" dithering tranformation to 'path' in 'source_bucket' and writes the results to 'target_bucket' returning
// a new relative path on success.
func (f *ColourSpaceProcess) Transform(ctx context.Context, source_bucket *blob.Bucket, target_bucket *blob.Bucket, path string) (string, error) {

	r, err := source_bucket.NewReader(ctx, path, nil)

	if err != nil {
		return "", fmt.Errorf("Failed to create new reader for %s, %w", path, err)
	}

	defer r.Close()

	// It won't matter because by now we will be working with a temp file that is missing profile info...

	/*
		profile, err := f.deriveColourSpace(ctx, r)

		if err == nil {
			log.Println("HELLO", profile)
		} else {
			log.Println("WUT", path, err)
		}

		_, err = r.Seek(0, 0)

		if err != nil {
			return "", fmt.Errorf("Failed to rewind reader for %s, %w", path, err)
		}
	*/

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

/*
func (f *ColourSpaceProcess) deriveColourSpace(ctx context.Context, r io.Reader) (string, error) {

	md, _, err := autometa.Load(r)

	if err != nil {
		return "", err
	}

	profile, err := md.ICCProfile()

	if err != nil {
		return "", err
	}

	if profile == nil {
		return "", fmt.Errorf("Missing profile")
	}

	return profile.Description()

}
*/
