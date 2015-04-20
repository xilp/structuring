package master

import (
	"errors"
	"flag"
	"fmt"
	"math"
	"sync"
	"time"
	"encoding/json"
	"ember/cli"
	"ember/http/rpc"
	"ember/structuring/slave"
	"ember/structuring/types"
	"regexp"
	"os"
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


func (p *Master) Tasks() (tasks []types.TaskInfo, err error) {
	return p.tasks, err
}

func (p *Master) UnTasks(tasks []types.TaskInfo) (err error) {
	p.tasks = tasks
	return err
}

func (p *Master) Dones() (urls []string, err error) {
	for url, _ := range p.dones {
		urls = append(urls, url)
	}
	return
}

func (p *Master) UnDones(urls []string) (err error) {
	for _, v := range urls {
		p.dones[v] = true
	}
	return err
}

func (p *Master) Slaves() (slaves map[string]int64, err error) {
	return p.slaves , err
}

func (p *Master) UnSlaves(slaves map[string]int64) (err error) {
	if slaves != nil {
		p.slaves = slaves
	}
	return
}

func (p *Master) donesSerialize(urls []string) (str string, err error) {
	for _, v := range urls {
		str = str + v + "\n"
	}
	return str, err
}

func (p *Master) donesUnSerialize(str string) (url []string, err error) {
	reg := regexp.MustCompile(`http://[^\n]+`)
	url = reg.FindAllString(str, -1)
	return url, err
}

func (p *Master) doingsSerialize(tasks []types.TaskInfo) (str string, err error) {
	for _, v := range tasks {
		data, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		str = str + string(data) + "\n"
	}
	return str, err
}

func (p *Master) doingsUnSerialize(str string) (tasks []types.TaskInfo, err error) {
	reg := regexp.MustCompile(`[^\n]+`)
	tasksJson := reg.FindAllString(str, -1)

	task := types.TaskInfo {}
	for _, v := range tasksJson {
		err := json.Unmarshal([]byte(v), &task)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, err
}

func (p *Master) Doings() (tasks []types.TaskInfo, err error) {
	for _, v := range p.doings {
		tasks = append(tasks, v)
	}
	return tasks, err
}

func (p *Master) UnDoings(tasks []types.TaskInfo) (err error) {
	//TODO
	for _, v := range tasks {
		p.doings[v.Url] = v
	}
	return
}

func (p *Master) tasksSerialize(tasks []types.TaskInfo) (str string, err error) {
	for _, v := range tasks {
		data, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		str = str + string(data) + "\n"
	}
	return str, err
}

func (p *Master) tasksUnSerialize(str string) (tasks []types.TaskInfo, err error) {
	reg := regexp.MustCompile(`[^\n]+`)
	tasksJson := reg.FindAllString(str, -1)

	task := types.TaskInfo {}
	for _, v := range tasksJson {
		err := json.Unmarshal([]byte(v), &task)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, err
}

func (p *Master) slavesSerialize(slaves map[string]int64) (str string, err error) {
	data, err := json.Marshal(slaves)
	if err != nil {
		return "", err
	}
	str = string(data) + "\n"
	return str, err
}

func (p *Master) slavesUnSerialize(str string) ( slaves map[string]int64, err error) {
	reg := regexp.MustCompile(`[^\n]+`)
	slavesJson := reg.FindAllString(str, -1)

	for _, v := range slavesJson {
		err := json.Unmarshal([]byte(v), &slaves)
		if err != nil {
			return nil, err
		}
	}
	return slaves, err
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

func (p *Master) save() (err error) {
	p.locker.Lock()
	defer p.locker.Unlock()

	//time.Sleep(1000*1000*1000*10)
	dones, err:= p.Dones()
	if err != nil {
		return err
	}

	slaves, err:= p.Slaves()
	if err != nil {
		return err
	}
	_ = slaves

	doings, err:= p.Doings()
	if err != nil {
		return err
	}

	tasks, err:= p.Tasks()
	if err != nil {
		return err
	}

	donesStr, err:= p.donesSerialize(dones)
	if err != nil {
		return
	}
	doingsStr, err:= p.doingsSerialize(doings)
	if err != nil {
		return
	}
	tasksStr, err:= p.tasksSerialize(tasks)
	if err != nil {
		return
	}
	slavesStr, err:= p.slavesSerialize(slaves)
	if err != nil {
		return
	}

	p.donesFile.write	(donesStr	, 0)
	p.doingsFile.write	(doingsStr	, 0)
	p.tasksFile.write	(tasksStr	, 0)
	p.slavesFile.write	(slavesStr	, 0)

	return
}

func (p *Master) load() (err error) {
	p.locker.Lock()
	defer p.locker.Unlock()

	donesStr, err := p.donesFile.read(0)
	if err != nil {
		return
	}
	doingsStr, err := p.doingsFile.read(0)
	if err != nil {
		return
	}
	slavesStr, err := p.slavesFile.read(0)
	if err != nil {
		return
	}
	tasksStr, err := p.tasksFile.read(0)
	if err != nil {
		return
	}
	//fmt.Printf("[tasksStr:%#v]\n", tasksStr)

	dones, err := p.donesUnSerialize(donesStr)
	if err != nil {
		return
	}

	doings, err:= p.doingsUnSerialize(doingsStr)
	if err != nil {
		return
	}

	slaves, err := p.slavesUnSerialize(slavesStr)
	if err != nil {
		return
	}

	tasks, err:= p.tasksUnSerialize(tasksStr)
	if err != nil {
		return
	}
	//fmt.Printf("[tasks:%#v]\n", tasks)


	err = p.UnDones(dones)
	if err != nil {
		return err
	}
	//fmt.Printf("[p.dones = %#v]\n", p.dones)

	err = p.UnDoings(doings)
	if err != nil {
		return err
	}

	err = p.UnSlaves(slaves)
	if err != nil {
		return err
	}
	fmt.Printf("[slaves:%#v]\n", slaves)

	err = p.UnTasks(tasks)
	if err != nil {
		return err
	}

	return
}

func (p *Master) scan() {
	go func() {
		for {
			for k, v := range p.doings {
				if time.Since(time.Unix(time.Unix(0, v.Created).Unix(), 0)) >= time.Minute {
					p.Push("repush", v)
					delete(p.doings, k)
				}
			}
			p.save()
			time.Sleep(time.Minute)
			fmt.Fprintf(os.Stderr, "hello\n")
		}
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
	slaves map[string]int64
	donesFile	Data
	doingsFile	Data
	tasksFile	Data
	slavesFile	Data
	locker sync.Mutex
}

var ErrNoTask = errors.New("no task")
