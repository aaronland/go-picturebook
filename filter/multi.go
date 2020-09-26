package filter

import (
	"context"
	"gocloud.dev/blob"
)

type MultiFilter struct {
	Filter
	filters []Filter
}

func NewMultiFilter(ctx context.Context, filters ...Filter) (Filter, error) {

	f := &MultiFilter{
		filters: filters,
	}

	return f, nil
}

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
