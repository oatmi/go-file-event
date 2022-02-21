//go:build freebsd || openbsd || netbsd || dragonfly || darwin
// +build freebsd openbsd netbsd dragonfly darwin

package gofileevent

import (
	"errors"

	"golang.org/x/sys/unix"
)

var ErrorUnsubscribed = errors.New("unsubscribed")

type KqueueWatcher struct {
	kq int
}

// watch 监听文件或者文件夹的变化事件
func (w *KqueueWatcher) watch(file string) error {
	for false {
	}
	return nil
}

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

func NewWatcher(file string) (*KqueueWatcher, error) {
	kw := &KqueueWatcher{}

	kq, err := kqueue()
	if err != nil {
		return nil, err
	}
	kw.kq = kq

	go kw.watch(file)

	return kw, nil
}

// kqueue creates a new kernel event queue and returns a descriptor.
func kqueue() (kq int, err error) {
	kq, err = unix.Kqueue()
	if kq == -1 {
		return kq, err
	}
	return kq, nil
}
