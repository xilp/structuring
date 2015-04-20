package slave

import (
	"ember/cli"
	"ember/http/rpc"
	"ember/structuring/types"
	"errors"
	"flag"
	"fmt"
	"strconv"
	"time"
	"math/rand"
	"os"
	"os/signal"
)

var ErrNoMatchSite = errors.New("no match site")

func Run(args []string) {
	flag := flag.NewFlagSet("slave", flag.ContinueOnError)
	var master string
	var id string
	var concurrent int
	var pause int
	flag.StringVar(&master, "master", "http://127.0.0.1:9000", "master address")
	flag.StringVar(&id, "id", "", "slave id. gen a random one if nil")
	flag.IntVar(&concurrent, "conc", 5, "goroutine number")
	flag.IntVar(&pause, "i", 5, "pause interval if no task, in second")

	cli.ParseFlag(flag, args, "master", "id", "conc")

	if id == "" {
		id = strconv.Itoa(rand.Int())
	}

	slave, err := NewSlave(master, id)
	cli.Check(err)

	rpc := rpc.NewServer()
	err = rpc.Reg(slave)
	if err != nil {
		return
	}
	err = rpc.Run(8888)

	slave.run(concurrent)
}

func (p *Slave) run(concurrent int) {
	p.catchSignal()
	for i := 0; i < concurrent - 1; i++ {
		go p.routine()
	}
	p.routine()
}

func (p *Slave) routine() {
	var err error
	for {
		err = p.invoke()
		if err != nil {
			println(err.Error())
			time.Sleep(time.Second * 3)
		}
	}
}

func (p *Slave) catchSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Printf("received ctrl+c(%v)\n", sig)
			for _, v := range p.sites {
				v.site.Flush()
			}
			os.Exit(0)
		}
	}()
}

func (p *Slave) invoke() (err error) {
	var task types.Task
	var info types.TaskInfo
	for {
		info, err = p.master.Pop(p.id)
		if err != nil {
			return err
		}
		if !info.Valid() {
			return err
		}
		task = p.sites.NewTask(info)
		if task == nil {
			return ErrNoMatchSite
		}
		fmt.Printf("start: %v\n", info)
		err = task.Run(p.append)
		if err != nil {
			return err
		}
		fmt.Printf("done: %v\n", info)
		p.master.Done(p.id, info)
	}
}

func (p *Slave) append(info types.TaskInfo) (err error) {
	return p.master.Push(p.id, info)
}

func NewSlave(addr string, id string) (p *Slave, err error) {
	var master Master
	err = rpc.NewClient(addr).Reg(&master, &MasterTrait{})
	if err != nil {
		return
	}
	p = &Slave{id, NewSites(), master}
	err = p.master.Register("http://127.0.0.1:8888", id)
	if err != nil {
		return
	}
	return
}

func (p *Slave) Trait() map[string][]string {
	return map[string][]string {
		"Search": {"key"},
	}
}

func (p *Slave) Search(key string) (ret string, err error) {
	//TODO 
	println("run", key)
	return "salve:hello json", err
}

type Slave struct {
	id string
	sites Sites
	master Master
}

type Master struct {
	Register func(addr, slave string) error
	Done func(slave string, info types.TaskInfo) error
	Push func(slave string, info types.TaskInfo) error
	Pop func(slave string) (info types.TaskInfo, err error)
}

func (p *MasterTrait) Trait() map[string][]string {
	return map[string][]string {
		"Register": {"addr", "slave"},
		"Done": {"slave", "task"},
		"Push": {"slave", "task"},
		"Pop": {"slave"},
	}
}

type MasterTrait struct {
}
