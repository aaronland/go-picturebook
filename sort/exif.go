package sort

import (
	"context"
	"fmt"
	"github.com/aaronland/go-picturebook/picture"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"gocloud.dev/blob"
	"log"
	"net/url"
	"sort"
)

func init() {

	ctx := context.Background()
	err := RegisterSorter(ctx, "exif", NewExifSorter)

	if err != nil {
		panic(err)
	}

	exif.RegisterParsers(mknote.All...)
}

type ExifSorter struct {
	Sorter
}

func NewExifSorter(ctx context.Context, uri string) (Sorter, error) {

	_, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	s := &ExifSorter{}
	return s, nil
}

func (f *ExifSorter) Sort(ctx context.Context, bucket *blob.Bucket, pictures []*picture.PictureBookPicture) ([]*picture.PictureBookPicture, error) {

	lookup := make(map[string]*picture.PictureBookPicture)
	candidates := make([]string, 0)

	for _, pic := range pictures {

		path := pic.Source

		fh, err := bucket.NewReader(ctx, path, nil)

		if err != nil {
			log.Println(path, err)
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
