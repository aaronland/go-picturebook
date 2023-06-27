// Package encode provides methods for encoding images.
package encode

import (
	"context"
	"fmt"
	"image"
	"io"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/aaronland/go-roster"
)

var encoders roster.Roster

// EncoderInitializationFunc is a function defined by individual encoder package and used to create
// an instance of that encoder.
type InitializeEncoderFunc func(context.Context, string) (Encoder, error)

// Encoder is an interface for writing data to multiple sources or targets.
type Encoder interface {
	// Encode encode and writes an `image.Image` instance to a `io.Writer` instance.
	Encode(context.Context, io.Writer, image.Image) error
}

func ensureRoster() error {

	if encoders == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return fmt.Errorf("Failed to create new encoder roster, %w", err)
		}

		encoders = r
	}

	return nil
}

// RegisterEncoder registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `Encoder` instances by the `NewEncoder` method.
func RegisterEncoder(ctx context.Context, f InitializeEncoderFunc, schemes ...string) error {

	err := ensureRoster()

	if err != nil {
		return err
	}

	for _, s := range schemes {

		err := encoders.Register(ctx, s, f)

		if err != nil {
			return err
		}
	}

	return nil
}

// NewEncoder returns a new `Encoder` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `EncoderInitializationFunc`
// function used to instantiate the new `Encoder`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterEncoder` method.
func NewEncoder(ctx context.Context, uri string) (Encoder, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	ext := filepath.Ext(u.Path)
	scheme := strings.TrimLeft(ext, ".")

	enc_u := url.URL{}
	enc_u.Path = u.Path
	enc_u.RawQuery = u.RawQuery

	i, err := encoders.Driver(ctx, scheme)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive encoder for %s, %w", scheme, err)
	}

	if i == nil {
		return nil, fmt.Errorf("Undefined encoder for %s", scheme)
	}

	f := i.(InitializeEncoderFunc)
	return f(ctx, enc_u.String())
}
