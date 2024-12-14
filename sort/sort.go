// package sort provides a common interface for sorting a list of images to be included in a picturebook.
package sort

import (
	"context"
	"net/url"
	"regexp"

	"github.com/aaronland/go-picturebook/bucket"
	"github.com/aaronland/go-picturebook/picture"
	"github.com/aaronland/go-roster"
)

// orthis_re is a regular expression pattern for matching files with names following the convention for (aaronland) Or This "original" photos.
var orthis_re *regexp.Regexp

func init() {
	orthis_re = regexp.MustCompile(`^(\d+)_[a-zA-Z0-9]+_o\.jpg$`)
}

// type Sorter provides a common interface for sorting a list of images to be included in a picturebook.
type Sorter interface {
	// Sort takes a list of `picture.PictureBookPicture` instances that are stored in a gocloud.dev/blob Bucket instance and returns new list of sorted picture.PictureBookPicture instances.
	Sort(context.Context, bucket.Bucket, []*picture.PictureBookPicture) ([]*picture.PictureBookPicture, error)
}

// type SorterInitializeFunc defined a common initialization function for instances implementing the Sorter interface.
// This is specified when the packages definining those instances call `RegisterSorter` and invoked with the `NewSorter`
// method is called.
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

// RegisterSorter associates a URI scheme with a `SorterInitializeFunc` initialization function.
func RegisterSorter(ctx context.Context, name string, fn SorterInitializeFunc) error {

	err := ensureRoster()

	if err != nil {
		return err
	}

	return sorters.Register(ctx, name, fn)
}

// NewSorter returns a new `Sorter` instance for 'uri' whose scheme is expected to have been associated
// with an `SorterInitializeFunc` (by the `RegisterSorter` method.
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

// AvailableSorters returns the list of schemes that have been registered with `SorterInitializeFunc` functions.
func AvailableSorters() []string {
	ctx := context.Background()
	return sorters.Drivers(ctx)
}
