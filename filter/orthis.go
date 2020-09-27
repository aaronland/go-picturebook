package filter

import (
	"context"
	"github.com/tidwall/gjson"
	"gocloud.dev/blob"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func init() {

	ctx := context.Background()
	err := RegisterFilter(ctx, "orthis", NewOrThisFilter)

	if err != nil {
		panic(err)
	}
}

type OrThisFilter struct {
	Filter
	year int
}

func NewOrThisFilter(ctx context.Context, uri string) (Filter, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	f := &OrThisFilter{
		year: 0,
	}

	q := u.Query()
	str_year := q.Get("year")

	if str_year != "" {

		t, err := time.Parse("2006", str_year)

		if err != nil {
			return nil, err
		}

		f.year = t.Year()
	}

	return f, nil
}

func (f *OrThisFilter) Continue(ctx context.Context, bucket *blob.Bucket, path string) (bool, error) {

	fname := filepath.Base(path)

	if !orthis_re.MatchString(fname) {
		return false, nil
	}

	if f.year != 0 {

		fname := filepath.Base(path)

		if !orthis_re.MatchString(fname) {
			return false, nil
		}

		root := filepath.Dir(path)
		root = filepath.Dir(root)

		index := filepath.Join(root, "index.json")

		fh, err := bucket.NewReader(ctx, index, nil)

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
	}

	return true, nil
}
