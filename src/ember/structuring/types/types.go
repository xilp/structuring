package types

type Site interface {
	NewTask(info TaskInfo) Task
}

type Task interface {
	Run(appender Appender) error
}

func (p *TaskInfo) Valid() bool {
	return p.Url != ""
}

type Appender func(task TaskInfo) error

type TaskInfo struct {
	Url     string `json:"url"`
	Type    string `json:"type"`
	Weight  int    `json:"weight"`
	Created int64  `json:"created"`
}
