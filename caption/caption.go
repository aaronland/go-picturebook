// package caption provides a common interface for different mechanisms to derive captions for images.
package caption

import (
	"context"
	"fmt"
	"net/url"
	"regexp"

	"github.com/aaronland/go-roster"
	"gocloud.dev/blob"
)

// flickr_re is a regular expression pattern for matching files with names following the convention for Flickr "original" photos.
var flickr_re *regexp.Regexp

// orthis_re is a regular expression pattern for matching files with names following the convention for (aaronland) Or This "original" photos.
var orthis_re *regexp.Regexp

func init() {
	flickr_re = regexp.MustCompile(`o_\.\.*$`)
	orthis_re = regexp.MustCompile(`^(\d+)_[a-zA-Z0-9]+_o\.jpg$`)
}

// type Caption provides a common interface for different mechanisms to derive captions for images.
type Caption interface {
	// Text produces a caption derived from a file contained in a gocloud.dev/blob Bucket instance.
	Text(context.Context, *blob.Bucket, string) (string, error)
}

// type CaptionInitializeFunc defined a common initialization function for instances implementing the Caption interface.
// This is specified when the packages definining those instances call `RegisterCaption` and invoked with the `NewCaption`
// method is called.
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

// RegisterCaption associates a URI scheme with a `CaptionInitializeFunc` initialization function.
func RegisterCaption(ctx context.Context, name string, fn CaptionInitializeFunc) error {

	err := ensureRoster()

	if err != nil {
		return fmt.Errorf("Failed to ensure captions roster, %w", err)
	}

	return captions.Register(ctx, name, fn)
}

// NewCaption returns a new `Caption` instance for 'uri' whose scheme is expected to have been associated
// with an `CaptionInitializeFunc` (by the `RegisterCaption` method.
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

// AvailableCaption returns the list of schemes that have been registered with `CaptionInitializeFunc` functions.
func AvailableCaptions() []string {
	ctx := context.Background()
	return captions.Drivers(ctx)
}
