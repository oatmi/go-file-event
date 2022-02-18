package gofileevent

import "sync"

type client struct {
	kq int        // File descriptor (as returned by the kqueue() syscall).
	mu sync.Mutex // Protects access to watcher data
}

func NewClient(path string) *client {
	return &client{}
}

func (c *client) SubscribeEvent(ch chan<- Event, eventSet int) {

}

func (c *client) Close() {

}
