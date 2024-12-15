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

	"github.com/aaronland/go-roster"
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
