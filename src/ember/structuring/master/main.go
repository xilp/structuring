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
	p.catchSignal()
	p.scan()
	err = rpc.Run(port)
	cli.Check(err)
}

func (p *Master) Fetch(url string) error {
	return p.Push("master", types.NewTaskInfo(url, "index", math.MaxInt64))
}

func (p *Master) Search(key string) (ret string, err error) {
	str := ""
	//fmt.Printf("[master][Search]\n")
	//fmt.Printf("[master][Search][p.slavesRemote:%#v]\n", p.slavesRemote)
	//for i, v := range p.slaves {
	for i, v := range p.slavesRemote {
		fmt.Printf("[i:%#v][v:%#v]\n", i, v)
		if i != "master" && i != "rpush" {
			ret, err = p.slavesRemote[i].Search(key)
			//fmt.Printf("[err:%#v]\n", err)
			if err != nil {
			} else {
				str = str + ret
			}
		}
	}
	//fmt.Printf("[master][Search][str:%s]\n", str)
	return "master:" + str, err
}

func (p *Master) Done(slave string, info types.TaskInfo) (err error) {
	p.locker.Lock()
	defer p.locker.Unlock()
	p.slaves[slave] = time.Now().UnixNano()

	delete(p.doings, info.Url)
	p.dones[info.Url] = true

	return
}

var count = 0
func (p *Master) Push(slave string, info types.TaskInfo) (err error) {
	fmt.Printf("appending %v\n", info)
	p.locker.Lock()
	defer p.locker.Unlock()

	count ++;
	fmt.Printf("[count:%d]\n", count)
	p.slaves[slave] = time.Now().UnixNano()

	if p.dones[info.Url] {
		return
	}

	if _, ok := p.doings[info.Url] ; ok {
		return
	}

	p.tasks = append(p.tasks, info)

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

	return
}

func (p *Master) Register(addr, slaveName string) (err error) {
	var slave Slave
	err = rpc.NewClient(addr).Reg(&slave, &SlaveTrait{})
	if err != nil {
		return
	}
	p.slavesRemote[slaveName] = slave
	return
}

func (p *Master) Trait() map[string][]string {
	st := slave.MasterTrait{}
	trait := st.Trait()
	trait["Fetch"] = []string{"url"}
	trait["Search"] = []string{"key"}
	trait["Slaves"] = []string{}
	trait["Dones"] = []string{}
	return trait
}

func NewMaster() *Master {
	p := &Master {
		dones: make(map[string]bool),
		doings: make(map[string]types.TaskInfo),
		slaves: make(map[string]int64),
		slavesRemote: make(map[string]Slave),
		donesFile: NewData("donesFile.txt"),
		doingsFile: NewData("doingFile.txt"),
		tasksFile: NewData("tasksFile.txt"),
		slavesFile: NewData("slavesFile.txt"),
	}
	p.load()
	return p
}

type Master struct {
	tasks []types.TaskInfo
	dones map[string]bool
	doings map[string]types.TaskInfo
	searchTasks map[string]types.TaskInfo
	slaves map[string]int64
	slavesRemote map[string]Slave
	donesFile	Data
	doingsFile	Data
	tasksFile	Data
	slavesFile	Data
	locker sync.Mutex
}

type Slave struct {
	Search func(key string)(ret string, err error)
}

func (p *SlaveTrait) Trait() map[string][]string {
	return map[string][]string {
		"Search": {"key"},
	}
}

type SlaveTrait struct {
}

var ErrNoTask = errors.New("no task")
