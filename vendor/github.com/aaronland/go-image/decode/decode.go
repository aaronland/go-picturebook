// Package decode provides methods for decoding images.
package decode

import (
	"context"
	"fmt"
	"image"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/aaronland/go-roster"
)

func DecodeFromPath(ctx context.Context, path string) (image.Image, string, error) {

	dec, err := NewDecoder(ctx, path)

	if err != nil {
		return nil, "", fmt.Errorf("Failed to create new decoder, %w", err)
	}

	r, err := os.Open(path)

	if err != nil {
		return nil, "", fmt.Errorf("Failed to open path for reading, %w", err)
	}

	defer r.Close()

	return dec.Decode(ctx, r)
}

var decoders roster.Roster

// DecoderInitializationFunc is a function defined by individual decoder package and used to create
// an instance of that decoder.
type InitializeDecoderFunc func(context.Context, string) (Decoder, error)

type Decoder interface {
	// Decode decodes an `io.ReaderSeeker` instance and returns an `image.Image` instance.
	Decode(context.Context, io.ReadSeeker) (image.Image, string, error)
}

func ensureRoster() error {

	if decoders == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return fmt.Errorf("Failed to create new decoder roster, %w", err)
		}

		decoders = r
	}

	return nil
}

// RegisterDecoder registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `Decoder` instances by the `NewDecoder` method.
func RegisterDecoder(ctx context.Context, f InitializeDecoderFunc, schemes ...string) error {

	err := ensureRoster()

	if err != nil {
		return err
	}

	for _, s := range schemes {

		err := decoders.Register(ctx, s, f)

		if err != nil {
			return err
		}
	}

	return nil
}

// NewDecoder returns a new `Decoder` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `DecoderInitializationFunc`
// function used to instantiate the new `Decoder`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterDecoder` method.
func NewDecoder(ctx context.Context, uri string) (Decoder, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	ext := filepath.Ext(uri)
	scheme := strings.TrimLeft(ext, ".")

	dec_u := url.URL{}
	dec_u.Path = u.Path
	dec_u.RawQuery = u.RawQuery

	i, err := decoders.Driver(ctx, scheme)

	if err != nil {
		return nil, fmt.Errorf("Failed to derive decoder for %s, %w", scheme, err)
	}

	if i == nil {
		return nil, fmt.Errorf("Undefined decoder for %s", scheme)
	}

	f := i.(InitializeDecoderFunc)
	return f(ctx, dec_u.String())
}
