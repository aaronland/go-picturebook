package encode

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/url"
	"strconv"
)

type JPEGEncoder struct {
	Encoder
	quality int
}

func init() {

	ctx := context.Background()
	RegisterEncoder(ctx, NewJPEGEncoder, "jpg", "jpeg")
}

func NewJPEGEncoder(ctx context.Context, uri string) (Encoder, error) {

	quality := 100

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URI, %w", err)
	}

	q := u.Query()

	q_quality := q.Get("quality")

	if q_quality != "" {

		v, err := strconv.Atoi(q_quality)

		if err != nil {
			return nil, fmt.Errorf("Invalid ?quality= parameter, %w", err)
		}

		quality = v
	}

	e := &JPEGEncoder{
		quality: quality,
	}

	return e, nil
}

func (e *JPEGEncoder) Encode(ctx context.Context, wr io.Writer, im image.Image) error {

	opts := &jpeg.Options{
		Quality: e.quality,
	}

	return jpeg.Encode(wr, im, opts)
}
