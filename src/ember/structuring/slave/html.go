package slave

import (
	"strings"
	"regexp"
)

func (p *Html) parse(body string) (ret []string, err error) {

	pattern := `<meta name="keywords" content="([^，]+)`
	reg := regexp.MustCompile(pattern)
	ret = reg.FindAllString(body, -1)
	var songName = ""
	var idx = 0
	if ret  != nil {
		idx = strings.Index(ret[0], "content=\"")
		idx  = idx + len("content=\"")
		songName = ret[0][idx:]
	}

	pattern = `<meta name="description" content="([^。]+)`
	reg = regexp.MustCompile(pattern)
	ret = reg.FindAllString(body, -1)
	var singer = ""
	if ret != nil {
		idx = strings.Index(ret[0], "content=\"歌手：")
		idx = idx + len("content=\"歌手：")
		singer = ret[0][idx:]
	}

	pattern = `<meta name="description" content="([^>]+)>`
	reg = regexp.MustCompile(pattern)
	ret = reg.FindAllString(body, -1)
	midBody := ret[0]

	pattern = `所属专辑：([^。]+)`
	var album = ""
	reg = regexp.MustCompile(pattern)
	ret = reg.FindAllString(midBody, -1)
	if ret != nil {
		idx = strings.Index(ret[0], "所属专辑：")
		idx  = idx + len("所属专辑：")
		album = ret[0][idx:]
	}

	pattern = `发行时间：([^。]+)`
	reg = regexp.MustCompile(pattern)
	ret = reg.FindAllString(midBody, -1)
	var issueDate = ""
	if ret != nil {
		idx = strings.Index(ret[0], "发行时间：")
		idx  = idx + len("发行时间：")
		issueDate = ret[0][idx:]
	}

	pattern = `发行公司：([^。]+)`
	reg = regexp.MustCompile(pattern)
	ret = reg.FindAllString(midBody, -1)
	var	issueCompany = ""
	if ret != nil {
		idx = strings.Index(ret[0], "发行公司：")
		idx  = idx + len("发行公司：")
		issueCompany = ret[0][idx:]
	}

	pattern = `。([^。]+)。"`
	reg = regexp.MustCompile(pattern)
	ret = reg.FindAllString(midBody, -1)
	var note = ""
	if ret != nil {
		idx  = 3
		note = ret[0][idx:len(ret[0]) - 1]
		note = strings.Replace(note , "\n", "", -1)
	}

	pattern = `<div class="bd bd-open f-brk f-ib">([^\/]+)</div>`
	reg = regexp.MustCompile(pattern)
	ret = reg.FindAllString(body, -1)
	var songLyric = ""
	if ret != nil  {
		idx  = len(`<div class="bd bd-open f-brk f-ib">`)
		songLyric = ret[0][idx + 1:len(ret[0]) - 6]
		songLyric = strings.Replace(songLyric, `<div id="flag_more" class="f-hide">`, "", -1)
		songLyric = strings.Replace(songLyric, `<br>`, ",", -1)
		songLyric = strings.Replace(songLyric, "\n", "", -1)
	}

	ret = nil
	ret = append(ret, songName)
	ret = append(ret, singer)
	ret = append(ret, album)
	ret = append(ret, issueDate)
	ret = append(ret, issueCompany)
	ret = append(ret, note)
	ret = append(ret, songLyric)

	return ret, err
}

func NewHtml() Html{
	return Html{}
}

type Html struct {
}
