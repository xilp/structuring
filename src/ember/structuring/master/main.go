package master

import (
	"errors"
	//"fmt"
	"math"
	"sync"
	"time"
	"ember/cli"
	"ember/http/rpc"
	"ember/structuring/types"
)

func Run(args []string) {
	master := NewMaster()
	master.catchSignal()
	master.scan()

	client := &types.Master{}
	hub := cli.NewRpcHub(args, master, client)
	hub.Run()
}

func (p *Master) Fetch(url string) error {
	return p.Push("master", types.NewTaskInfo(url, "index", math.MaxInt64))
}

func (p *Master) Search(key string) (ret [][][]string, err error) {
	var res [][]string
	var x [][][]string
	for i, _ := range p.slavesRemote {
		if i != "master" && i != "rpush" {
			res, err = p.slavesRemote[i].Search(key)
			if err != nil {
				println(err.Error())
			} else {
				x = append(x, res)
			}

		}
	}
	return x, err
}

func (p *Master) Done(slave string, info types.TaskInfo) (err error) {
	p.locker.Lock()
	defer p.locker.Unlock()
	p.slaves[slave] = time.Now().UnixNano()

	delete(p.doings, info.Url)
	p.dones[info.Url] = true

	return
}

func (p *Master) Push(slave string, info types.TaskInfo) (err error) {
	//fmt.Printf("appending %v\n", info)
	p.locker.Lock()
	defer p.locker.Unlock()

	p.slaves[slave] = time.Now().UnixNano()

	if p.dones[info.Url] {
		return
	}

	if p.todos[info.Url] {
		return
	}

	if _, ok := p.doings[info.Url] ; ok {
		return
	}

	p.tasks = append(p.tasks, info)
	p.todos[info.Url] = true

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
	delete(p.todos, info.Url)

	p.doings[info.Url] = info

	return
}

func (p *Master) Register(addr, slave string) (err error) {
	sc := &types.Slave{}
	err = rpc.NewClient(addr).Reg(sc)
	if err != nil {
		return
	}
	p.slavesRemote[slave] = sc
	return
}

func NewMaster() *Master {
	p := &Master {
		todos: make(map[string]bool),
		dones: make(map[string]bool),
		doings: make(map[string]types.TaskInfo),
		slaves: make(map[string]int64),
		slavesRemote: make(map[string]*types.Slave),
		donesFile: NewData("donesFile.txt"),
		doingsFile: NewData("doingFile.txt"),
		tasksFile: NewData("tasksFile.txt"),
	}
	p.load()
	return p
}

type Master struct {
	tasks []types.TaskInfo
	todos map[string]bool
	dones map[string]bool
	doings map[string]types.TaskInfo
	searchTasks map[string]types.TaskInfo
	slaves map[string]int64
	slavesRemote map[string]*types.Slave
	donesFile	Data
	doingsFile	Data
	tasksFile	Data
	locker sync.Mutex
}

var ErrNoTask = errors.New("no task")
