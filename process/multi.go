package process

import (
	"context"
	"gocloud.dev/blob"
)

type MultiProcess struct {
	Process
	processes []Process
}

func NewMultiProcess(ctx context.Context, processes ...Process) (Process, error) {

	p := &MultiProcess{
		processes: processes,
	}

	return p, nil
}

func (p *MultiProcess) Transform(ctx context.Context, source_bucket *blob.Bucket, target_bucket *blob.Bucket, path string) (string, error) {

	final_path := path

	for _, current_p := range p.processes {

		new_path, err := current_p.Transform(ctx, source_bucket, target_bucket, final_path)

		if err != nil {
			return "", err
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
