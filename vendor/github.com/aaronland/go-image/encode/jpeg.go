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

// JPEGEncoder is a struct that implements the `Encoder` interface for
// encoding JPEG images.
type JPEGEncoder struct {
	Encoder
	quality int
}

func init() {

	ctx := context.Background()
	RegisterEncoder(ctx, NewJPEGEncoder, "jpg", "jpeg")
}

// NewJPEGEncoder returns a new `JPEGEncoder` instance.
// 'uri' in the form of:
//
//	/path/to/image.jpg?{OPTIONS}
//
// Where {OPTIONS} are:
//   - ?quality={QUALITY} - an optional value specifying the quality of the
//     JPEG image; default is 100.
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

// Encode will encode 'im' using the `image/jpeg` package and write the results to 'wr'
func (e *JPEGEncoder) Encode(ctx context.Context, wr io.Writer, im image.Image) error {

	opts := &jpeg.Options{
		Quality: e.quality,
	}

	return jpeg.Encode(wr, im, opts)
}
