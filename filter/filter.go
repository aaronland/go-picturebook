package filter

import (
	"context"
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

type Filter interface {
	Continue(context.Context, *blob.Bucket, string) (bool, error)
}

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

func RegisterFilter(ctx context.Context, name string, fn FilterInitializeFunc) error {

	err := ensureRoster()

	if err != nil {
		return err
	}

	return filters.Register(ctx, name, fn)
}

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

func AvailableFilters() []string {
	ctx := context.Background()
	return filters.Drivers(ctx)
}
