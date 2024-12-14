package picturebook

import (
	"context"

	"github.com/go-pdf/fpdf"
)

type Target interface {
	Save(context.Context, string, *fpdf.Fpdf) error
	Close() error
}
