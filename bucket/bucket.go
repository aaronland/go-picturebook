package bucket

import (
	"context"
	"io"
	"iter"
	"time"

	"github.com/aaronland/go-picturebook/picture"
)

type Attributes struct {
	ModTime time.Time
}

// type GatherPicturesProcessFunc defines a method for processing the path to an image file in to a `picture.PictureBookPicture` instance.
type GatherPicturesProcessFunc func(context.Context, string) (*picture.PictureBookPicture, error)

type Bucket interface {
	GatherPictures(context.Context, GatherPicturesProcessFunc, ...string) iter.Seq2[*picture.PictureBookPicture, error]

	// This needs to implement ModTime and Size...
	// https://pkg.go.dev/gocloud.dev/blob#Reader
	// https://github.com/google/go-cloud/blob/v0.40.0/blob/blob.go#L99
	NewReader(context.Context, string, any) (io.ReadSeekCloser, error)
	
	NewWriter(context.Context, string, any) (io.WriteCloser, error)
	Attributes(context.Context, string) (*Attributes, error)
	Close() error
}
