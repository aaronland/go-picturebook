package text

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"gocloud.dev/blob"
)

func init() {

	ctx := context.Background()
	err := RegisterText(ctx, "json", NewJsonText)

	if err != nil {
		panic(err)
	}
}

// type JsonText implements the `Text` interface and returns empty text strings.
type JsonText struct {
	Text
	texts_table map[string]string
}

// NewJsonText return a new instance of `JsonText` for 'url'
func NewJsonText(ctx context.Context, uri string) (Text, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URL, %w", err)
	}

	texts_path := u.Path

	texts_r, err := os.Open(texts_path)

	if err != nil {
		return nil, fmt.Errorf("Failed to open %s for reading, %w", texts_path, err)
	}

	defer texts_r.Close()

	var texts_table map[string]string

	dec := json.NewDecoder(texts_r)
	err = dec.Decode(&texts_table)

	if err != nil {
		return nil, fmt.Errorf("Failed to decode %s, %w", texts_path, err)
	}

	c := &JsonText{
		texts_table: texts_table,
	}
	return c, nil
}

// Body returns ...
func (c *JsonText) Body(ctx context.Context, bucket *blob.Bucket, path string) (string, error) {

	text, exists := c.texts_table[path]

	if !exists {
		return "", nil
	}

	return text, nil
}
