package slave

import (
	"flag"
	"time"
	"ember/cli"
	"ember/http/rpc"
	"ember/structuring/types"
)

func Run(args []string) {
	flag := flag.NewFlagSet("slave", flag.ContinueOnError)
	var master string
	var id string
	var concurrent int
	var pause int
	flag.StringVar(&master, "master", "127.0.0.1", "master address")
	flag.StringVar(&id, "id", "", "slave id. gen a random one if nil")
	flag.IntVar(&concurrent, "conc", 5, "goroutine number")
	flag.IntVar(&pause, "i", 5, "pause interval if no task, in second")

	cli.ParseFlag(flag, args, "master", "id", "conc")

	slave, err := NewSlave(master, id)
	cli.Check(err)
	slave.run(concurrent)
}

func (p *Slave) run(concurrent int) {
	for i := 0; i < concurrent; i++ {
		go p.routine()
	}
}

func (p *Slave) routine() (err error) {
	for {
		info, err := p.master.Pop(p.id)
		if err != nil {
			return err
		}
		if !info.Valid() {
			time.Sleep(time.Second * 3)
		}
		task := p.sites.NewTask(info)
		err = task.Run(p.append)
		if err != nil {
			return err
		}
	}
}

func (p *Slave) append(task types.TaskInfo) (err error) {
	// TODO
	return
}

func NewSlave(addr string, id string) (p *Slave, err error) {
	var master Master
	err = rpc.NewClient(addr).Reg(&master, &MasterTrait{})
	if err != nil {
		return
	}
	p = &Slave{id, NewSites(), master}
	return
}

type Slave struct {
	id string
	sites Sites
	master Master
}

type Master struct {
	Done func(slave string, task types.TaskInfo) error
	Push func(slave string, task types.TaskInfo) error
	Pop func(slave string) (task types.TaskInfo, err error)
}

func (p *MasterTrait) Trait() map[string][]string {
	return map[string][]string {
		"Done": {"slave", "task"},
		"Push": {"slave", "task"},
		"Pop": {"slave"},
	}
}

type MasterTrait struct {
}
