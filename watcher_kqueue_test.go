package gofileevent

import (
	"fmt"
	"testing"
	"time"
)

func TestKqueueWatcher(t *testing.T) {
	dirname := "test"
	w, _ := NewWatcher(dirname)
	sub, _ := w.Subscribe(Create | Write | Remove)

	go func() {
		for e := range sub.Events {
			fmt.Printf("%s\n", e)
		}
	}()

	time.Sleep(time.Minute)
}
