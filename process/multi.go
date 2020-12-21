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

func (p *MultiProcess) Transform(ctx context.Context, bucket *blob.Bucket, path string) (string, error) {

	final_path := path

	for _, current_p := range p.processes {

		new_path, err := current_p.Transform(ctx, bucket, final_path)

		if err != nil {
			return "", err
		}

		if new_path != "" {
			final_path = new_path
		}
	}

	return final_path, nil
}
