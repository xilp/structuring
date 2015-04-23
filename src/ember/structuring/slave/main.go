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
	args1, args2 := cli.SplitArgs(args, "master", "id", "conc", "i")

	var master string
	var id string
	var concurrent int
	var pause int
	flag := flag.NewFlagSet("slave", flag.ContinueOnError)
	flag.StringVar(&master, "master", "http://127.0.0.1:9000", "master address")
	flag.StringVar(&id, "id", "", "slave id. gen a random one if nil")
	flag.IntVar(&concurrent, "conc", 5, "goroutine number")
	flag.IntVar(&pause, "i", 5, "pause interval if no task, in second")
	cli.ParseFlag(flag, args1, "master", "id", "conc")

	if id == "" {
		id = strconv.Itoa(rand.Int())
	}

	slave, err := NewSlave(master, id)
	cli.Check(err)
	slave.sites.Register("music.163.com")
	slave.run(concurrent)

	client := &types.Slave{}
	hub := cli.NewRpcHub(args2, slave, client)
	hub.Run()
}

func (p *Slave) run(concurrent int) {
	p.catchSignal()
	for i := 0; i < concurrent ; i++ {
		go p.routine()
	}
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
	var master types.Master
	err = rpc.NewClient(addr).Reg(&master)
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

func (p *Slave) Search(key string) (ret [][]string, err error) {
	//TODO
	var x [][]string
	for i, v := range p.sites {
		_ = i
		x, err = v.site.Search(key)
		if err != nil {
			continue
		}
		for _, m := range x {
			ret = append(ret, m)
		}
	}
	return ret, err
}

type Slave struct {
	id string
	sites Sites
	master types.Master
}
