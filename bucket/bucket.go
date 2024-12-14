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
	Size    int64
}

// type GatherPicturesProcessFunc defines a method for processing the path to an image file in to a `picture.PictureBookPicture` instance.
type GatherPicturesProcessFunc func(context.Context, string) (*picture.PictureBookPicture, error)

type Bucket interface {
	GatherPictures(context.Context, GatherPicturesProcessFunc, ...string) iter.Seq2[*picture.PictureBookPicture, error]
	NewReader(context.Context, string, any) (io.ReadSeekCloser, error)
	NewWriter(context.Context, string, any) (io.WriteCloser, error)
	Delete(context.Context, string) error
	Attributes(context.Context, string) (*Attributes, error)
	Close() error
}
