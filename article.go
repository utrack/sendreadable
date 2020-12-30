package main

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/go-shiori/go-readability"
	"github.com/k3a/html2text"
	"github.com/pkg/errors"
)

type Article struct {
	URL           string
	Content       string
	Title         string
	Author        string
	Excerpt       string
	Source        string
	Image         string
	AvgTimeString string
}

func NewArticle(art readability.Article, url string) (Article, error) {

	ret := Article{
		URL:           url,
		Content:       "",
		Title:         art.Title,
		Author:        art.Byline,
		Excerpt:       art.Excerpt,
		Source:        art.SiteName,
		Image:         art.Image,
		AvgTimeString: "",
	}
	ret.Author = strings.TrimSpace(ret.Author)

	tex, err := htmlToTex(art.Content)
	if err != nil {
		return Article{}, errors.Wrap(err, "can't convert HTML to Tex")
	}
	ret.Content = tex

	count := wc(art.TextContent)

	avgTime := float32(count) / 140

	t := time.Duration(avgTime) * time.Minute
	tMins := int(t.Minutes())
	humanStr := "minutes"
	if tMins == 1 {
		humanStr = "minute"
	}
	ret.AvgTimeString = fmt.Sprintf("%v %v", tMins, humanStr)

	return ret, nil
}

func wc(s string) int {
	var count int

	for _, c := range []rune(s) {
		if unicode.IsSpace(c) {
			count++
		}
	}
	return count + 1
}

var latexSpecialSym = map[string]string{
	`&`: `\&`,
	`%`: `\%`,
	`$`: `\$`,
	`#`: `\#`,
	`_`: `\_`,
	`{`: `\{`,
	`}`: `\}`,
	`~`: `\textasciitilde`,
	`^`: `\textasciicircum`,
	`\`: `/`,
}

func formattedText(ht string) string {
	plain := html2text.HTML2Text(ht)
	for c, r := range latexSpecialSym {
		plain = strings.ReplaceAll(plain, c, r)
	}
	return plain
}
