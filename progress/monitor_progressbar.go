package progress

import (
	"context"

	"github.com/schollz/progressbar/v3"
)

type ProgressbarMonitor struct {
	Monitor
	progressbar *progressbar.ProgressBar
}

func NewProgressbarMonitor(ctx context.Context, uri string) (Monitor, error) {
	m := &ProgressbarMonitor{}
	return m, nil
}

func (m *ProgressbarMonitor) Signal(ctx context.Context, ev *ProgressEvent) error {

	if m.progressbar == nil {
		m.progressbar = progressbar.Default(int64(ev.Pages))
	}

	m.progressbar.ChangeMax(ev.Pages)
	m.progressbar.Set(ev.Page)
	return nil
}

func (m *ProgressbarMonitor) Close() error {
	return nil
}
