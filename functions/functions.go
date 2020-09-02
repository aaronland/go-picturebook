package functions

import (
	"context"
)

type PictureBookFilterFunc func(context.Context, string) (bool, error)

type PictureBookPreProcessFunc func(context.Context, string) (string, error)

type PictureBookCaptionFunc func(context.Context, string) (string, error)

func DefaultFilterFunc(ctx context.Context, path string) (bool, error) {
	return true, nil
}
