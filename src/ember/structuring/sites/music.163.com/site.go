package m1c

import (
	"ember/structuring/types"
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
	task := types.NewTaskInfo(p.url, "song", 0)
	return appender(task)
}

type Song struct {
	url string
}

func (p *Site) NewTask(info types.TaskInfo) types.Task {
	switch info.Type {
	}
	return &Song{info.Url}
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

func New() *Site {
	return &Site{"music.163.com", NewUrl(), NewHtml(), NewCrawl()}
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
	url Url
	html Html
	crawl Crawl
}
