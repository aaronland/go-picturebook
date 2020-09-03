package sort

import (
	"context"
	"github.com/aaronland/go-picturebook/picture"
)

type Sorter interface {
	Sort(context.Context, []*picture.PictureBookPicture) error
}
