// package caption provides a common interface for different mechanisms to derive captions for images.
package caption

import (
	"context"
	"fmt"
	"github.com/aaronland/go-roster"
	"gocloud.dev/blob"
	"net/url"
	"regexp"
)

var flickr_re *regexp.Regexp
var orthis_re *regexp.Regexp

func init() {
	flickr_re = regexp.MustCompile(`o_\.\.*$`)
	orthis_re = regexp.MustCompile(`^(\d+)_[a-zA-Z0-9]+_o\.jpg$`)
}

type Caption interface {
	Text(context.Context, *blob.Bucket, string) (string, error)
	// Text(context.Context, io.ReadSeeker, string) (string, error)
}

type CaptionInitializeFunc func(context.Context, string) (Caption, error)

var captions roster.Roster

func ensureRoster() error {

	if captions == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return fmt.Errorf("Failed to create new roster for captions, %w", err)
		}

		captions = r
	}

	return nil
}

func RegisterCaption(ctx context.Context, name string, fn CaptionInitializeFunc) error {

	err := ensureRoster()

	if err != nil {
		return fmt.Errorf("Failed to ensure captions roster, %w", err)
	}

	return captions.Register(ctx, name, fn)
}

func NewCaption(ctx context.Context, uri string) (Caption, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI for NewCaption, %w", err)
	}

	scheme := u.Scheme

	i, err := captions.Driver(ctx, scheme)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive driver for '%s' caption scheme, %w", scheme, err)
	}

	fn := i.(CaptionInitializeFunc)

	caption, err := fn(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("CaptionInitializeFunc failed, %w", err)
	}

	return caption, nil
}

func AvailableCaptions() []string {
	ctx := context.Background()
	return captions.Drivers(ctx)
}
