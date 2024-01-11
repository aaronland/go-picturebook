package caption

import (
	"context"
	"fmt"
	"net/url"

	"gocloud.dev/blob"
)

func init() {

	ctx := context.Background()
	err := RegisterCaption(ctx, "modtime", NewModtimeCaption)

	if err != nil {
		panic(err)
	}
}

// type ModtimeCaption implements the `Caption` interface and derives caption text from image modification times.
type ModtimeCaption struct {
	Caption
	format string
}

// NewExifCaption return a new instance of `ModtimeCaption` for 'url' which is expected to take
// the form of:
//
//	modtime://
func NewModtimeCaption(ctx context.Context, uri string) (Caption, error) {

	_, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URL, %w", err)
	}

	c := &ModtimeCaption{
		format: "January 02, 2006",
	}

	return c, nil
}

// Text returns a caption string derived from the modification time of 'path'
func (c *ModtimeCaption) Text(ctx context.Context, bucket *blob.Bucket, path string) (string, error) {

	attrs, err := bucket.Attributes(ctx, path)

	if err != nil {
		return "", fmt.Errorf("Failed to derive attributes for %s, %w", path, err)
	}

	t := attrs.ModTime

	return t.Format(c.format), nil
}
