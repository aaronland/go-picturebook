package functions

import (
	"context"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	_ "log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func PictureBookFilterFuncFromString(str_filter string) (PictureBookFilterFunc, error) {

	var filter PictureBookFilterFunc

	switch str_filter {

	case "cooperhewitt":
		filter = CooperHewittShoeboxFilterFunc
	case "default":
		filter = DefaultFilterFunc
	case "flickr":
		filter = FlickrArchiveFilterFunc
	case "orthis":
		filter = OrThisFilterFunc
	default:
		return nil, errors.New("Invalid filter type")
	}

	return filter, nil
}

func DefaultFilterFunc(ctx context.Context, path string) (bool, error) {
	return true, nil
}

func OrThisFilterFunc(ctx context.Context, path string) (bool, error) {

	fname := filepath.Base(path)

	if !orthis_re.MatchString(fname) {
		return false, nil
	}

	return true, nil
}

func MakeOrThisYearFilterFunc(ctx context.Context, year int) PictureBookFilterFunc {

	fn := func(ctx context.Context, path string) (bool, error) {

		fname := filepath.Base(path)

		if !orthis_re.MatchString(fname) {
			return false, nil
		}

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

		str_year := strconv.Itoa(year)

		if !strings.HasPrefix(date_rsp.String(), str_year) {
			return false, nil
		}

		return true, nil
	}

	return fn
}

func FlickrArchiveFilterFunc(ctx context.Context, path string) (bool, error) {

	if !flickr_re.MatchString(path) {
		return false, nil
	}

	return true, nil
}

func CooperHewittShoeboxFilterFunc(ctx context.Context, path string) (bool, error) {

	if !strings.HasSuffix(path, "_b.jpg") {
		return false, nil
	}

	root := filepath.Dir(path)
	info := filepath.Join(root, "index.json")

	_, err := os.Stat(info)

	if os.IsNotExist(err) {
		return true, nil
	}

	if err != nil {
		return true, err
	}

	info_fh, err := os.Open(info)

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

	_, err = os.Stat(object_info)

	if os.IsNotExist(err) {
		return true, nil
	}

	if err != nil {
		return true, err
	}

	object_fh, err := os.Open(object_info)

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
