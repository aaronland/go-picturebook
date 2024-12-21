package progress

// Event defines details about a (picturebook) page being processed.
type Event struct {
	// The total number of pages being processed
	Pages int
	// The current page being processed
	Page int
	// An optional string describing the event
	Message string
}

// NewEvent returns a new `Event` instance derived from 'page' and 'count'.
func NewEvent(page int, count int) *Event {
	return &Event{
		Page:  page,
		Pages: count,
	}
}
