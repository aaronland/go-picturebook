package sort

import (
	"context"
	"github.com/aaronland/go-picturebook/picture"
	"github.com/aaronland/go-roster"
	"gocloud.dev/blob"
	"net/url"
	"regexp"
)

var orthis_re *regexp.Regexp

func init() {
	orthis_re = regexp.MustCompile(`^(\d+)_[a-zA-Z0-9]+_o\.jpg$`)
}

type Sorter interface {
	Sort(context.Context, *blob.Bucket, []*picture.PictureBookPicture) ([]*picture.PictureBookPicture, error)
}

type SorterInitializeFunc func(context.Context, string) (Sorter, error)

var sorters roster.Roster

func ensureRoster() error {

	if sorters == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		sorters = r
	}

	return nil
}

func RegisterSorter(ctx context.Context, name string, fn SorterInitializeFunc) error {

	err := ensureRoster()

	if err != nil {
		return err
	}

	return sorters.Register(ctx, name, fn)
}

func NewSorter(ctx context.Context, uri string) (Sorter, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := sorters.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	fn := i.(SorterInitializeFunc)

	sorter, err := fn(ctx, uri)

	if err != nil {
		return nil, err
	}

	return sorter, nil
}

func AvailableSorters() []string {
	ctx := context.Background()
	return sorters.Drivers(ctx)
}
