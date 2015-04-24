package slave

import (
	"regexp"
	"strings"
	"ember/structuring/types"
)

func (p *Sites) Reg(domain string, site types.Site) {
	(*p)[domain] = site
}

func (p *Sites) NewTask(info types.TaskInfo) (task types.Task) {
	domain := p.domain(info.Url)
	site, ok := (*p)[domain]
	if !ok {
		return nil
	}
	return site.NewTask(info)
}

func (p *Sites) domain(url string) string {
	reg := regexp.MustCompile(`http://[^/:]+`)
	domain := reg.FindString(url)
	domain = strings.Replace(domain, "http://", "", -1)
	return domain
}

func NewSites() Sites {
	return make(map[string]types.Site)
}

type Sites map[string]types.Site
