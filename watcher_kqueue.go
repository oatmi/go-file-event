//go:build freebsd || openbsd || netbsd || dragonfly || darwin
// +build freebsd openbsd netbsd dragonfly darwin

package gofileevent

import (
	"os"
	"path/filepath"
	"time"

	"golang.org/x/sys/unix"
)

var (
	// keventWaitTime 每次从 kevent 读取事件的阻塞时间
	keventWaitTime = durationToTimespec(1 * time.Second)
)

// pathInfo 表示一个路径信息
type pathInfo struct {
	name  string
	isDir bool
}

type Subscription struct {
	Err      chan error
	Events   chan Event
	EventSet Op
}

type kqueueWatcher struct {
	kq          int
	path        string
	files       map[int]pathInfo
	subscribers []Subscription
}

// NewWatcher 创建一个针对制定路径的文件变更事件Watcher
func NewWatcher(path string) (*kqueueWatcher, error) {
	kw := &kqueueWatcher{
		path:  path,
		files: make(map[int]pathInfo),
	}

	kq, err := kqueue()
	if err != nil {
		return nil, err
	}

	kw.kq = kq
	go kw.watch()

	return kw, nil
}

// Subscribe 订阅事件
//
// 可以指定期望得到的文件变更事件，如文件创建、文件写入等
func (w *kqueueWatcher) Subscribe(eventSet Op) (*Subscription, error) {
	subscriber := Subscription{
		Err:      make(chan error),
		Events:   make(chan Event, 100),
		EventSet: eventSet,
	}
	w.subscribers = append(w.subscribers, subscriber)
	return &subscriber, nil
}

// watch 开始监听文件变更
func (w *kqueueWatcher) watch() error {
	fi, err := os.Lstat(w.path)
	if err != nil {
		return err
	}

	files := []string{w.path}

	if fi.IsDir() {
		files, err = fileInDir(w.path)
		if err != nil {
			return err
		}

		for _, file := range files {
			files = append(files, file)
		}
	}

	fds := []int{}
	for _, file := range files {
		fd, err := unix.Open(file, openMode, 0700)
		if fd == -1 {
			return err
		}

		fds = append(fds, fd)
		w.files[fd] = pathInfo{name: file, isDir: false}
	}

	const registerAdd = unix.EV_ADD | unix.EV_CLEAR | unix.EV_ENABLE
	if err := register(w.kq, fds, registerAdd, noteAllEvents); err != nil {
		for _, fd := range fds {
			unix.Close(fd)
		}
		return err
	}

	eventBuffer := make([]unix.Kevent_t, 10)

	for {
		kevents, err := read(w.kq, eventBuffer, &keventWaitTime)
		if err != nil {
			break
		}

		for len(kevents) > 0 {
			kevent := &kevents[0]
			watchfd := int(kevent.Ident)
			mask := uint32(kevent.Fflags)
			path := w.files[watchfd]

			event := newEvent(path.name, mask)

			if len(w.subscribers) > 0 {
				for _, sub := range w.subscribers {
					sub.Events <- event
				}
			}

			kevents = kevents[1:]
		}
	}

	err = unix.Close(w.kq)

	return err
}

// kqueue 创建unix.Kqueue队列
func kqueue() (kq int, err error) {
	kq, err = unix.Kqueue()
	if kq == -1 {
		return kq, err
	}
	return kq, nil
}

// read retrieves pending events, or waits until an event occurs.
// A timeout of nil blocks indefinitely, while 0 polls the queue.
func read(kq int, events []unix.Kevent_t, timeout *unix.Timespec) ([]unix.Kevent_t, error) {
	n, err := unix.Kevent(kq, nil, events, timeout)
	if err != nil {
		return nil, err
	}
	return events[0:n], nil
}

func durationToTimespec(d time.Duration) unix.Timespec {
	return unix.NsecToTimespec(d.Nanoseconds())
}

// register 在kq上注册感兴趣的事件
func register(kq int, fds []int, flags int, fflags uint32) error {
	changes := make([]unix.Kevent_t, len(fds))

	for i, fd := range fds {
		// SetKevent converts int to the platform-specific types:
		unix.SetKevent(&changes[i], fd, unix.EVFILT_VNODE, flags)
		changes[i].Fflags = fflags
	}

	// register the events
	success, err := unix.Kevent(kq, changes, nil, nil)
	if success == -1 {
		return err
	}
	return nil
}

// newEvent 从kqueue的fflags解析出事件的信息
func newEvent(name string, mask uint32) Event {
	e := Event{Name: name}
	if mask&unix.NOTE_DELETE == unix.NOTE_DELETE {
		e.Op |= Remove
	}
	if mask&unix.NOTE_WRITE == unix.NOTE_WRITE {
		e.Op |= Write
	}
	if mask&unix.NOTE_RENAME == unix.NOTE_RENAME {
		e.Op |= Rename
	}
	if mask&unix.NOTE_ATTRIB == unix.NOTE_ATTRIB {
		e.Op |= Chmod
	}
	return e
}

// fileInDir 便利读出文件夹中的所有文件
func fileInDir(path string) ([]string, error) {
	var files []string

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return files, nil
}
