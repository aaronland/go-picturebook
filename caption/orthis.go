package caption

import (
	"context"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
)

func init() {

	ctx := context.Background()
	err := RegisterCaption(ctx, "orthis", NewOrThisCaption)

	if err != nil {
		panic(err)
	}
}

type OrThisCaption struct {
	Caption
}

func NewOrThisCaption(ctx context.Context, uri string) (Caption, error) {

	_, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	c := &OrThisCaption{}

	return c, nil
}

func (c *OrThisCaption) Text(ctx context.Context, path string) (string, error) {

	fname := filepath.Base(path)

	if !orthis_re.MatchString(fname) {
		return "", nil
	}

	root := filepath.Dir(path)
	root = filepath.Dir(root)

	index := filepath.Join(root, "index.json")

	fh, err := os.Open(index)

	if err != nil {
		return "", err
	}

	defer fh.Close()

	body, err := ioutil.ReadAll(fh)

	if err != nil {
		return "", err
	}

	m := orthis_re.FindStringSubmatch(fname)
	caption := fmt.Sprintf("untitled #%s", m[1])

	caption_rsp := gjson.GetBytes(body, "caption")

	if caption_rsp.Exists() {
		caption = fmt.Sprintf("%s / %s", caption, caption_rsp.String())
	}

	return caption, nil
}