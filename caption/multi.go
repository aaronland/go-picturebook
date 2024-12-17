package caption

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/aaronland/go-picturebook/bucket"
)

func init() {

	ctx := context.Background()
	err := RegisterCaption(ctx, "multi", NewMultiCaption)

	if err != nil {
		panic(err)
	}
}

type MultiCaptionOptions struct {
	Captions   []Caption
	Combined   bool
	AllowEmpty bool
}

// type MultiCaption implements the `Caption` interface and derives caption text from image multis.
type MultiCaption struct {
	Caption
	captions    []Caption
	combined    bool
	allow_empty bool
}

// NewExifCaption return a new instance of `MultiCaption` for 'url'
func NewMultiCaption(ctx context.Context, uri string) (Caption, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URL, %w", err)
	}

	q := u.Query()

	caption_uris := q["uri"]

	if len(caption_uris) == 0 {
		return nil, fmt.Errorf("No caption URIs")
	}

	captions := make([]Caption, len(caption_uris))

	for idx, _uri := range caption_uris {

		pr, err := NewCaption(ctx, _uri)

		if err != nil {
			return nil, fmt.Errorf("Failed to create caption caption for %s, %w", _uri, err)
		}

		captions[idx] = pr
	}

	opts := &MultiCaptionOptions{
		Captions:   captions,
		Combined:   false,
		AllowEmpty: true,
	}

	return NewMultiCaptionWithOptions(ctx, opts)
}

func NewMultiCaptionWithOptions(ctx context.Context, opts *MultiCaptionOptions) (Caption, error) {

	c := &MultiCaption{
		captions:    opts.Captions,
		combined:    opts.Combined,
		allow_empty: opts.AllowEmpty,
	}

	return c, nil
}

// Text returns a caption string derived from the base name of 'path'
func (c *MultiCaption) Text(ctx context.Context, source_bucket bucket.Bucket, path string) (string, error) {

	texts := make([]string, len(c.captions))

	for idx, pr := range c.captions {

		txt, err := pr.Text(ctx, source_bucket, path)

		if err != nil {
			return "", fmt.Errorf("Failed to derive text for caption, %w", err)
		}

		if txt != "" && !c.combined {
			return txt, nil
		}

		texts[idx] = txt
	}

	combined := strings.Join(texts, " ")

	if !c.allow_empty && combined == "" {
		return "", fmt.Errorf("Unable to derive caption for %s", path)
	}

	return combined, nil
}
