package functions

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// when it comes to returning strings (HTML) see also:
// https://github.com/straup/go-image-tools/issues/7

func PictureBookCaptionFuncFromString(caption string) (PictureBookCaptionFunc, error) {

	var capt PictureBookCaptionFunc

	switch caption {

	case "cooperhewitt":
		capt = CooperHewittShoeboxCaptionFunc
	case "default":
		capt = FilenameCaptionFunc
	case "filename":
		capt = FilenameCaptionFunc
	case "flickr":
		capt = FlickrArchiveCaptionFunc
	case "orthis":
		capt = OrThisCaptionFunc
	case "parent":
		capt = FilenameAndParentCaptionFunc
	case "none":
		capt = NoneCaptionFunc
	default:
		return nil, errors.New("Invalid caption type")
	}

	return capt, nil
}

func DefaultCaptionFunc(path string) (string, error) {
	return FilenameCaptionFunc(path)
}

func FilenameCaptionFunc(path string) (string, error) {

	fname := filepath.Base(path)
	return fname, nil
}

func FilenameAndParentCaptionFunc(path string) (string, error) {

	root := filepath.Dir(path)
	parent := filepath.Base(root)
	fname := filepath.Base(path)

	return filepath.Join(parent, fname), nil
}

func OrThisCaptionFunc(path string) (string, error) {

	fname := filepath.Base(path)
	pat := "-or-this.jpg"

	if !strings.HasSuffix(fname, pat) {
		return "", nil
	}

	ymd := strings.Replace(fname, pat, "", -1)

	tm, err := time.Parse("2006-01-02", ymd)

	if err != nil {
		return "", nil
	}

	dt := tm.Format("January 02 2006")
	caption := fmt.Sprintf("%s / or this...", strings.ToLower(dt))

	return caption, nil
}

func NoneCaptionFunc(path string) (string, error) {
	return "", nil
}

func FlickrArchiveCaptionFunc(path string) (string, error) {

	ext := filepath.Ext(path)

	img_ext := fmt.Sprintf("_o%s", ext)
	info_ext := "_i.json"

	info := strings.Replace(path, img_ext, info_ext, -1)

	_, err := os.Stat(info)

	if err != nil {
		return "", err
	}

	fh, err := os.Open(info)

	if err != nil {
		return "", err
	}

	defer fh.Close()

	body, err := ioutil.ReadAll(fh)

	var item interface{}
	err = json.Unmarshal(body, &item)

	if err != nil {
		return "", err
	}

	var rsp gjson.Result
	var photo_id int64
	var title string
	var taken string

	rsp = gjson.GetBytes(body, "photo.id")

	if !rsp.Exists() {
		return "", errors.New("Missing photo ID")
	}

	photo_id = rsp.Int()

	rsp = gjson.GetBytes(body, "photo.title._content")

	if !rsp.Exists() {
		return "", errors.New("Missing title")
	}

	title = rsp.String()

	rsp = gjson.GetBytes(body, "photo.dates.taken")

	if !rsp.Exists() {
		return "", errors.New("Missing date")
	}

	taken = rsp.String()

	// go... Y U SO WEIRD...
	// https://golang.org/src/time/format.go

	tm, err := time.Parse("2006-01-02 15:04:05", taken)

	if err != nil {
		return "", nil
	}

	dt := tm.Format("Jan 02, 2006")

	caption := fmt.Sprintf("<b>%s</b><br />%s / %d", title, dt, photo_id)
	return caption, nil
}

func CooperHewittShoeboxCaptionFunc(path string) (string, error) {

	root := filepath.Dir(path)
	info := filepath.Join(root, "index.json")

	_, err := os.Stat(info)

	if err != nil {
		return "", err
	}

	fh, err := os.Open(info)

	if err != nil {
		return "", err
	}

	defer fh.Close()

	body, err := ioutil.ReadAll(fh)

	var item interface{}
	err = json.Unmarshal(body, &item)

	if err != nil {
		return "", err
	}

	var rsp gjson.Result
	var title string
	var acc string
	var object_id int64
	var created int64

	rsp = gjson.GetBytes(body, "refers_to_a")

	if !rsp.Exists() {
		return "", errors.New("Unknown shoebox item")
	}

	isa := rsp.String()

	if isa != "object" {
		return "", errors.New("Unsuported shoebox item")
	}

	rsp = gjson.GetBytes(body, "refers_to.title")

	if !rsp.Exists() {
		return "", errors.New("Object information missing title")
	}

	title = rsp.String()

	rsp = gjson.GetBytes(body, "created")

	if rsp.Exists() {
		created = rsp.Int()
	}

	rsp = gjson.GetBytes(body, "refers_to.accession_number")

	if rsp.Exists() {
		acc = rsp.String()
	}

	rsp = gjson.GetBytes(body, "refers_to.id")

	if rsp.Exists() {
		object_id = rsp.Int()
	}

	tm := time.Unix(created, 0)
	dt := tm.Format("Jan 02, 2006")

	caption := fmt.Sprintf("<b>%s</b><br />%s (%d)<br /><i>%s</i>", title, acc, object_id, dt)
	return caption, nil
}
