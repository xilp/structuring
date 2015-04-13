package slave

import (
	//"fmt"
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
	//fmt.Println("%#v", site)
	task := site.site.NewTask(info)
	return task
}

func Domain(url string) string {
	// TODO
	return "music.163.com"
}

func (p *Sites) NewSite(domain string) (site *SiteInfo) {
	switch domain {
	case "music.163.com":
		return &SiteInfo{site: m1c.New()}
	}
	return
}

func NewSites() Sites {
	return Sites{}
}

type Sites map[string]*SiteInfo

type SiteInfo struct {
	site types.Site
	cookie string
}
