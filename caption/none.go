package caption

import (
	"context"
	"fmt"
	"net/url"

	"gocloud.dev/blob"
)

func init() {

	ctx := context.Background()
	err := RegisterCaption(ctx, "none", NewNoneCaption)

	if err != nil {
		panic(err)
	}
}

// type ExifCaption implements the `Caption` interface and returns empty caption strings.
type NoneCaption struct {
	Caption
}

// NewNoneCaption return a new instance of `NoneCaption` for 'url'
func NewNoneCaption(ctx context.Context, uri string) (Caption, error) {

	_, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URL, %w", err)
	}

	c := &NoneCaption{}
	return c, nil
}

// Text returns an empty caption string
func (c *NoneCaption) Text(ctx context.Context, bucket *blob.Bucket, path string) (string, error) {
	return "", nil
}
