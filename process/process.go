package process

import (
	"context"
	"github.com/aaronland/go-roster"
	"gocloud.dev/blob"
	"net/url"
)

type Process interface {
	Transform(context.Context, *blob.Bucket, string) (string, error)
}

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

func RegisterProcess(ctx context.Context, name string, fn ProcessInitializeFunc) error {

	err := ensureRoster()

	if err != nil {
		return err
	}

	return processes.Register(ctx, name, fn)
}

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

func AvailableProcesses() []string {
	ctx := context.Background()
	return processes.Drivers(ctx)
}
