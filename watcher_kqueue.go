package gofileevent

import "errors"

var ErrorUnsubscribed = errors.New("unsubscribed")

type KqueueWatcher struct {
}

func (w *KqueueWatcher) watch(file string) {}

func (w *KqueueWatcher) Subscribe(ch chan<- Event, eventSet int) (Subscription, error) {
	return kqueueSubscription{
		err:   make(chan error),
		event: make(chan Event),
	}, nil
}

type kqueueSubscription struct {
	err   chan error
	event chan Event
}

func (s kqueueSubscription) Unsubscribe() {
	s.err <- ErrorUnsubscribed
	close(s.err)
}

func (s kqueueSubscription) Err() <-chan error {
	return s.err
}

func (s kqueueSubscription) Receive() Event {
	e := <-s.event
	return e
}

func NewWatcher() KqueueWatcher {
	kw := KqueueWatcher{}
	kw.watch(".")

	return kw
}
