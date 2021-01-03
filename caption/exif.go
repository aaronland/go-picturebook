package caption

import (
	"context"
	"errors"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"gocloud.dev/blob"
	"net/url"
)

func init() {

	ctx := context.Background()
	err := RegisterCaption(ctx, "exif", NewExifCaption)

	if err != nil {
		panic(err)
	}

	exif.RegisterParsers(mknote.All...)
}

type ExifCaption struct {
	Caption
	property string
}

func NewExifCaption(ctx context.Context, uri string) (Caption, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	property := q.Get("property")

	if property != "datetime" {
		return nil, errors.New("Unsupported EXIF field")
	}

	c := &ExifCaption{
		property: property,
	}

	return c, nil
}

func (c *ExifCaption) Text(ctx context.Context, bucket *blob.Bucket, path string) (string, error) {

	fh, err := bucket.NewReader(ctx, path, nil)

	if err != nil {
		return "", err
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
