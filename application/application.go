package application

import (
	"context"
)

type Application interface {
	Run(context.Context) error
}
