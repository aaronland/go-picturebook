package process

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aaronland/go-picturebook/bucket"
)

func init() {

	ctx := context.Background()
	err := RegisterProcess(ctx, "null", NewNullProcess)

	if err != nil {
		panic(err)
	}
}

// type NullProcess implements the `Process` interface but does not apply any transformations to an image.
type NullProcess struct {
	Process
}

// NullProcess returns a new instance of `NullProcess` for 'uri' which must be parsable as a valid `net/url` URL instance.
func NewNullProcess(ctx context.Context, uri string) (Process, error) {

	_, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI for NewNullProcess, %w", err)
	}

	f := &NullProcess{}

	return f, nil
}

// Tranform is a no-op, does not apply any tranformations to 'path' and returns an empty string.
func (f *NullProcess) Transform(ctx context.Context, source_bucket bucket.Bucket, target_bucket bucket.Bucket, path string) (string, error) {
	return "", nil
}
