package gofileevent

import (
	"bytes"
	"fmt"

	"golang.org/x/sys/unix"
)

// Event 代表一个文件变更事件
type Event struct {
	Name string // 文件的相对路径
	Op   Op     // 事件的类型
}

// String returns a string representation of the event in the form
// "file: REMOVE|WRITE|..."
func (e Event) String() string {
	return fmt.Sprintf("%q: %s", e.Name, e.Op.String())
}

type Op uint32

// 事件的类型
const (
	Create Op = 1 << iota
	Write
	Remove
	Rename
	Chmod
)

const noteAllEvents = unix.NOTE_DELETE | unix.NOTE_WRITE | unix.NOTE_ATTRIB | unix.NOTE_RENAME

func (op Op) String() string {
	var buffer bytes.Buffer

	if op&Create == Create {
		buffer.WriteString("|CREATE")
	}
	if op&Remove == Remove {
		buffer.WriteString("|REMOVE")
	}
	if op&Write == Write {
		buffer.WriteString("|WRITE")
	}
	if op&Rename == Rename {
		buffer.WriteString("|RENAME")
	}
	if op&Chmod == Chmod {
		buffer.WriteString("|CHMOD")
	}
	if buffer.Len() == 0 {
		return ""
	}
	return buffer.String()[1:]
}
