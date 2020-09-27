package caption

import (
	"context"
	"gocloud.dev/blob"
	"net/url"
)

func init() {

	ctx := context.Background()
	err := RegisterCaption(ctx, "none", NewNoneCaption)

	if err != nil {
		panic(err)
	}
}

type NoneCaption struct {
	Caption
}

func NewNoneCaption(ctx context.Context, uri string) (Caption, error) {

	_, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	c := &NoneCaption{}
	return c, nil
}

func (c *NoneCaption) Text(ctx context.Context, bucket *blob.Bucket, path string) (string, error) {
	return "", nil
}
