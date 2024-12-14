// package process provides a common interfaces for manipulating images before adding them to a picturebook.
package process

import (
	"context"
	"net/url"

	"github.com/aaronland/go-picturebook/bucket"
	"github.com/aaronland/go-roster"
)

// type Process provides a common interfaces for manipulating images before adding them to a picturebook.
type Process interface {
	// Transform reads a file from a `blob.Bucket` instance, processes it and writes the result to a
	// second `blob.Bucket` instance returning a new filename.
	Transform(context.Context, bucket.Bucket, bucket.Bucket, string) (string, error)
}

// type ProcessInitializeFunc defined a common initialization function for instances implementing the Process interface.
// This is specified when the packages definining those instances call `RegisterProcess` and invoked with the `NewProcess`
// method is called.
type ProcessInitializeFunc func(context.Context, string) (Process, error)

var processes roster.Roster

func ensureRoster() error {

	if processes == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		processes = r
	}

	return nil
}

// RegisterProcess associates a URI scheme with a `ProcessInitializeFunc` initialization function.
func RegisterProcess(ctx context.Context, name string, fn ProcessInitializeFunc) error {

	err := ensureRoster()

	if err != nil {
		return err
	}

	return processes.Register(ctx, name, fn)
}

// NewProcess returns a new `Process` instance for 'uri' whose scheme is expected to have been associated
// with an `ProcessInitializeFunc` (by the `RegisterProcess` method.
func NewProcess(ctx context.Context, uri string) (Process, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := processes.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	fn := i.(ProcessInitializeFunc)

	process, err := fn(ctx, uri)

	if err != nil {
		return nil, err
	}

	return process, nil
}

// AvailableProcess returns the list of schemes that have been registered with `ProcessInitializeFunc` functions.
func AvailableProcesses() []string {
	ctx := context.Background()
	return processes.Drivers(ctx)
}
