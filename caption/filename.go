package caption

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"

	"gocloud.dev/blob"
)

func init() {

	ctx := context.Background()
	err := RegisterCaption(ctx, "filename", NewFilenameCaption)

	if err != nil {
		panic(err)
	}
}

// type FilenameCaption implements the `Caption` interface and derives caption text from image filenames.
type FilenameCaption struct {
	Caption
	parent bool
}

// NewExifCaption return a new instance of `FilenameCaption` for 'url'
func NewFilenameCaption(ctx context.Context, uri string) (Caption, error) {

	_, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URL, %w", err)
	}

	c := &FilenameCaption{
		parent: false,
	}

	return c, nil
}

// Text returns a caption string derived from the base name of 'path'
func (c *FilenameCaption) Text(ctx context.Context, bucket *blob.Bucket, path string) (string, error) {

	if c.parent {

		root := filepath.Dir(path)
		parent := filepath.Base(root)
		fname := filepath.Base(path)

		return filepath.Join(parent, fname), nil
	}

	fname := filepath.Base(path)
	return fname, nil
}
