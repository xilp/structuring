package m1c

import (
	"ember/structuring/types"
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

func (p *Site) ParseHtml(body []byte) (ret []string, err error) {
	return p.html.parse(body)
}

func (p *Site) ExtractUrl(body string) (ret []string, err error) {
	return p.url.extract(body)
}

func New() *Site {
	return &Site{NewUrl(), NewHtml(), NewCrawl()}
}

type Site struct {
	url Url
	html Html
	crawl Crawl
}
