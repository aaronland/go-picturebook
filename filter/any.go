package filter

import (
	"context"
	"net/url"
	"gocloud.dev/blob"		
)

func init() {

	ctx := context.Background()
	err := RegisterFilter(ctx, "any", NewAnyFilter)

	if err != nil {
		panic(err)
	}
}

type AnyFilter struct {
	Filter
}

func NewAnyFilter(ctx context.Context, uri string) (Filter, error) {

	_, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	f := &AnyFilter{}

	return f, nil
}

func (f *AnyFilter) Continue(ctx context.Context, bucket *blob.Bucket, path string) (bool, error) {
	return true, nil
}
