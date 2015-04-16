package m1c

import (
	"ember/structuring/types"
	"fmt"
	"net/http"
)

func (p *Song) Run(appender types.Appender) (err error) {
	// TODO
	// extrat infos
	// get similars
	// new song task -> chan
	if len(p.url) > 100 {
		return
	}
	//task := types.NewTaskInfo(p.url + "*", "song", 0)

	ret, err := p.Crawl(p.url)
	if err != nil {
		return err
	}

	task := types.NewTaskInfo(p.url, "song", 0)
	for _, v := range ret {
		task.Url = "http://" + p.site.domain + "/" + v
		err = appender(task)
	}

	return err
}

func (p *Song) Crawl(url string) (ret []string, err error) {
	body, err := p.site.FetchHtml(url)
	if err != nil {
		return nil, err
	}
	pv, err := p.site.ParseHtml(body)
	if pv == nil || err != nil {
		return nil, err
	}
	p.site.Write(url, pv)
	return p.site.ExtractUrl(body)
}

type Song struct {
	site *Site
	url string
}

func (p *Site) NewTask(info types.TaskInfo) types.Task {
	switch info.Type {
	}
	return &Song{p, info.Url}
}

func (p *Site) FetchHtml(url string) (ret []byte, err error) {
	cookie, err := p.GetCookie()
	// TODO check err 
	p.html.cookie = cookie
	return p.html.fetch(url)
}

func (p *Site) ParseHtml(body []byte) (ret []string, err error) {
	return p.html.parse(body)
}

func (p *Site) ExtractUrl(body []byte) (ret []string, err error) {
	return p.url.extract(body)
}

func (p *Site) Write(url string, ret []string) (err error) {
	str := p.version + "\t" + url
	for _, v := range ret {
		str = str + "\t" + v
	}
	str = str + "\n"
	return p.data.write(str, 0)
}

func (p *Site) Serialize() (ret []byte, err error) {
	fmt.Printf("[]")
	return ret, err
}

func New() *Site {
	return &Site{"music.163.com", "01", NewUrl(), NewHtml(), NewCrawl(), NewData()}
}

func (p *Site) GetCookie() (cookie string, err error) {
	for i := 0; i < 3; i++ {
		resp, err := http.Head(p.domain)
		if err != nil {
			continue
		}
		cookie = resp.Header.Get("Set-Cookie")
		break
	}
	return cookie, err
}

type Site struct {
	domain string
	version string
	url Url
	html Html
	crawl Crawl
	data Data
}
