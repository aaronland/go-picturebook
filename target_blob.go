package picturebook

import (
	"context"

	"github.com/aaronland/gocloud-blob/bucket"
	"github.com/go-pdf/fpdf"
	"gocloud.dev/blob"
)

type BlobTarget struct {
	bucket *blob.Bucket
}

func NewBlobTarget(ctx context.Context, uri string) (Target, error) {

	b, err := bucket.OpenBucket(ctx, uri)

	if err != nil {
		return nil, err
	}

	t := &BlobTarget{
		bucket: b,
	}

	return t, nil
}

func (t *BlobTarget) Save(ctx context.Context, path string, doc *fpdf.Fpdf) error {

	wr, err := t.bucket.NewWriter(ctx, path, nil)

	if err != nil {
		return err
	}

	err = doc.Output(wr)

	if err != nil {
		return err
	}

	return wr.Close()
}

func (t *BlobTarget) Close() error {
	return t.bucket.Close()
}
