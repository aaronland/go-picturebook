package bucket

import (
	"context"
	"fmt"
	"io"
	"iter"

	"github.com/aaronland/go-picturebook/picture"
	aa_bucket "github.com/aaronland/gocloud-blob/bucket"
	"gocloud.dev/blob"
)

// type GatherPicturesProcessFunc defines a method for processing the path to an image file in to a `picture.PictureBookPicture` instance.
type GatherPicturesProcessFunc func(context.Context, string) (*picture.PictureBookPicture, error)

type BlobBucket struct {
	Bucket
	bucket *blob.Bucket
}

func NewBlobBucket(ctx context.Context, uri string) (Bucket, error) {

	b, err := aa_bucket.OpenBucket(ctx, uri)

	if err != nil {
		return nil, err
	}

	s := &BlobBucket{
		bucket: b,
	}

	return s, nil
}

func (s *BlobBucket) GatherPictures(ctx context.Context, process_func GatherPicturesProcessFunc, uris ...string) iter.Seq2[*picture.PictureBookPicture, error] {

	return func(yield func(*picture.PictureBookPicture, error) bool) {

		for _, uri := range uris {
			for p, err := range s.gatherPictures(ctx, process_func, uri) {
				yield(p, err)
			}
		}
	}
}

func (s *BlobBucket) gatherPictures(ctx context.Context, process_func GatherPicturesProcessFunc, uri string) iter.Seq2[*picture.PictureBookPicture, error] {

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

func (s *BlobBucket) NewReader(ctx context.Context, key string, opts any) (io.ReadCloser, error) {
	return s.bucket.NewReader(ctx, key, nil)
}

func (s *BlobBucket) Close() error {
	return s.bucket.Close()
}
