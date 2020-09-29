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
		return nil, err
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
			log.Println(path, err)
			continue
		}

		defer r.Close()

		mtime := r.ModTime()
		sz := r.Size()

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
