package gofileevent

type EventType uint32

const (
	EventCreate EventType = 1 << iota
	EventWrite
	EventRemove
	EventRename
	EventChmod
)

type Event struct {
	filePath string
	et       EventType
}

func (e *Event) File() string {
	return e.filePath
}

func (e *Event) Type() EventType {
	return e.et
}

type Subscription interface {
	Unsubscribe()
	Err() <-chan error
	Receive() Event
}

type Watcher interface {
	watch(file string) error

	Subscribe(ch chan<- Event, topic int) (Subscription, error)
}
