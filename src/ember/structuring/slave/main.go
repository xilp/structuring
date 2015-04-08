package slave

import (
	"flag"
	"os"
	"ember/cli"
	"ember/http/rpc"
)

func Run(args []string) {
	flag := flag.NewFlagSet("slave", flag.ContinueOnError)
	var master string
	var id string
	var cocurrent int
	var pause int
	flag.VarStr(&master, "master", "127.0.0.1", "master address")
	flag.VarStr(&id, "id", "", "slave id. gen a random one if nil")
	flag.VarInt(&cocurrent, "coc", 5, "goroutine number")
	flag.VarInt(&pause, "i", 5, "pause interval if no task, in second")

	cli.ParseFlag(flag, args, "master", "id", "coc")

	slave, err := NewSlave(master, id)
	cli.Check(err)
	slave.Run(concurrent)
}

func (p *Slave) Run(cocurrent int) {
	for i := 0; i < cocurrent {
		go p.run()
	}
}

func (p *Slave) run() (err error) {
	for {
		info, err := p.tasks.Pop()
		if !info.Valid() {
			time.Sleep(time.Second * 3)
		}
		task := p.site.NewTask(info)
		err = task.Run(p.tasks)
		if err != nil {
			return
		}
	}
}

func NewSlave(master string, id string) (p *Slave, err error) {
	var master Master
	err = rpc.Reg(master, task, trait)
	if err != nil {
		return
	}
	p = &Slave{site, master}
	return
}

type Slave struct {
	sites map[string]Site
	master Master
}

type Master struct {
	Done(task TaskInfo) error
	Push(task TaskInfo) error
	Pop() (task TaskInfo, err error)
}

type TaskInfo {
	Url string `json:"url"`
	Type string `json:"type"`
	Weight int `json:"weight"`
	Created int64 `json:"created"`
}
