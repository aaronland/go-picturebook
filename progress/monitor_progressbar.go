package progress

import (
	"context"

	"github.com/schollz/progressbar/v3"
)

// ProgressBarMonitor implements the `Monitor` interface for receiving progress reports using the `schollz/progressbar/v3` package.
type ProgressBarMonitor struct {
	Monitor
	progressbar *progressbar.ProgressBar
}

func init() {

	ctx := context.Background()
	err := RegisterMonitor(ctx, "progressbar", NewProgressBarMonitor)

	if err != nil {
		panic(err)
	}
}

// NewProgressBarMonitor returns a new `ProgressBarMonitor` instance implementing the `Monitor` interface.
func NewProgressBarMonitor(ctx context.Context, uri string) (Monitor, error) {
	m := &ProgressBarMonitor{}
	return m, nil
}

// Signal updates the progress bar with details from 'ev'.
func (m *ProgressBarMonitor) Signal(ctx context.Context, ev *Event) error {

	if m.progressbar == nil {
		m.progressbar = progressbar.Default(int64(ev.Pages))
	}

	m.progressbar.Describe(ev.Message)
	m.progressbar.ChangeMax(ev.Pages)
	m.progressbar.Set(ev.Page)
	return nil
}

// Clear removes the current progress bar.
func (m *ProgressBarMonitor) Clear() error {
	m.progressbar.Reset()
	return m.progressbar.Clear()
}

// Close terminates the progress bar.
func (m *ProgressBarMonitor) Close() error {
	return m.progressbar.Close()
}
