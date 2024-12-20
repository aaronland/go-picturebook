package progress

import (
	"context"
)

type NullMonitor struct {
	Monitor
}

func NewNullMonitor(ctx context.Context, uri string) (Monitor, error) {
	m := &NullMonitor{}
	return m, nil
}

func (m *NullMonitor) Signal(ctx context.Context, ev *ProgressEvent) error {
	return nil
}
