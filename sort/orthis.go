package sort

import (
	"context"
	"fmt"
	"github.com/aaronland/go-picturebook/picture"
	"github.com/tidwall/gjson"
	"io/ioutil"
	_ "log"
	"net/url"
	"os"
	"path/filepath"
	"sort"
)

func init() {

	ctx := context.Background()
	err := RegisterSorter(ctx, "orthis", NewOrThisSorter)

	if err != nil {
		panic(err)
	}
}

type OrThisSorter struct {
	Sorter
}

func NewOrThisSorter(ctx context.Context, uri string) (Sorter, error) {

	_, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	s := &OrThisSorter{}
	return s, nil
}

func (f *OrThisSorter) Sort(ctx context.Context, pictures []*picture.PictureBookPicture) ([]*picture.PictureBookPicture, error) {

	lookup := make(map[string]*picture.PictureBookPicture)
	candidates := make([]string, 0)

	for _, pic := range pictures {

		path := pic.Source

		root := filepath.Dir(path)
		root = filepath.Dir(root)

		index := filepath.Join(root, "index.json")

		fh, err := os.Open(index)

		if err != nil {
			return nil, err
		}

		defer fh.Close()

		body, err := ioutil.ReadAll(fh)

		if err != nil {
			return nil, err
		}

		id_rsp := gjson.GetBytes(body, "id")

		if !id_rsp.Exists() {
			continue
		}

		date_rsp := gjson.GetBytes(body, "date")

		if !date_rsp.Exists() {
			continue
		}

		key := fmt.Sprintf("%s-%s", date_rsp.String(), id_rsp.String())
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
