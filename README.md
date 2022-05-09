# 基于 kqueue api 的文件变更通知组件

受 [fsnotify/fsnotify](https://github.com/fsnotify/fsnotify) 和 [fis3](https://fis.baidu.com/) 启发使用 go 开发的基于 kqueue 的文件文件变更监听组建，更多平台的轮训 api 计划中。

| Adapter | OS     |Stat                                    |
| ------  | ------ | -------------------------------------------- |
| kqueue                | BSD, macOS, iOS\*                | Supported |
| inotify               | Linux 2.6.27 or later, Android\* | Planned  |
| Polling               | *All*                            | Planned | 
| ReadDirectoryChangesW | Windows                          | Maybe |
| USN Journals          | Windows                          | Maybe                     |

## 基础能力

1. 支持对文件夹下的所有文件进行监听；
2. 多方订阅，一个监听支持多个订阅者，订阅者可以只订阅自己感兴趣的事件；

## Usage

```go
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
```

```
❯ go test -run TestKqueueWatcher
"test/bfoo.txt": WRITE
"test/bfoo.txt": WRITE
"test/bfoo.txt": CHMOD
"test/bfoo.txt": CHMOD
"test": WRITE
"test": WRITE
"test/foo.txt": WRITE
"test/foo.txt": WRITE
"test/foo.txt": CHMOD
"test/foo.txt": CHMOD
"test/foo.txt": REMOVE
"test/foo.txt": REMOVE
"test": WRITE
"test": WRITE
"test": WRITE
"test": WRITE
```

## 关联项目

- [fsnotify](https://github.com/fsnotify/fsnotify)
