package slave

import (
	"regexp"
)

func (p *Url) extract(body string) (ret []string, err error) {
	pattern := `song\?id=[\d]+`
	reg := regexp.MustCompile(pattern)
	return reg.FindAllString(body, -1), err
}

func NewUrl() Url {
	return Url{}
}

type Url struct {
}
