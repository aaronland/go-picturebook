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

// type AnyFilter implements the `Filter` interface and allows any image to be processed
type AnyFilter struct {
	Filter
}

func NewAnyFilter(ctx context.Context, uri string) (Filter, error) {

	_, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URL, %w", err)
	}

	f := &AnyFilter{}

	return f, nil
}

func (f *AnyFilter) Continue(ctx context.Context, bucket *blob.Bucket, path string) (bool, error) {
	return true, nil
}
