package process

import (
	"context"
	"fmt"

	"gocloud.dev/blob"
)

// type MultiProcess implements the `Process` interface and allows multiple `Process` instances to
// to tranform an image before adding it to a picturebook.
type MultiProcess struct {
	Process
	processes []Process
}

// NewMultiProcess returns a new instance of `MultiProcess` for 'processes'
func NewMultiProcess(ctx context.Context, processes ...Process) (Process, error) {

	p := &MultiProcess{
		processes: processes,
	}

	return p, nil
}

// Tranform applies the `Tranform` method for all its internal `Process` instances. All processes must succeed
// in order for this method to succeed.
func (p *MultiProcess) Transform(ctx context.Context, source_bucket *blob.Bucket, target_bucket *blob.Bucket, path string) (string, error) {

	final_path := path

	for _, current_p := range p.processes {

		new_path, err := current_p.Transform(ctx, source_bucket, target_bucket, final_path)

		if err != nil {
			return "", fmt.Errorf("Failed to execute transformation for %s with %v, %w", path, p, err)
		}

		if new_path != "" && new_path != path {

			final_path = new_path

			// See this - it's important. Because we're in a loop we need to make sure
			// we know where to find the second (and third...) temporary file that's been
			// processed which will be in the "target" (or temporary) bucket and not the
			// source bucket. (20210223/straup)

			source_bucket = target_bucket
		}
	}

	return final_path, nil
}
