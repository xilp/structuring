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
	var path string
	var host, port string
	flag := flag.NewFlagSet("slave", flag.ContinueOnError)
	flag.StringVar(&master, "master", "http://127.0.0.1:9000", "master address")
	flag.StringVar(&id, "id", "", "slave id. gen a random one if nil")
	flag.IntVar(&concurrent, "conc", 5, "goroutine number")
	flag.IntVar(&pause, "i", 5, "pause interval if no task, in second")
	flag.StringVar(&path, "path", "data", "path to store data")
	flag.StringVar(&host, "host", "127.0.0.1", "slave ip register to master")
	flag.StringVar(&port, "port", "8888", "slave port regist to master")
	cli.ParseFlag(flag, args1, "master", "id", "conc")

	if id == "" {
		id = strconv.Itoa(rand.Int())
	}

	slave, err := NewSlave(master, id, path, host, port)
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
			for _, v := range p.sites.sites {
				v.site.Flush()
			}
			os.Exit(0)
		}
	}()
}

func (p *Slave) invoke() (err error) {
	var task types.Task
	var info types.TaskInfo
	var i = 0
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
		i ++
		if i > 100 {
			i = 0
			p.master.Register("http://" + p.host + ":" + p.port, p.id)
		}
	}
}

func (p *Slave) append(info types.TaskInfo) (err error) {
	return p.master.Push(p.id, info)
}

func NewSlave(addr string, id string, path string, host, port string) (p *Slave, err error) {
	var master types.Master
	err = rpc.NewClient(addr).Reg(&master)
	if err != nil {
		return
	}
	p = &Slave{id, host, port, path, NewSites(path), master}
	err = p.master.Register("http://" + host + ":" + port, id)
	if err != nil {
		return
	}
	return
}

func (p *Slave) Search(key string) (ret [][]string, err error) {
	//TODO
	var x [][]string
	for i, v := range p.sites.sites {
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
	host string
	port string
	path string
	sites Sites
	master types.Master
}
