package slave

import (
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"ember/cli"
	"ember/http/rpc"
	"ember/structuring/types"
	"regexp"
	"net/http"
	"io/ioutil"
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
	slave.run(concurrent)
}

func (p *Slave) run(concurrent int) {
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
		//err = task.Run(p.append)
		err = task.Run(p.processTask)
		if err != nil {
			return err
		}
		fmt.Printf("done: %v\n", info)
		p.master.Done(p.id, info)
	}
}

func (p *Slave) processTask(info types.TaskInfo) (err error) {
	ret, err := p.Crawl(info.Url)
	if err != nil {
		return err
	}
	host := Domain(info.Url)
	for _, v := range ret {
		info.Url = "http://" + host + "/"+ v
		p.master.Push(p.id, info)
	}
	return err
}

func (p *Slave) Crawl(url string)(ret []string, err error) {
	body, err := p.fetchHtml(url)
	if err != nil {
		return nil, err
	}
	return p.extractUrl(body)
}

func (p *Slave) extractUrl(body string) (ret []string, err error) {
	pattern := `song\?id=[\d]+`
	reg := regexp.MustCompile(pattern)
	return reg.FindAllString(body, -1), err
}

func (p *Slave) getCookie(host string) (cookie string, err error) {
	for i := 0; i < 3; i++ {
		resp, err := http.Head(host)
		if err != nil {
			continue
		}
		cookie = resp.Header.Get("Set-Cookie")
		break
	}
	return cookie, err
}

func (p *Slave) fetchHtml(url string) (body string, err error) {
	domain := Domain(url)
	cookie := p.sites[domain].cookie
	if cookie == "" {
		cookie, err = p.getCookie("http://" + domain + "/")
		if err != nil {
			return "", err
		}
		p.sites[domain].cookie = cookie
	}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Cookie", cookie)
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:31.0) Gecko/20100101 Firefox/31.0")
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	return string(data), err
}

func (p *Slave) append(info types.TaskInfo) (err error) {
	fmt.Printf("appending: %v\n", info)
	return p.master.Push(p.id, info)
}

func NewSlave(addr string, id string) (p *Slave, err error) {
	var master Master
	err = rpc.NewClient(addr).Reg(&master, &MasterTrait{})
	if err != nil {
		return
	}
	p = &Slave{id, NewSites(), master}
	return
}

type Slave struct {
	id string
	sites Sites
	master Master
}

type Master struct {
	Done func(slave string, info types.TaskInfo) error
	Push func(slave string, info types.TaskInfo) error
	Pop func(slave string) (info types.TaskInfo, err error)
}

func (p *MasterTrait) Trait() map[string][]string {
	return map[string][]string {
		"Done": {"slave", "task"},
		"Push": {"slave", "task"},
		"Pop": {"slave"},
	}
}

type MasterTrait struct {
}

//var cookie = ""
