package caption

import (
	"context"
	"fmt"
	"net/url"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"gocloud.dev/blob"
)

func init() {

	ctx := context.Background()
	err := RegisterCaption(ctx, "exif", NewExifCaption)

	if err != nil {
		panic(err)
	}

	exif.RegisterParsers(mknote.All...)
}

// type ExifCaption implements the `Caption` interface and derives caption text from EXIF properties.
type ExifCaption struct {
	Caption
	property string
}

// NewExifCaption return a new instance of `ExifCaption` for 'uri'
func NewExifCaption(ctx context.Context, uri string) (Caption, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse URL, %w", err)
	}

	q := u.Query()

	property := q.Get("property")

	if property != "datetime" {
		return nil, fmt.Errorf("Unsupported EXIF field '%s'", property)
	}

	c := &ExifCaption{
		property: property,
	}

	return c, nil
}

// Text returns a caption string derived from EXIF data in 'path'
func (c *ExifCaption) Text(ctx context.Context, bucket *blob.Bucket, path string) (string, error) {

	fh, err := bucket.NewReader(ctx, path, nil)

	if err != nil {
		return "", fmt.Errorf("Failed to create new bucket, %w", err)
	}

	defer fh.Close()

	x, err := exif.Decode(fh)

	if err != nil {
		return "", nil
	}

	dt, err := x.DateTime()

	if err != nil {
		return "", nil
	}

	return dt.Format("January 02, 2006"), nil
}
