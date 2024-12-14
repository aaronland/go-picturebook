package picturebook

import (
	"context"
	"io"
	"iter"

	"github.com/aaronland/go-picturebook/picture"
)

type Source interface {
	GatherPictures(context.Context, GatherPicturesProcessFunc, ...string) iter.Seq2[*picture.PictureBookPicture, error]
	NewReader(context.Context, string) (io.ReadSeekCloser, error)
	Close() error
}
