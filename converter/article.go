package converter

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/go-shiori/go-readability"
	"github.com/pkg/errors"
	"github.com/utrack/sendreadable/htmllatex"
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

func NewArticle(ctx context.Context,
	art readability.Article,
	url string,
	dwn htmllatex.ImageDownloader,
) (Article, error) {

	conv := htmllatex.New(dwn, url)

	ret := Article{
		URL:           conv.DoPlain(url),
		Content:       "",
		Title:         conv.DoPlain(art.Title),
		Author:        conv.DoPlain(art.Byline),
		Excerpt:       conv.DoPlain(art.Excerpt),
		Source:        conv.DoPlain(art.SiteName),
		AvgTimeString: "",
	}
	ret.Author = strings.TrimSpace(ret.Author)

	tex, err := conv.Do(ctx, art.Content)
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
