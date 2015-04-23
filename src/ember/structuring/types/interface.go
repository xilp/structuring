package types

type Slave struct {
	Search func(key string)(ret [][]string, err error) `args:"key" return:"ret"`
}

type Master struct {
	Fetch func(url string) error `args:"url"`
	Register func(addr, slave string) error `args:"addr,slave"`
	Done func(slave string, info TaskInfo) error `args:"slave,info"`
	Push func(slave string, info TaskInfo) error `args:"slave,info"`
	Pop func(slave string) (info TaskInfo, err error) `args:"slave" return:"info"`
	Dones func() (urls []string, err error) `return:"urls"`
	Doings func() (tasks []TaskInfo, err error) `return:"tasks"`
	Tasks func() (tasks []TaskInfo, err error) `return:"tasks"`
}

