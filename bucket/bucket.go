package bucket

import (
	"context"
	"io"
	"iter"
	"time"
)

type Attributes struct {
	ModTime time.Time
	Size    int64
}

type Bucket interface {
	GatherPictures(context.Context, ...string) iter.Seq2[string, error]
	NewReader(context.Context, string, any) (io.ReadSeekCloser, error)
	NewWriter(context.Context, string, any) (io.WriteCloser, error)
	Delete(context.Context, string) error
	Attributes(context.Context, string) (*Attributes, error)
	Close() error
}
