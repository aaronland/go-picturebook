package filter

import (
	"context"

	"gocloud.dev/blob"
)

// type MultiFilter implements the `Filter` interface and allows multiple `Filter` instances to be tested
// to defermine whether an image should  be included in a picturebook.
type MultiFilter struct {
	Filter
	filters []Filter
}

// NewMultiFilter returns a new instance of `MultiFilter` for 'filters'
func NewMultiFilter(ctx context.Context, filters ...Filter) (Filter, error) {

	f := &MultiFilter{
		filters: filters,
	}

	return f, nil
}

// Continues returns a boolean value signaling whether or not 'path' should be included in a picturebook
// by testing each of the `Filter` instances passed to constructor. All filters must return true for this
// method to return true.
func (f *MultiFilter) Continue(ctx context.Context, bucket *blob.Bucket, path string) (bool, error) {

	for _, current_f := range f.filters {

		ok, err := current_f.Continue(ctx, bucket, path)

		if err != nil {
			return false, err
		}

		if !ok {
			return false, nil
		}
	}

	return true, nil
}
