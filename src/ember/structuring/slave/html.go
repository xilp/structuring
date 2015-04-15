package slave

import (
	"bytes"
	"regexp"
)

func (p *Html) splitHtml(content, key []byte, pattern string ) (word []byte) {
	var idx = 0
	reg := regexp.MustCompile(pattern)
	b := reg.Find(content)
	if b != nil {
		idx = bytes.Index(b, key)
		idx  = idx + len(key)
		word = b[idx:]
	} else {
		word = nil
	}
	return word
}

func (p *Html) parse(body []byte) (ret []string, err error) {
	var songName, singer, album, issueDate, issueCompany, note, songLyric []byte
	var word []byte
	var idx = 0
	var b []byte

	pattern := `<meta name="keywords" content="([^，]+)`
	key := []byte("content=\"")
	word = p.splitHtml(body, key, pattern)
	if word != nil {
		songName = word
	}

	pattern = `<meta name="description" content="([^。]+)`
	key = []byte("content=\"歌手：")
	word = p.splitHtml(body, key, pattern)
	if word != nil {
		singer = word
	}

	pattern = `<meta name="description" content="([^>]+)>`
	reg := regexp.MustCompile(pattern)
	b = reg.Find(body)
	midBody := b

	pattern = `所属专辑：([^。]+)`
	key = []byte("所属专辑：")
	word = p.splitHtml(body, key, pattern)
	if word != nil {
		album = word
	}

	pattern = `发行时间：([^。]+)`
	key = []byte("发行时间：")
	word = p.splitHtml(body, key, pattern)
	if word != nil {
		issueDate = word
	}

	pattern = `发行公司：([^。]+)`
	key = []byte("发行公司：")
	word = p.splitHtml(body, key, pattern)
	if word != nil {
		issueCompany = word
	}

	pattern = `。([^。]+)。"`
	reg = regexp.MustCompile(pattern)
	b = reg.Find(midBody)
	if b != nil {
		idx  = 3
		note = b[idx:len(b) - 1]
		note = bytes.Replace(note , []byte("\n"), []byte(""), -1)
	}

	pattern = `<div class="bd bd-open f-brk f-ib">([^\/]+)</div>`
	reg = regexp.MustCompile(pattern)
	b = reg.Find(body)
	if b != nil {
		idx  = len(`<div class="bd bd-open f-brk f-ib">`)
		songLyric = b[idx + 1:len(b) - 6]
		songLyric = bytes.Replace(songLyric, []byte(`<div id="flag_more" class="f-hide">`), []byte(""), -1)
		songLyric = bytes.Replace(songLyric, []byte(`<br>`), []byte(","), -1)
		songLyric = bytes.Replace(songLyric, []byte("\n"), []byte(""), -1)
	}

	ret = append(ret, string(songName))
	ret = append(ret, string(singer))
	ret = append(ret, string(album))
	ret = append(ret, string(issueDate))
	ret = append(ret, string(issueCompany))
	ret = append(ret, string(note))
	ret = append(ret, string(songLyric))

	return ret, err
}

func NewHtml() Html{
	return Html{}
}

type Html struct {
}
