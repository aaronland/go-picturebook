package bucket

import (
	"context"
	"fmt"
	"io"
	"iter"
	_ "log/slog"
	"strings"

	aa_bucket "github.com/aaronland/gocloud/blob/bucket"
	aa_walk "github.com/aaronland/gocloud/blob/walk"
	"gocloud.dev/blob"
	"sync"
)

var bucket_mu = new(sync.Map)

// BlobBucket implements the `Bucket` interface using a `gocloud.dev/blob.Bucket` instance.
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

// NewBlobBucket returns a new instantiation of the `Bucket` interface using a `gocloud.dev/blob.Bucket` instance.
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

// GatherPictures will return a iterator listing items in 'uris' (parented by 'b').
func (b *BlobBucket) GatherPictures(ctx context.Context, uris ...string) iter.Seq2[string, error] {

	return func(yield func(string, error) bool) {

		for _, uri := range uris {

			// Strip the leading slash because that's what gocloud.dev/blob likes
			walk_uri := strings.TrimLeft(uri, "/")

			for obj, err := range aa_walk.WalkBucketWithIter(ctx, b.bucket, walk_uri) {

				if err != nil {
					if !yield("", err) {
						break
					}
				}

				// Put the leading slash back because that's what everything else likes
				path := fmt.Sprintf("/%s", obj.Key)

				if !yield(path, nil) {
					break
				}
			}
		}
	}
}

// NewReader returns an `io.ReadSeekCloser` instance for the record named 'key' in 'b'.
func (b *BlobBucket) NewReader(ctx context.Context, key string, opts any) (io.ReadSeekCloser, error) {

	// exists, err := b.bucket.Exists(ctx, key)
	// slog.Info("WTF", "key", key, "exists", exists, "error", err)

	r, err := b.bucket.NewReader(ctx, key, nil)
	return r, err
}

// NewWriter returns an `io.WriterCloser` instance for writing to the record named 'key' in 'b'.
func (b *BlobBucket) NewWriter(ctx context.Context, key string, opts any) (io.WriteCloser, error) {
	return b.bucket.NewWriter(ctx, key, nil)
}

// Attributes returns an `Attributes` struct for the record named 'key' in'b'.
func (b *BlobBucket) Attributes(ctx context.Context, path string) (*Attributes, error) {

	blob_attrs, err := b.bucket.Attributes(ctx, path)

	if err != nil {
		return nil, err
	}

	attrs := &Attributes{
		ModTime: blob_attrs.ModTime,
		Size:    blob_attrs.Size,
	}

	return attrs, nil
}

// Delete removes the record named 'key' in 'b'.
func (b *BlobBucket) Delete(ctx context.Context, key string) error {
	return b.bucket.Delete(ctx, key)
}

// Close tells 'b' to wrap things up.
func (b *BlobBucket) Close() error {
	return b.bucket.Close()
}
