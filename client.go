package gofileevent

import "sync"

type client struct {
	path string
	kq   int        // File descriptor (as returned by the kqueue() syscall).
	mu   sync.Mutex // Protects access to watcher data
}

func NewClient(path string) *client {
	return &client{
		path: path,
		kq:   1,
		mu:   sync.Mutex{},
	}
}

func (c *client) SubscribeEvent(ch chan<- Event, eventSet int) {

}

func (c *client) Close() {

}
