// package walk provides methods for iterating (walking) over files in a bucket.
package walk

import (
	"context"
	"fmt"
	"io"

	"gocloud.dev/blob"
)

// WalkBucketCallback is a custom function for processing a `blob.ListObject` instance, used by the `WalkBucket` method.
type WalkBucketCallback func(context.Context, *blob.ListObject) error

// WalkBucket will crawl 'bucket' and invoke 'cb' for each file (each `blob.ListObject`) it encounters.
func WalkBucket(ctx context.Context, bucket *blob.Bucket, cb WalkBucketCallback) error {

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

			err = cb(ctx, obj)

			if err != nil {
				return fmt.Errorf("Callback function for %s returned an error, %w", path, err)
			}
		}

		return nil
	}

	err := list(ctx, bucket, "")

	if err != nil {
		return fmt.Errorf("Failed to walk bucket, %w", err)
	}

	return nil
}
