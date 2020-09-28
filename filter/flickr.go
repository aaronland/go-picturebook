package filter

import (
	"context"
	"gocloud.dev/blob"
	"net/url"
)

func init() {

	ctx := context.Background()
	err := RegisterFilter(ctx, "flickr", NewFlickrFilter)

	if err != nil {
		panic(err)
	}
}

type FlickrFilter struct {
	Filter
}

func NewFlickrFilter(ctx context.Context, uri string) (Filter, error) {

	_, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	f := &FlickrFilter{}

	return f, nil
}

func (f *FlickrFilter) Continue(ctx context.Context, bucket *blob.Bucket, path string) (bool, error) {

	if !flickr_re.MatchString(path) {
		return false, nil
	}

	return true, nil
}
