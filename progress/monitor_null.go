package progress

import (
	"context"
)

// NullMonitor implements the `Monitor` interface for receiving progress reports but not doing anything with them.
type NullMonitor struct {
	Monitor
}

func init() {

	ctx := context.Background()
	err := RegisterMonitor(ctx, "null", NewNullMonitor)

	if err != nil {
		panic(err)
	}
}

// NewNullMonitor returns a new `NullMonitor` instance implementing the `Monitor` interface.
func NewNullMonitor(ctx context.Context, uri string) (Monitor, error) {
	m := &NullMonitor{}
	return m, nil
}

// Signal receives 'ev' but doesn't do anything with it.
func (m *NullMonitor) Signal(ctx context.Context, ev *Event) error {
	return nil
}

// Clear doesn't do anything.
func (m *NullMonitor) Clear() error {
	return nil
}

// Close doesn't do anything.
func (m *NullMonitor) Close() error {
	return nil
}
