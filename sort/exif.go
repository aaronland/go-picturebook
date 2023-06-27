package sort

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"sort"

	"github.com/aaronland/go-picturebook/picture"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"gocloud.dev/blob"
)

func init() {

	ctx := context.Background()
	err := RegisterSorter(ctx, "exif", NewExifSorter)

	if err != nil {
		panic(err)
	}

	exif.RegisterParsers(mknote.All...)
}

// type ExifSorter implements the `Sorter` interface to sort a list of `picture.PictureBookPicture` by their EXIF DateTime properties.
type ExifSorter struct {
	Sorter
}

// NewExifSorter returns a new instance of `ExifSorter` for 'uri' which must be parsable as a valid `net/url` URL instance.
func NewExifSorter(ctx context.Context, uri string) (Sorter, error) {

	_, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI for NewExifSorter, %w", err)
	}

	s := &ExifSorter{}
	return s, nil
}

// Sort sorts a list of `picture.PictureBookPicture` by their EXIF DateTime properties. If an image does not have an EXIF DateTime property it is
// excluded from the sorted result set.
func (f *ExifSorter) Sort(ctx context.Context, bucket *blob.Bucket, pictures []*picture.PictureBookPicture) ([]*picture.PictureBookPicture, error) {

	lookup := make(map[string]*picture.PictureBookPicture)
	candidates := make([]string, 0)

	for _, pic := range pictures {

		path := pic.Source

		fh, err := bucket.NewReader(ctx, path, nil)

		if err != nil {
			log.Printf("Failed to open %s for exif sorting, %v\n", path, err)
			continue
		}

		defer fh.Close()

		mtime := fh.ModTime()
		sz := fh.Size()

		ts := mtime.Unix()

		x, err := exif.Decode(fh)

		if err == nil {

			dt, err := x.DateTime()

			if err == nil {
				ts = dt.Unix()
			}
		} else {
			log.Printf("Failed to decode EXIF data for %s, %v\n", path, err)
		}

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
