package sort

import (
	"context"
	"errors"
	"github.com/aaronland/go-picturebook/picture"
	_ "github.com/tidwall/gjson"
	_ "io/ioutil"
	"net/url"
	_ "os"
	_ "path/filepath"
	_ "strconv"
	_ "strings"
	_ "time"
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

func (f *OrThisSorter) Sort(ctx context.Context, pictures []*picture.PictureBookPicture) error {

	return errors.New("Not implemented")

	/*
			root := filepath.Dir(path)
			root = filepath.Dir(root)

			index := filepath.Join(root, "index.json")

			fh, err := os.Open(index)

			if err != nil {
				return false, err
			}

			defer fh.Close()

			body, err := ioutil.ReadAll(fh)

			if err != nil {
				return false, err
			}

			date_rsp := gjson.GetBytes(body, "date")

			if !date_rsp.Exists() {
				return false, nil
			}

			str_year := strconv.Itoa(f.year)

			if !strings.HasPrefix(date_rsp.String(), str_year) {
				return false, nil
			}

		return true, nil
	*/
}
