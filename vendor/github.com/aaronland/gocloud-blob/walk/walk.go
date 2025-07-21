// package walk provides methods for iterating (walking) over files in a bucket.
package walk

import (
	"context"
	"fmt"
	"io"
	"iter"

	"gocloud.dev/blob"
)

// WalkBucketCallback is a custom function for processing a `blob.ListObject` instance, used by the `WalkBucket` method.
type WalkBucketCallback func(context.Context, *blob.ListObject) error

// WalkBucket will crawl 'bucket' and invoke 'cb' for each file (each `blob.ListObject`) it encounters.
func WalkBucket(ctx context.Context, bucket *blob.Bucket, cb WalkBucketCallback) error {

	return WalkBucketWithPrefix(ctx, bucket, "", cb)
}

// WalkBucketWithPrefix will crawl 'bucket' for files parented by 'prefix' and invoke 'cb' for each file (each `blob.ListObject`) it encounters.
func WalkBucketWithPrefix(ctx context.Context, bucket *blob.Bucket, prefix string, cb WalkBucketCallback) error {

	for obj, err := range WalkBucketWithIter(ctx, bucket, prefix) {

		if err != nil {
			return err
		}

		err = cb(ctx, obj)

		if err != nil {
			return err
		}
	}

	return nil
}

// WalkBucketWithIter will iterate 'bucket' for files parented by 'prefix' and yield an `iter.Seq2[*blob.ListObject, error]` instance for each file it encounters. 
func WalkBucketWithIter(ctx context.Context, bucket *blob.Bucket, prefix string) iter.Seq2[*blob.ListObject, error] {

	return func(yield func(obj *blob.ListObject, err error) bool) {

		var list func(context.Context, *blob.Bucket, string) error

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

					if !yield(nil, err) {
						return nil
					}
				}

				path := obj.Key

				if obj.IsDir {

					err := list(ctx, bucket, path)

					if err != nil {
						if !yield(nil, fmt.Errorf("Failed to list bucket for %s, %w", path, err)) {
							return nil
						}
					}

					continue
				}

				if !yield(obj, nil) {
					return nil
				}

			}

			return nil
		}

		err := list(ctx, bucket, prefix)

		if err != nil {
			yield(nil, err)
		}
	}
}
