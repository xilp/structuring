package master

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"sync"
	"time"
	"ember/cli"
	"ember/http/rpc"
	"ember/structuring/slave"
	"ember/structuring/types"
)

func Run(args []string) {
	flag := flag.NewFlagSet("run", flag.ContinueOnError)
	var port int
	flag.IntVar(&port, "port", 9000, "master service port")
	cli.ParseFlag(flag, args)

	p := NewMaster()
	rpc := rpc.NewServer()
	err := rpc.Reg(p)
	if err != nil {
		return
	}
	err = rpc.Run(port)
	cli.Check(err)
	p.scan()
}

func (p *Master) Fetch(url string) error {
	return p.Push("master", types.NewTaskInfo(url, "index", math.MaxInt64))
}

func (p *Master) Dones() (urls []string, err error) {
	p.locker.Lock()
	defer p.locker.Unlock()

	for url, _ := range p.dones {
		urls = append(urls, url)
	}
	return
}

func (p *Master) Slaves() (slaves []string, err error) {
	p.locker.Lock()
	defer p.locker.Unlock()

	for slave, _ := range p.slaves {
		slaves = append(slaves, slave)
	}
	return
}

func (p *Master) Done(slave string, info types.TaskInfo) (err error) {
	p.locker.Lock()
	defer p.locker.Unlock()
	p.slaves[slave] = time.Now().UnixNano()

	delete(p.doings, info.Url)
	p.dones[info.Url] = true

	p.save()
	return
}

var count = 0
func (p *Master) Push(slave string, info types.TaskInfo) (err error) {
	fmt.Printf("appending %v\n", info)
	p.locker.Lock()

	count ++;
	fmt.Printf("[count:%d]\n", count)
	defer p.locker.Unlock()
	p.slaves[slave] = time.Now().UnixNano()

	if p.dones[info.Url] {
		return
	}

	if _, ok := p.doings[info.Url] ; ok {
		return
	}

	p.tasks = append(p.tasks, info)

	p.save()
	return
}

func (p *Master) Pop(slave string) (info types.TaskInfo, err error) {
	p.locker.Lock()
	defer p.locker.Unlock()
	p.slaves[slave] = time.Now().UnixNano()

	if len(p.tasks) == 0 {
		err = ErrNoTask
		return
	}
	info = p.tasks[0]
	p.tasks = p.tasks[1:]

	p.doings[info.Url] = info

	p.save()
	return
}

func (p *Master) save() (err error) {
	// TODO
	return
}

func (p *Master) load() (err error) {
	// TODO
	return
}

func (p *Master) scan() {
	go func() {
		for k, v := range p.doings {
			if time.Since(time.Unix(time.Unix(0, v.Created).Unix(), 0)) >= time.Minute {
				p.Push("repush", v)
				delete(p.doings, k)
			}
		}
		time.Sleep(time.Minute)
	}()
}

func (p *Master) Trait() map[string][]string {
	st := slave.MasterTrait{}
	trait := st.Trait()
	trait["Fetch"] = []string{"url"}
	trait["Slaves"] = []string{}
	trait["Dones"] = []string{}
	return trait
}

func NewMaster() *Master {
	p := &Master {
		dones: make(map[string]bool),
		doings: make(map[string]types.TaskInfo),
		slaves: make(map[string]int64),
	}
	p.load()
	return p
}

type Master struct {
	tasks []types.TaskInfo
	dones map[string]bool
	doings map[string]types.TaskInfo
	slaves map[string]int64
	locker sync.Mutex
}

var ErrNoTask = errors.New("no task")
