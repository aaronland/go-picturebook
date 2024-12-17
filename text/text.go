// package text provides a common interface for different mechanisms to derive texts for images.
package text

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aaronland/go-picturebook/bucket"
	"github.com/aaronland/go-roster"
)

// type Text provides a common interface for different mechanisms to derive texts for images.
type Text interface {
	// Text produces a text derived from a file contained in a gocloud.dev/blob Bucket instance.
	Body(context.Context, bucket.Bucket, string) (string, error)
}

// type TextInitializeFunc defined a common initialization function for instances implementing the Text interface.
// This is specified when the packages definining those instances call `RegisterText` and invoked with the `NewText`
// method is called.
type TextInitializeFunc func(context.Context, string) (Text, error)

var texts roster.Roster

func ensureRoster() error {

	if texts == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return fmt.Errorf("Failed to create new roster for texts, %w", err)
		}

		texts = r
	}

	return nil
}

// RegisterText associates a URI scheme with a `TextInitializeFunc` initialization function.
func RegisterText(ctx context.Context, name string, fn TextInitializeFunc) error {

	err := ensureRoster()

	if err != nil {
		return fmt.Errorf("Failed to ensure texts roster, %w", err)
	}

	return texts.Register(ctx, name, fn)
}

// NewText returns a new `Text` instance for 'uri' whose scheme is expected to have been associated
// with an `TextInitializeFunc` (by the `RegisterText` method.
func NewText(ctx context.Context, uri string) (Text, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI for NewText, %w", err)
	}

	scheme := u.Scheme

	i, err := texts.Driver(ctx, scheme)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive driver for '%s' text scheme, %w", scheme, err)
	}

	fn := i.(TextInitializeFunc)

	text, err := fn(ctx, uri)

	if err != nil {
		return nil, fmt.Errorf("TextInitializeFunc failed, %w", err)
	}

	return text, nil
}

// AvailableText returns the list of schemes that have been registered with `TextInitializeFunc` functions.
func AvailableTexts() []string {
	ctx := context.Background()
	return texts.Drivers(ctx)
}
