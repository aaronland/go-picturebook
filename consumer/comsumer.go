package consumer

import (
	"context"
	"github.com/aaronland/go-picturebook/picture"
	"github.com/aaronland/go-roster"
	"net/url"
)

type Consumer interface {
	GatherPictures(context.Context, ...string) ([]*picture.PictureBookPicture, error)
}

type ConsumerInitializeFunc func(context.Context, string) (Consumer, error)

var consumers roster.Roster

func ensureRoster() error {

	if consumers == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		consumers = r
	}

	return nil
}

func RegisterConsumer(ctx context.Context, name string, fn ConsumerInitializeFunc) error {

	err := ensureRoster()

	if err != nil {
		return err
	}

	return consumers.Register(ctx, name, fn)
}

func NewConsumer(ctx context.Context, uri string) (Consumer, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := consumers.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	fn := i.(ConsumerInitializeFunc)

	consumer, err := fn(ctx, uri)

	if err != nil {
		return nil, err
	}

	return consumer, nil
}

func AvailableConsumers() []string {
	ctx := context.Background()
	return consumers.Drivers(ctx)
}
