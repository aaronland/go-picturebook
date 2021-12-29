package sort

import (
	"context"
	"fmt"
	"github.com/aaronland/go-picturebook/picture"
	"gocloud.dev/blob"
	"log"
	"net/url"
	"sort"
)

func init() {

	ctx := context.Background()
	err := RegisterSorter(ctx, "modtime", NewModTimeSorter)

	if err != nil {
		panic(err)
	}
}

type ModTimeSorter struct {
	Sorter
}

func NewModTimeSorter(ctx context.Context, uri string) (Sorter, error) {

	_, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI for NewModTimeSorter, %w", err)
	}

	s := &ModTimeSorter{}
	return s, nil
}

func (f *ModTimeSorter) Sort(ctx context.Context, bucket *blob.Bucket, pictures []*picture.PictureBookPicture) ([]*picture.PictureBookPicture, error) {

	lookup := make(map[string]*picture.PictureBookPicture)
	candidates := make([]string, 0)

	for _, pic := range pictures {

		path := pic.Source

		r, err := bucket.NewReader(ctx, path, nil)

		if err != nil {
			log.Printf("Failed to open %s for modtime sorting, %v\n", path, err)
			continue
		}

		mtime := r.ModTime()
		sz := r.Size()
		r.Close()

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
