package bucket

import (
	"context"
	"fmt"
	"io"
	"iter"
	"net/url"
	"sort"
	"strings"
	"time"

	"log/slog"

	"github.com/aaronland/go-roster"
)

type Attributes struct {
	ModTime time.Time
	Size    int64
}

// Bucket implements a simple interface for reading and writing Picturebook images to and from
// different storage implementations. It is modeled after the `gocloud.dev/blob.Bucket` interface
// which is what this package used to use. This simplified interface reflects the limited methods
// from the original interface that were used. The goal is to make it easier to implement a variety
// of Picturebook "sources" (or buckets) without having to implement the entirety of the `gocloud.dev/blob.Bucket`
// interface.
type Bucket interface {
	// GatherPictures returns an iterator for listing Picturebook images URIs that can passed to the (bucket implementation's) `NewReader` method.
	GatherPictures(context.Context, ...string) iter.Seq2[string, error]
	// NewReader returns an `io.ReadSeekCloser` instance for a record in the bucket.
	NewReader(context.Context, string, any) (io.ReadSeekCloser, error)
	// NewWriter returns an `io.WriterCloser` instance for writing a record to the bucket.
	NewWriter(context.Context, string, any) (io.WriteCloser, error)
	// Delete removed a record from the bucket.
	Delete(context.Context, string) error
	// Attributes returns an `Attributes` struct for a record in the bucket.
	Attributes(context.Context, string) (*Attributes, error)
	// Close signals the implementation to wrap things up (internally).
	Close() error
}

var bucket_roster roster.Roster

// BucketInitializationFunc is a function defined by individual bucket package and used to create
// an instance of that bucket
type BucketInitializationFunc func(ctx context.Context, uri string) (Bucket, error)

// RegisterBucket registers 'scheme' as a key pointing to 'init_func' in an internal lookup table
// used to create new `Bucket` instances by the `NewBucket` method.
func RegisterBucket(ctx context.Context, scheme string, init_func BucketInitializationFunc) error {

	err := ensureBucketRoster()

	if err != nil {
		return err
	}

	return bucket_roster.Register(ctx, scheme, init_func)
}

func ensureBucketRoster() error {

	if bucket_roster == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		bucket_roster = r
	}

	return nil
}

// NewBucket returns a new `Bucket` instance configured by 'uri'. The value of 'uri' is parsed
// as a `url.URL` and its scheme is used as the key for a corresponding `BucketInitializationFunc`
// function used to instantiate the new `Bucket`. It is assumed that the scheme (and initialization
// function) have been registered by the `RegisterBucket` method.
func NewBucket(ctx context.Context, uri string) (Bucket, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	slog.Info("Scheme", "s", scheme)

	i, err := bucket_roster.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	init_func := i.(BucketInitializationFunc)
	return init_func(ctx, uri)
}

// BucketSchemes returns the list of schemes that have been registered.
func BucketSchemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensureBucketRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range bucket_roster.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
}
