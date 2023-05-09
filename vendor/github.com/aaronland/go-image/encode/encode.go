package encode

import (
	"context"
	"image"
	"io"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/aaronland/go-roster"
)

type InitializeEncoderFunc func(context.Context, string) (Encoder, error)

type Encoder interface {
	Encode(context.Context, io.Writer, image.Image) error
}

var encoders roster.Roster

func ensureRoster() error {

	if encoders == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		encoders = r
	}

	return nil
}

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

func NewEncoder(ctx context.Context, uri string) (Encoder, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(uri)
	scheme := strings.TrimLeft(ext, ".")

	enc_u := url.URL{}
	enc_u.Path = u.Path
	enc_u.RawQuery = u.RawQuery

	i, err := encoders.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	f := i.(InitializeEncoderFunc)
	return f(ctx, enc_u.String())
}
