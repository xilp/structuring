package slave

import (
	"ember/structuring/types"
	m1c "ember/structuring/sites/music.163.com"
	"os"
	"regexp"
	"strings"
)

func (p *Sites) NewTask(info types.TaskInfo) (types.Task) {
	domain := Domain(info.Url)
	site, ok := (*p).sites[domain]
	if !ok {
		site = p.NewSite(domain)
		if site == nil {
			return nil
		}
		(*p).sites[domain] = site
	}
	task := site.site.NewTask(info)
	return task
}

func Domain(url string) string {
	// TODO
	reg := regexp.MustCompile(`http://[^/:]+`)
	domain := reg.FindString(url)
	if domain  == "" {
		return "music.163.com"
	}
	domain  = strings.Replace(domain, "http://", "", -1)
	return domain
}

func (p *Sites) Register(domain string) (err error) {
	site, ok := (*p).sites[domain]
	if !ok {
		site = p.NewSite(domain)
		if site == nil {
			return nil
		}
		(*p).sites[domain] = site
	}
	return
}

func (p *Sites) NewSite(domain string) (site *SiteInfo) {
	switch domain {
	case "music.163.com":
		return &SiteInfo{site: m1c.New(p.root)}
	}
	return
}

func NewSites(root string) Sites {
	err := os.MkdirAll(root, 0755)
	if err != nil {
		println(err.Error())
	}
	return Sites{root:root, sites:make(map[string]*SiteInfo)}
}

type Sites struct {
	root string
	sites map[string]*SiteInfo
}

type SiteInfo struct {
	site types.Site
}
