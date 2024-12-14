package sort

import (
	"context"
	"fmt"
	"net/url"
	"sort"

	"github.com/aaronland/go-picturebook/bucket"
	"github.com/aaronland/go-picturebook/picture"
)

func init() {

	ctx := context.Background()
	err := RegisterSorter(ctx, "modtime", NewModTimeSorter)

	if err != nil {
		panic(err)
	}
}

// type ModTimeSorter implements the `Sorter` interface to sort a list of `picture.PictureBookPicture` by their modification dates.
type ModTimeSorter struct {
	Sorter
}

// NewModTimeSorter returns a new instance of `ModTimeSorter` for 'uri' which must be parsable as a valid `net/url` URL instance.
func NewModTimeSorter(ctx context.Context, uri string) (Sorter, error) {

	_, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI for NewModTimeSorter, %w", err)
	}

	s := &ModTimeSorter{}
	return s, nil
}

// Sort sorts a list of `picture.PictureBookPicture` by their modification dates.
func (f *ModTimeSorter) Sort(ctx context.Context, source_bucket bucket.Bucket, pictures []*picture.PictureBookPicture) ([]*picture.PictureBookPicture, error) {

	lookup := make(map[string]*picture.PictureBookPicture)
	candidates := make([]string, 0)

	for _, pic := range pictures {

		path := pic.Source

		attrs, err := source_bucket.Attributes(ctx, path)

		if err != nil {
			return nil, fmt.Errorf("Failed to dereive attributes from %s for modtime sorting, %v\n", path, err)
		}

		mtime := attrs.ModTime
		sz := attrs.Size

		ts := mtime.Unix()

		key := fmt.Sprintf("%d-%d", ts, sz)
		lookup[key] = pic
	}

	for key, _ := range lookup {
		candidates = append(candidates, key)
	}

	sort.Strings(candidates)

	sorted := make([]*picture.PictureBookPicture, len(candidates))

	for idx, key := range candidates {
		sorted[idx] = lookup[key]
	}

	return sorted, nil
}
