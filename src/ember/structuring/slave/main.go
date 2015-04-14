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
	"strings"
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
	//pv, err:= p.parseHtml(body)
	p.parseHtml(body)
	//fmt.Printf("[pv : %#v]\n", pv)
	return p.extractUrl(body)
}

func (p *Slave) extractUrl(body string) (ret []string, err error) {
	pattern := `song\?id=[\d]+`
	reg := regexp.MustCompile(pattern)
	return reg.FindAllString(body, -1), err
}

func (p *Slave) parseHtml(body string) (ret []string, err error) {
	keywords := `<meta name="keywords" content="为爱痴狂，收获 新歌+精选，刘若英，奶茶，Ren'e Liu，2001-06-01，滚石唱片" />`
	description := `<meta name="description" content="歌手：刘若英。所属专辑：收获 新歌+精选。发行时间：2001-06-01。发行公司：滚石唱片。 在刘若英这张精选辑中，除了新歌、过去的歌以外，滚石唱片还特别重新Mastering，加进了刘若英这些年来的工作、生活的声音纪实，让整张精选辑听来更有新意，并且更贴近奶茶的工作与生活。虽然可能只是一个机场的环境音、拍戏片场的现场音，或者是奶茶和制作人对Key的录音室记录，然而这些都是刘若英，工作中的刘若英，生活的刘若英，我们因而听见她的歌曲以外，另外一个有趣、生动的刘若英。而在正版专辑的POWER CD中，也特别收录了刘若英98年纯净影像代表作“很爱很爱你”的Music Video，以及历年来奶茶唱片专辑电视CF广告的精华版，以回馈给一直以来支持奶茶的歌迷朋友。。" />`
	//description := `<meta name="description" content="歌手：EXO。所属专辑：EXODUS。发行时间：2015-03-30。发行公司：(주)KT뮤직。EXO正规2辑《EXODUS》于3月30日发行，专辑共收录中韩两版各10首风格多样的歌曲，>`
	_ = keywords
	_ = description
	pattern := `<meta name="keywords" content="([^，]+)`
	reg := regexp.MustCompile(pattern)
	//return reg.FindAllString(body, -1), err
	ret = reg.FindAllString(body, -1)
	idx := strings.Index(ret[0], "content=\"")
	idx  = idx + len("content=\"")
	//println(idx)
	//println(ret[0])
	songName := ret[0][idx:]

	pattern = `<meta name="description" content="([^。]+)`
	reg = regexp.MustCompile(pattern)
	//return reg.FindAllString(body, -1), err
	ret = reg.FindAllString(body, -1)
	idx = strings.Index(ret[0], "content=\"歌手：")
	idx  = idx + len("content=\"歌手：")
	//println(ret[0])
	singer := ret[0][idx:]

	pattern = `<meta name="description" content="([^/>]+)`
	reg = regexp.MustCompile(pattern)
	//return reg.FindAllString(body, -1), err
	ret = reg.FindAllString(body, -1)
	midBody := ret[0]
	//println(midBody)

	pattern = `所属专辑：([^。]+)`
	reg = regexp.MustCompile(pattern)
	//return reg.FindAllString(body, -1), err
	ret = reg.FindAllString(midBody, -1)
	idx = strings.Index(ret[0], "所属专辑：")
	idx  = idx + len("所属专辑：")
	//println(ret[0])
	album := ret[0][idx:]

	pattern = `发行时间：([^。]+)`
	reg = regexp.MustCompile(pattern)
	//return reg.FindAllString(body, -1), err
	ret = reg.FindAllString(midBody, -1)
	idx = strings.Index(ret[0], "发行时间：")
	idx  = idx + len("发行时间：")
	//println(ret[0])
	issueDate := ret[0][idx:]

	pattern = `发行公司：([^。]+)`
	reg = regexp.MustCompile(pattern)
	//return reg.FindAllString(body, -1), err
	ret = reg.FindAllString(midBody, -1)
	idx = strings.Index(ret[0], "发行公司：")
	idx  = idx + len("发行公司：")
	//println(ret[0])
	issueCompany := ret[0][idx:]

	idx = strings.Index(midBody, "发行公司：")
	idx = idx + len(ret[0])  + 3
	note := midBody[idx:len(midBody) - 4]


	pattern = `<div class="bd bd-open f-brk f-ib">([^\/]+)`
	reg = regexp.MustCompile(pattern)
	//return reg.FindAllString(body, -1), err
	ret = reg.FindAllString(body, -1)
	//idx = strings.Index(ret[0], "发行时间：")
	idx  = len(`<div class="bd bd-open f-brk f-ib">`)
	//println(ret[0])
	songLyric := ret[0][idx + 1:len(ret[0]) - 2]
	//songLyric := ret[0]
	songLyric = strings.Replace(songLyric, `<div id="flag_more" class="f-hide">`, "", -1)
	songLyric = strings.Replace(songLyric, `<br>`, ",", -1)



	println("songName",songName)
	println("singer:",singer)
	println("album:",album)
	println("issueDate:",issueDate)
	println("issueCompany:",issueCompany)
	println("note:", note)
	println("songLyric :", songLyric)

	//ret = res

	return ret, err
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
