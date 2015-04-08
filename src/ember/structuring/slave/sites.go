package slave

import (
	"ember/structuring/types"
	m1c "ember/structuring/sites/music.163.com"
)

func (p *Sites) NewTask(task types.TaskInfo) types.Task {
	domain := domain(task.Url)
	site, ok := (*p)[domain]
	if !ok {
		site = p.NewSite(domain)
		if site == nil {
			return nil
		}
		(*p)[domain] = site
	}
	return site.NewTask(task)
}

func domain(url string) string {
	// TODO
	return "music.163.com"
}

func (p *Sites) NewSite(domain string) types.Site {
	switch domain {
	case "music.163.com":
		return m1c.New()
	}
	return nil
}

func NewSites() Sites {
	return Sites{}
}

type Sites map[string]types.Site
