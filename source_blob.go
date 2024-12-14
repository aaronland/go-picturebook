package picturebook

import (
	"context"
	"fmt"
	"io"
	"iter"

	"github.com/aaronland/go-picturebook/picture"
	"github.com/aaronland/gocloud-blob/bucket"
	"gocloud.dev/blob"
)

type BlobSource struct {
	Source
	bucket *blob.Bucket
}

func NewBlobSource(ctx context.Context, uri string) (Source, error) {

	b, err := bucket.OpenBucket(ctx, uri)

	if err != nil {
		return nil, err
	}

	s := &BlobSource{
		bucket: b,
	}

	return s, nil
}

func (s *BlobSource) GatherPictures(ctx context.Context, process_func GatherPicturesProcessFunc, uris ...string) iter.Seq2[*picture.PictureBookPicture, error] {

	return func(yield func(*picture.PictureBookPicture, error) bool) {

		for _, uri := range uris {
			for p, err := range s.gatherPictures(ctx, process_func, uri) {
				yield(p, err)
			}
		}
	}
}

func (s *BlobSource) gatherPictures(ctx context.Context, process_func GatherPicturesProcessFunc, uri string) iter.Seq2[*picture.PictureBookPicture, error] {

	var list func(context.Context, *blob.Bucket, string) error

	return func(yield func(*picture.PictureBookPicture, error) bool) {

		list = func(ctx context.Context, bucket *blob.Bucket, prefix string) error {

			iter := bucket.List(&blob.ListOptions{
				Delimiter: "/",
				Prefix:    prefix,
			})

			for {
				obj, err := iter.Next(ctx)

				if err == io.EOF {
					break
				}

				if err != nil {
					return fmt.Errorf("Failed to iterate next in bucket for %s, %w", prefix, err)
				}

				path := obj.Key

				if obj.IsDir {

					err := list(ctx, bucket, path)

					if err != nil {
						return fmt.Errorf("Failed to list bucket for %s, %w", path, err)
					}

					continue
				}

				pic, err := process_func(ctx, path)

				if err != nil {
					return err
				}

				if pic == nil {
					continue
				}

				yield(pic, nil)
			}

			return nil
		}

		err := list(ctx, s.bucket, uri)

		if err != nil {
			yield(nil, err)
		}
	}
}

func (s *BlobSource) Close() error {
	return s.bucket.Close()
}
