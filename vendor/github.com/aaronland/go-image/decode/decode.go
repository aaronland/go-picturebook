// Package decode provides a common interface for decoding image file handles.
package decode

import (
	"context"
	"image"
	"io"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/aaronland/go-roster"
)

type InitializeDecoderFunc func(context.Context, string) (Decoder, error)

type Decoder interface {
	Decode(context.Context, io.ReadSeeker) (image.Image, string, error)
}

var decoders roster.Roster

func ensureRoster() error {

	if decoders == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		decoders = r
	}

	return nil
}

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

func NewDecoder(ctx context.Context, uri string) (Decoder, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(uri)
	scheme := strings.TrimLeft(ext, ".")

	dec_u := url.URL{}
	dec_u.Path = u.Path
	dec_u.RawQuery = u.RawQuery

	i, err := decoders.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	f := i.(InitializeDecoderFunc)
	return f(ctx, dec_u.String())
}
