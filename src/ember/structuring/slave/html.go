package slave

import (
	"fmt"
	"strings"
	"regexp"
)

func (p *Html) parse(body string) (ret []string, err error) {
	keywords := `<meta name="keywords" content="为爱痴狂，收获 新歌+精选，刘若英，奶茶，Ren'e Liu，2001-06-01，滚石唱片" />`
	description := `<meta name="description" content="歌手：刘若英。所属专辑：收获 新歌+精选。
		发行时间：2001-06-01。发行公司：滚石唱片。 在刘若英这张精选辑中，除了新歌、过去的歌以外，
		滚石唱片还特别重新Mastering，加进了刘若英这些年来的工作、生活的声音纪实，
		让整张精选辑听来更有新意，并且更贴近奶茶的工作与生活。虽然可能只是一个机场的环境音、
		拍戏片场的现场音，或者是奶茶和制作人对Key的录音室记录，然而这些都是刘若英，工作中的刘若英，
		生活的刘若英，我们因而听见她的歌曲以外，另外一个有趣、生动的刘若英。而在正版专辑的POWER CD中，
		也特别收录了刘若英98年纯净影像代表作“很爱很爱你”的Music Video，
		以及历年来奶茶唱片专辑电视CF广告的精华版，以回馈给一直以来支持奶茶的歌迷朋友。。" />`
	_ = keywords
	_ = description


	pattern := `<meta name="keywords" content="([^，]+)`
	reg := regexp.MustCompile(pattern)
	ret = reg.FindAllString(body, -1)
	idx := strings.Index(ret[0], "content=\"")
	idx  = idx + len("content=\"")
	songName := ret[0][idx:]

	pattern = `<meta name="description" content="([^。]+)`
	reg = regexp.MustCompile(pattern)
	ret = reg.FindAllString(body, -1)
	idx = strings.Index(ret[0], "content=\"歌手：")
	idx  = idx + len("content=\"歌手：")
	singer := ret[0][idx:]

	//body = `<meta name="description" content="歌手：AaA。所属专辑：ソウルエッジボーイ/キモノジェットガール。发行时间：2006-07-19。发行公司：エイベックス.トラックス。网易云音乐是一款专注于发现与分享的音乐产品，依托专业音乐人、DJ、好友推荐及社交功能，为用户打造全新的音乐生活。" />`
	//body = `<meta name="description" content="歌手：何耀珊。所属专辑：Without Love。发行时间：2004-10-06。发行公司：Virgin Records。网易云音乐是一款专注于发现与分享的音乐产品，依托专业音乐人、DJ、好友推荐及社交功能，为用户打造全新的音乐生活。" />`
	body = `<meta name="description" content="歌手：何耀珊。所属专辑：China Wine。发行时间：2007-07-24。发行公司：Carnival House Records Inc./JH Music, Inc。2007年何耀珊与美国歌手Wyclef Jean, Elephant Man &amp; Tony Matteron合作的单曲，MV里面Sun颠覆以往给人的印象，大跳艳舞，造型抢眼，媒体粉丝火热关注...并且热浪持续蔓延至亚洲。" />`

	//pattern = `<meta name="description" content="([^/>]+)`
	pattern = `<meta name="description" content="([^>]+)>`
	reg = regexp.MustCompile(pattern)
	ret = reg.FindAllString(body, -1)
	midBody := ret[0]
	//println(midBody)

	pattern = `所属专辑：([^。]+)`
	reg = regexp.MustCompile(pattern)
	ret = reg.FindAllString(midBody, -1)
	idx = strings.Index(ret[0], "所属专辑：")
	idx  = idx + len("所属专辑：")
	album := ret[0][idx:]

	pattern = `发行时间：([^。]+)`
	reg = regexp.MustCompile(pattern)
	ret = reg.FindAllString(midBody, -1)
	if ret == nil {
		return ret, err
	}
	idx = strings.Index(ret[0], "发行时间：")
	idx  = idx + len("发行时间：")
	issueDate := ret[0][idx:]

	pattern = `发行公司：([^。]+)`
	reg = regexp.MustCompile(pattern)
	ret = reg.FindAllString(midBody, -1)
	var	issueCompany = ""
	var note = ""
	if ret != nil {
		idx = strings.Index(ret[0], "发行公司：")
		idx  = idx + len("发行公司：")
		issueCompany = ret[0][idx:]

		idx = strings.Index(midBody, "发行公司：")
		println(midBody)
		println(idx)
		println("len(midBody):", len(midBody))
		if idx == -1 {
			return ret, err
		}
		idx = idx + len(ret[0])  + 3
		note = midBody[idx:len(midBody) - 4]
		note = strings.Replace(note , "\n", "", -1)
	} else {
		return ret, err
		/*
		idx = strings.Index(midBody, "发行时间：")
		idx = idx + len("发行时间：") + len(issueDate)  + 3
		note = midBody[idx:len(midBody) - 4]
		note = strings.Replace(note , "\n", "", -1)
		*/
	}


	pattern = `<div class="bd bd-open f-brk f-ib">([^\/]+)</div>`
	reg = regexp.MustCompile(pattern)
	ret = reg.FindAllString(body, -1)
	fmt.Printf("[ret:%#v]\n")
	var songLyric = ""
	//if ret != nil  {
	idx  = len(`<div class="bd bd-open f-brk f-ib">`)
	println("idx:", idx)
	println("ret[0]:", ret[0])
	songLyric = ret[0][idx + 1:len(ret[0]) - 2]
	songLyric = strings.Replace(songLyric, `<div id="flag_more" class="f-hide">`, "", -1)
	songLyric = strings.Replace(songLyric, `<br>`, ",", -1)
	songLyric = strings.Replace(songLyric, "\n", "", -1)
	//}

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
