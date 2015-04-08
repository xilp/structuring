package slave

import (
	"ember/structuring/types"
	m1c "ember/structuring/sites/music.163.com"
)

// TODO: rpc

func (p *Sites) NewTask(task types.TaskInfo) types.Task {
	domain := domain(task.Url)
	site, ok := (*p)[domain]
	if !ok {
		site := p.NewSite(domain)
		(*p)[domain] = site
	}
	return site.NewTask(task)
}

func domain(url string) string {
	// TODO
	return url
}

func (p *Sites) NewSite(url string) types.Site {
	switch url {
	case "music.163.com":
		return m1c.New()
	}
	return nil
}

func NewSites() Sites {
	return Sites{}
}

type Sites map[string]types.Site
