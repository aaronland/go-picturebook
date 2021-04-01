package application

import (
	"context"
)

// Application is a common interface for all picturebook-related applications.
type Application interface {
	Run(context.Context) error
}
