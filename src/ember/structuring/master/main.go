package master

import (
	"flag"
	"math"
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
}

func (p *Master) Fetch(url string) error {
	return p.Push("master", types.NewTaskInfo(url, "index", math.MaxInt64))
}

func (p *Master) Done(slave string, task types.TaskInfo) (err error) {
	return
}

func (p *Master) Push(slave string, task types.TaskInfo) (err error) {
	p.save()
	return
}

func (p *Master) Pop(slave string) (task types.TaskInfo, err error) {
	return
}

func (p *Master) save() (err error) {
	return
}

func (p *Master) load() (err error) {
	return
}

func (p *Master) Trait() map[string][]string {
	st := slave.MasterTrait{}
	trait := st.Trait()
	trait["Fetch"] = []string{"url"}
	return trait
}

func NewMaster() *Master {
	p := &Master{}
	p.load()
	return p
}

type Master struct {
	tasks []types.TaskInfo
}
