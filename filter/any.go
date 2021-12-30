package filter

import (
	"context"
	"fmt"
	"gocloud.dev/blob"
	"net/url"
)

func init() {

	ctx := context.Background()
	err := RegisterFilter(ctx, "any", NewAnyFilter)

	if err != nil {
		panic(err)
	}
}

// type AnyFilter implements the `Filter` interface and allows any image to be included in a picturebook.
type AnyFilter struct {
	Filter
}

// NewAnyFilter returns a new instance of `AnyFilter` for 'uri'
func NewAnyFilter(ctx context.Context, uri string) (Filter, error) {

	_, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URL, %w", err)
	}

	f := &AnyFilter{}

	return f, nil
}

// Continues returns a boolean value signaling whether or not 'path' should be included in a picturebook.
func (f *AnyFilter) Continue(ctx context.Context, bucket *blob.Bucket, path string) (bool, error) {
	return true, nil
}
