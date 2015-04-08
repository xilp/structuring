package slave

import (
	m1c "ember/structuring/sites/music.163.com"
)

func NewSite(url string) Site {
	switch url: {
	case "music.163.com":
		return m1c.NewMusic163Com()
	}
	return nil
}


