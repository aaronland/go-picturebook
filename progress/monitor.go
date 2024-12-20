package progress

import (
	"context"
)

type Monitor interface {
	Signal(context.Context, *ProgressEvent) error
	Close() error
}
