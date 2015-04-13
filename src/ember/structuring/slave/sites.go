package slave

import (
	"ember/structuring/types"
	m1c "ember/structuring/sites/music.163.com"
)

func (p *Sites) NewTask(info types.TaskInfo) types.Task {
	domain := Domain(info.Url)
	site, ok := (*p)[domain]
	if !ok {
		site = p.NewSite(domain)
		if site == nil {
			return nil
		}
		(*p)[domain] = site
	}
	task := site.NewTask(info)
	return task
}

//func domain(url string) string {
func Domain(url string) string {
	// TODO
	return "music.163.com"
}

func (p *Sites) NewSite(domain string) (site types.Site) {
	switch domain {
	case "music.163.com":
		site = m1c.New()
	}
	return
}

func NewSites() Sites {
	return Sites{}
}

type Sites map[string]types.Site
