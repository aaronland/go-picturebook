// package filter provides a common interfaces for filtering mechanisms used to determine inclusion in a picturebook.
package filter

import (
	"context"
	"net/url"
	"regexp"

	"github.com/aaronland/go-picturebook/bucket"
	"github.com/aaronland/go-roster"
)

// flickr_re is a regular expression pattern for matching files with names following the convention for Flickr "original" photos.
var flickr_re *regexp.Regexp

// orthis_re is a regular expression pattern for matching files with names following the convention for (aaronland) Or This "original" photos.
var orthis_re *regexp.Regexp

func init() {
	flickr_re = regexp.MustCompile(`o_\.\.*$`)
	orthis_re = regexp.MustCompile(`^(\d+)_[a-zA-Z0-9]+_o\.jpg$`)
}

// type Filter provides a common interfaces for filtering mechanisms used to determine inclusion in a picturebook.
type Filter interface {
	// Continue determines whether a file contained in a gocloud.dev/blob Bucket instance should be included in a picturebook.
	Continue(context.Context, bucket.Bucket, string) (bool, error)
}

// type FilterInitializeFunc defined a common initialization function for instances implementing the Filter interface.
// This is specified when the packages definining those instances call `RegisterFilter` and invoked with the `NewFilter`
// method is called.
type FilterInitializeFunc func(context.Context, string) (Filter, error)

var filters roster.Roster

func ensureRoster() error {

	if filters == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		filters = r
	}

	return nil
}

// RegisterFilter associates a URI scheme with a `FilterInitializeFunc` initialization function.
func RegisterFilter(ctx context.Context, name string, fn FilterInitializeFunc) error {

	err := ensureRoster()

	if err != nil {
		return err
	}

	return filters.Register(ctx, name, fn)
}

// NewFilter returns a new `Filter` instance for 'uri' whose scheme is expected to have been associated
// with an `FilterInitializeFunc` (by the `RegisterFilter` method.
func NewFilter(ctx context.Context, uri string) (Filter, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := filters.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	fn := i.(FilterInitializeFunc)

	filter, err := fn(ctx, uri)

	if err != nil {
		return nil, err
	}

	return filter, nil
}

// AvailableFilter returns the list of schemes that have been registered with `FilterInitializeFunc` functions.
func AvailableFilters() []string {
	ctx := context.Background()
	return filters.Drivers(ctx)
}
