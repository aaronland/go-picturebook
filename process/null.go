package process

import (
	"context"
	"net/url"
)

func init() {

	ctx := context.Background()
	err := RegisterProcess(ctx, "null", NewNullProcess)

	if err != nil {
		panic(err)
	}
}

type NullProcess struct {
	Process
}

func NewNullProcess(ctx context.Context, uri string) (Process, error) {

	_, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	f := &NullProcess{}

	return f, nil
}

func (f *NullProcess) Continue(ctx context.Context, path string) (string, error) {
	return "", nil
}
