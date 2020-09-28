package caption

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"gocloud.dev/blob"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"time"
)

func init() {

	ctx := context.Background()
	err := RegisterCaption(ctx, "cooperhewitt", NewCooperHewittCaption)

	if err != nil {
		panic(err)
	}
}

type CooperHewittCaption struct {
	Caption
}

func NewCooperHewittCaption(ctx context.Context, uri string) (Caption, error) {

	_, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	c := &CooperHewittCaption{}

	return c, nil
}

func (c *CooperHewittCaption) Text(ctx context.Context, bucket *blob.Bucket, path string) (string, error) {

	root := filepath.Dir(path)
	info := filepath.Join(root, "index.json")

	exists, err := bucket.Exists(ctx, info)

	if err != nil {
		return "", err
	}

	if !exists {
		return "", errors.New("Missing index.json")
	}

	fh, err := bucket.NewReader(ctx, info, nil)

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
