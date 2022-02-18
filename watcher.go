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
	FilePath string
	Type     EventType
}

type Watcher interface {
	Watch(file string) error
}

type Subscription interface {
	Unsubscribe()
	Err() <-chan error
}

type Subscriber interface {
	Subscribe(ch chan<- Event) (Subscription, error)
}
