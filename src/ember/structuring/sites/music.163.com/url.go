package m1c

import (
	"regexp"
)

func (p *Url) extract(body []byte) (ret []string, err error) {
	pattern := `song\?id=[\d]+`
	reg := regexp.MustCompile(pattern)
	return reg.FindAllString(string(body), -1), err
}

func NewUrl() Url {
	return Url{}
}

type Url struct {
}
