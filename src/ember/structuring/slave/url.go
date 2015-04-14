package slave

import (
	"regexp"
)

func (p *Slave) extractUrl(body string) (ret []string, err error) {
	pattern := `song\?id=[\d]+`
	reg := regexp.MustCompile(pattern)
	return reg.FindAllString(body, -1), err
}

