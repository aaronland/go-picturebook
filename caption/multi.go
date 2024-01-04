package caption

import (
	"context"
	"fmt"
	"net/url"

	"gocloud.dev/blob"
)

func init() {

	ctx := context.Background()
	err := RegisterCaption(ctx, "multi", NewMultiCaption)

	if err != nil {
		panic(err)
	}
}

// type MultiCaption implements the `Caption` interface and derives caption text from image multis.
type MultiCaption struct {
	Caption
	providers []Caption
}

// NewExifCaption return a new instance of `MultiCaption` for 'url'
func NewMultiCaption(ctx context.Context, uri string) (Caption, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URL, %w", err)
	}

	q := u.Query()

	provider_uris := q["uri"]

	if len(provider_uris) == 0 {
		return nil, fmt.Errorf("No provider URIs")
	}

	providers := make([]Caption, len(provider_uris))

	for idx, _uri := range provider_uris {

		pr, err := NewCaption(ctx, _uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to create caption provider for %s, %w", _uri, err)
		}

		providers[idx] = pr
	}

	c := &MultiCaption{
		providers: providers,
	}

	return c, nil
}

// Text returns a caption string derived from the base name of 'path'
func (c *MultiCaption) Text(ctx context.Context, bucket *blob.Bucket, path string) (string, error) {

	for _, pr := range c.providers {

		txt, err := pr.Text(ctx, bucket, path)

		if err != nil {
			return "", fmt.Errorf("Failed to derive text for provider, %w", err)
		}

		if txt != "" {
			return txt, nil
		}
	}

	return "", nil
}
