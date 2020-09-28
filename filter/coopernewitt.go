package filter

import (
	"context"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"gocloud.dev/blob"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"strings"
)

func init() {

	ctx := context.Background()
	err := RegisterFilter(ctx, "cooperhewitt", NewCooperHewittFilter)

	if err != nil {
		panic(err)
	}
}

type CooperHewittFilter struct {
	Filter
}

func NewCooperHewittFilter(ctx context.Context, uri string) (Filter, error) {

	_, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	f := &CooperHewittFilter{}

	return f, nil
}

func (f *CooperHewittFilter) Continue(ctx context.Context, bucket *blob.Bucket, path string) (bool, error) {

	if !strings.HasSuffix(path, "_b.jpg") {
		return false, nil
	}

	root := filepath.Dir(path)
	info := filepath.Join(root, "index.json")

	exists, err := bucket.Exists(ctx, info)

	if err != nil {
		return true, err
	}

	if !exists {
		return true, nil
	}

	info_fh, err := bucket.NewReader(ctx, info, nil)

	if err != nil {
		return true, err
	}

	defer info_fh.Close()

	info_body, err := ioutil.ReadAll(info_fh)

	if err != nil {
		return true, err
	}

	var rsp gjson.Result

	rsp = gjson.GetBytes(info_body, "refers_to_uid")

	if !rsp.Exists() {
		return true, errors.New("Unable to determine refers_to_uid")
	}

	uid := rsp.Int()

	object_fname := fmt.Sprintf("%d.json", uid)
	object_info := filepath.Join(root, object_fname)

	exists, err = bucket.Exists(ctx, object_info)

	if err != nil {
		return true, err
	}

	if !exists {
		return true, nil
	}

	object_fh, err := bucket.NewReader(ctx, object_info, nil)

	if err != nil {
		return true, err
	}

	defer object_fh.Close()

	object_body, err := ioutil.ReadAll(object_fh)

	if err != nil {
		return true, err
	}

	rsp = gjson.GetBytes(object_body, "object.images")

	if !rsp.Exists() {
		return true, errors.New("Unable to determine object.images")
	}

	fname := filepath.Base(path)
	ok := false

	for _, im := range rsp.Array() {

		for k, details := range im.Map() {

			if k != "b" {
				continue
			}

			rsp = details.Get("is_primary")

			if !rsp.Exists() {
				continue
			}

			if rsp.Int() != 1 {
				continue
			}

			rsp = details.Get("url")

			if !rsp.Exists() {
				continue
			}

			url := rsp.String()
			url_fname := filepath.Base(url)

			if fname != url_fname {
				continue
			}

			ok = true
			break
		}

		if ok == true {
			break
		}
	}

	return ok, nil
}
