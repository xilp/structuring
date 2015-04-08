package m1c

import (
	"ember/structuring/types"
)

func (p *Song) Run(appender types.Appender) (err error) {
	// TODO
	// extrat infos
	// get similars
	// new song task -> chan
	return
}

type Song struct {
	url string
}

func (p *Site) NewTask(info types.TaskInfo) types.Task {
	switch info.Type {
	case "song":
		return &Song{info.Url}
	}
	return nil
}

func New() *Site {
	return &Site{}
}

type Site struct {
}
