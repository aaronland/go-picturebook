package caption

import (
	"context"
	"net/url"
	"path/filepath"
)

func init() {

	ctx := context.Background()
	err := RegisterCaption(ctx, "filename", NewFilenameCaption)

	if err != nil {
		panic(err)
	}
}

type FilenameCaption struct {
	Caption
	parent bool
}

func NewFilenameCaption(ctx context.Context, uri string) (Caption, error) {

	_, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	c := &FilenameCaption{
		parent: false,
	}

	return c, nil
}

func (c *FilenameCaption) Text(ctx context.Context, path string) (string, error) {

	if c.parent {

		root := filepath.Dir(path)
		parent := filepath.Base(root)
		fname := filepath.Base(path)

		return filepath.Join(parent, fname), nil
	}

	fname := filepath.Base(path)
	return fname, nil
}
