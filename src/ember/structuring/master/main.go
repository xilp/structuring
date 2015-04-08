package structuring

import (
	"os"
	"ember/structuring/slave"
)

func Run() {
	if len(os.Args) < 2 {
		println("usage: bin url")
		os.Exit(1)
	}

	master := NewMaster(site.New())
	err := master.Run()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}

func (p *Server) Run() {
	task := site.NewTask()
	p.tasks.Push(task)
}

func (p *Master) Done(task slave.TaskInfo) error {
}

func (p *Master) Push(task slave.TaskInfo) error {
}

func (p *Master) Pop() (task slave.TaskInfo, err error) {
}

func (p *Master) Heartbeat(slave string) error {
}

func NewMaster() *Master {
	return &Master{}
}

type Master struct {
	tasks []slave.TaskInfo
	counter chan interface{}
}
