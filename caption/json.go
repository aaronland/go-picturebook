package caption

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/aaronland/go-picturebook/bucket"
)

func init() {

	ctx := context.Background()
	err := RegisterCaption(ctx, "json", NewJsonCaption)

	if err != nil {
		panic(err)
	}
}

// type JsonCaption implements the `Caption` interface and returns empty caption strings.
type JsonCaption struct {
	Caption
	captions_table map[string]string
}

// NewJsonCaption return a new instance of `JsonCaption` for 'url'
func NewJsonCaption(ctx context.Context, uri string) (Caption, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URL, %w", err)
	}

	captions_path := u.Path

	captions_r, err := os.Open(captions_path)

	if err != nil {
		return nil, fmt.Errorf("Failed to open %s for reading, %w", captions_path, err)
	}

	defer captions_r.Close()

	var captions_table map[string]string

	dec := json.NewDecoder(captions_r)
	err = dec.Decode(&captions_table)

	if err != nil {
		return nil, fmt.Errorf("Failed to decode %s, %w", captions_path, err)
	}

	c := &JsonCaption{
		captions_table: captions_table,
	}
	return c, nil
}

// Text returns an empty caption string
func (c *JsonCaption) Text(ctx context.Context, source_bucket bucket.Bucket, path string) (string, error) {

	text, exists := c.captions_table[path]

	if !exists {
		return "", nil
	}

	return text, nil
}
