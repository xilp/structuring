package slave

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"ember/cli"
	"ember/http/rpc"
	"ember/structuring/types"
	m1c "ember/structuring/sites/music.163.com"
)

func Run(args []string) {
	args1, args2 := cli.SplitArgs(args, "master", "id", "conc", "i", "path", "host", "port")

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

	m1c, err := m1c.New(path)
	cli.Check(err)

	slave, err := NewSlave(master, id, host, port)
	cli.Check(err)
	slave.Reg("music.163.com", m1c)

	slave.Run(concurrent)

	client := &types.Slave{}
	hub := cli.NewRpcHub(args2, slave, client)
	go hub.Run()

	holder := make(chan os.Signal)
	signal.Notify(holder, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	sig := <-holder
	slave.Close()
}

func (p *Slave) Run(concurrent int) {
	for i := 0; i < concurrent ; i++ {
		go p.routine()
	}
}

func (p *Slave) Close() {
	for _, site := range p.sites {
		site.Close()
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

func (p *Slave) report() {
	p.master.Register("http://" + p.host + ":" + p.port, p.id)
	for _ = range time.Tick(time.Second * 3) {
		p.master.Register("http://" + p.host + ":" + p.port, p.id)
	}
}

func (p *Slave) append(info types.TaskInfo) (err error) {
	return p.master.Push(p.id, info)
}

func NewSlave(addr string, id string, host string, port string) (p *Slave, err error) {
	if id == "" {
		id = strconv.Itoa(rand.Int())
	}

	var master types.Master
	err = rpc.NewClient(addr).Reg(&master)
	if err != nil {
		return
	}
	p = &Slave{id, host, port, NewSites(), master}

	go p.report()
	return
}

func (p *Slave) Search(key string) (ret [][]string, err error) {
	var x [][]string
	for i, site := range p.sites {
		_ = i
		x, err = site.Search(key)
		if err != nil {
			continue
		}
		for _, m := range x {
			ret = append(ret, m)
		}
	}
	return ret, err
}

func (p *Slave) Reg(domain string, site types.Site) {
	p.sites.Reg(domain, site)
}

type Slave struct {
	id string
	host string
	port string
	sites Sites
	master types.Master
}

var ErrNoMatchSite = errors.New("no match site")

