package bucket

import (
	"context"
	"io"
	"iter"

	"github.com/aaronland/go-picturebook/picture"
)

// Maybe call it something other than source?
// NewWriter?
// Attributes?

type Bucket interface {
	GatherPictures(context.Context, GatherPicturesProcessFunc, ...string) iter.Seq2[*picture.PictureBookPicture, error]
	NewReader(context.Context, string, any) (io.ReadCloser, error)	
	Close() error
}
