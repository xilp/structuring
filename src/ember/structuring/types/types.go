package types

import (
	"time"
)

type Site interface {
	NewTask(info TaskInfo) Task
	Flush() (err error)
	Search(key string) (ret [][]string, err error)
}

type Task interface {
	Run(appender Appender) error
}

type SongInfo struct {
	Version string
	Url string
	SongName, Singer, Album, IssueDate string
	IssueCompany, Note, SongLyric string
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
