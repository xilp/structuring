package m1c

import (
	"bytes"
	"fmt"
	"regexp"
)

func (p *Html) splitHtml(content, key []byte, pattern string, word *[]byte) () {
	var idx = 0
	reg := regexp.MustCompile(pattern)
	b := reg.Find(content)
	if b != nil {
		idx = bytes.Index(b, key)
		idx  = idx + len(key)
		*word = b[idx:]
	}
}

func (p *Html) parse(body []byte) (ret []string, err error) {
	var songName, singer, album, issueDate, issueCompany, note, songLyric []byte
	var idx = 0
	var b []byte

	p.splitHtml(body, []byte("content=\""), `<meta name="keywords" content="([^，]+)`, &songName)
	p.splitHtml(body, []byte("content=\"歌手："), `<meta name="description" content="([^。]+)`, &singer)

	reg := regexp.MustCompile(`<meta name="description" content="([^>]+)>`)
	midBody := reg.Find(body)
	if midBody == nil {
		return nil, err
	}

	p.splitHtml(midBody, []byte("所属专辑："), `所属专辑：([^。]+)`, &album)
	p.splitHtml(midBody, []byte("发行时间："), `发行时间：([^。]+)`, &issueDate)
	p.splitHtml(midBody, []byte("发行公司："), `发行公司：([^。]+)`, &issueCompany)

	reg = regexp.MustCompile(`。([^。]+)。"`)
	b = reg.Find(midBody)
	if b != nil {
		idx  = 3
		note = b[idx:len(b) - 1]
		note = bytes.Replace(note , []byte("\n"), []byte(""), -1)
	}

	reg = regexp.MustCompile(`<div class="bd bd-open f-brk f-ib">([^\/]+)</div>`)
	b = reg.Find(body)
	if b != nil {
		idx  = len(`<div class="bd bd-open f-brk f-ib">`)
		songLyric = b[idx + 1:len(b) - 6]
		songLyric = bytes.Replace(songLyric, []byte(`<div id="flag_more" class="f-hide">`), []byte(""), -1)
		songLyric = bytes.Replace(songLyric, []byte(`<br>`), []byte(","), -1)
		songLyric = bytes.Replace(songLyric, []byte("\n"), []byte(""), -1)
	}

	ret = append(ret, string(songName), string(singer), string(album))
	ret = append(ret, string(issueDate), string(issueCompany), string(note), string(songLyric))
	fmt.Printf("[ret:%#v]\n", ret)

	return ret, err
}

func NewHtml() Html{
	return Html{}
}

type Html struct {
}
