package progress

import (
	"context"
	"net/url"

	"github.com/aaronland/go-roster"
)

// Monitor provides an interface for reporting the progress of a picturebook creation process.
type Monitor interface {
	// Signal() dispatches a processing event to report.
	Signal(context.Context, *Event) error
	// Clear() resets the progress reporting interface.
	Clear() error
	// Close() terminates the progress reporter.
	Close() error
}

// type MonitorInitializeFunc defined a common initialization function for instances implementing the Monitor interface.
// This is specified when the packages definining those instances call `RegisterMonitor` and invoked with the `NewMonitor`
// method is called.
type MonitorInitializeFunc func(context.Context, string) (Monitor, error)

var monitors roster.Roster

func ensureRoster() error {

	if monitors == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		monitors = r
	}

	return nil
}

// RegisterMonitor associates a URI scheme with a `MonitorInitializeFunc` initialization function.
func RegisterMonitor(ctx context.Context, name string, fn MonitorInitializeFunc) error {

	err := ensureRoster()

	if err != nil {
		return err
	}

	return monitors.Register(ctx, name, fn)
}

// NewMonitor returns a new `Monitor` instance for 'uri' whose scheme is expected to have been associated
// with an `MonitorInitializeFunc` (by the `RegisterMonitor` method.
func NewMonitor(ctx context.Context, uri string) (Monitor, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := monitors.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	fn := i.(MonitorInitializeFunc)

	monitor, err := fn(ctx, uri)

	if err != nil {
		return nil, err
	}

	return monitor, nil
}

// AvailableMonitors returns the list of schemes that have been registered with `MonitorInitializeFunc` functions.
func AvailableMonitors() []string {
	ctx := context.Background()
	return monitors.Drivers(ctx)
}
