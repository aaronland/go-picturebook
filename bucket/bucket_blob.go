package bucket

import (
	"context"
	"fmt"
	"io"
	"iter"

	aa_bucket "github.com/aaronland/gocloud-blob/bucket"
	"gocloud.dev/blob"
	"sync"
)

var bucket_mu = new(sync.Map)

type BlobBucket struct {
	Bucket
	bucket *blob.Bucket
}

func init() {

	ctx := context.Background()
	err := RegisterGoCloudBuckets(ctx)

	if err != nil {
		panic(err)
	}
}

// RegisterGoCloudBuckets will explicitly register all the schemes associated with the `gocloud.dev/blob.Bucket` interface.
func RegisterGoCloudBuckets(ctx context.Context) error {

	for _, scheme := range blob.DefaultURLMux().BucketSchemes() {

		_, exists := bucket_mu.LoadOrStore(scheme, true)

		if exists {
			continue
		}

		err := RegisterBucket(ctx, scheme, NewBlobBucket)

		if err != nil {
			return fmt.Errorf("Failed to register scheme '%s', %w", scheme, err)
		}
	}

	return nil
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

func (s *BlobBucket) GatherPictures(ctx context.Context, uris ...string) iter.Seq2[string, error] {

	return func(yield func(string, error) bool) {

		for _, uri := range uris {
			for p, err := range s.gatherPictures(ctx, uri) {
				yield(p, err)
			}
		}
	}
}

func (s *BlobBucket) gatherPictures(ctx context.Context, uri string) iter.Seq2[string, error] {

	var list func(context.Context, *blob.Bucket, string) error

	return func(yield func(string, error) bool) {

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

				yield(path, nil)

				/*
					pic, err := process_func(ctx, path)

					if err != nil {
						return err
					}

					if pic == nil {
						continue
					}

					yield(pic, nil)
				*/
			}

			return nil
		}

		err := list(ctx, s.bucket, uri)

		if err != nil {
			yield("", err)
		}
	}
}

func (s *BlobBucket) NewReader(ctx context.Context, key string, opts any) (io.ReadSeekCloser, error) {

	r, err := s.bucket.NewReader(ctx, key, nil)
	return r, err

	/*
		if err != nil {
			return nil, err
		}

		return ioutil.NewReadSeekCloser(r)
	*/
}

func (s *BlobBucket) NewWriter(ctx context.Context, key string, opts any) (io.WriteCloser, error) {
	return s.bucket.NewWriter(ctx, key, nil)
}

func (s *BlobBucket) Attributes(ctx context.Context, path string) (*Attributes, error) {

	blob_attrs, err := s.bucket.Attributes(ctx, path)

	if err != nil {
		return nil, err
	}

	attrs := &Attributes{
		ModTime: blob_attrs.ModTime,
		Size:    blob_attrs.Size,
	}

	return attrs, nil
}

func (s *BlobBucket) Delete(ctx context.Context, key string) error {
	return s.bucket.Delete(ctx, key)
}

func (s *BlobBucket) Close() error {
	return s.bucket.Close()
}
