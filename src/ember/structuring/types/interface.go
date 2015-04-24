package types

import (
	"time"
)

type Slave struct {
	Search func(key string)(ret [][]string, err error) `args:"key" return:"ret"`
}

type Master struct {
	Register func(addr, slave string) error `args:"addr,slave"`
	Pop func(slave string) (info TaskInfo, err error) `args:"slave" return:"info"`
	Fetch func(url string) error `args:"url"`
	Done func(slave string, info TaskInfo) error `args:"slave,info"`
	Push func(slave string, info TaskInfo) error `args:"slave,info"`

	Search func(key string)(ret [][][]string, err error) `args:"key" return:"ret"`
	Dones func() (urls []string, err error) `return:"urls"`
	Doings func() (tasks []TaskInfo, err error) `return:"tasks"`
	Tasks func() (tasks []TaskInfo, err error) `return:"tasks"`
}

type Site interface {
	NewTask(info TaskInfo) Task
	Search(key string) (ret [][]string, err error)
	Close() (err error)
}

type Task interface {
	Run(appender Appender) error
}

type Appender func(info TaskInfo) error

func NewTaskInfo(url string, typ string, weight int) TaskInfo {
	return TaskInfo{url, typ, weight, time.Now().UnixNano()}
}

func (p *TaskInfo) Valid() bool {
	return p.Url != ""
}

type TaskInfo struct {
	Url     string `json:"url"`
	Type    string `json:"type"`
	Weight  int    `json:"weight"`
	Created int64  `json:"created"`
}
