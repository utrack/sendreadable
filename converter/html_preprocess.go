package converter

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/ryboe/q"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type htmlData struct {
	Language string
	Body     []byte
}

func preprocessHtml(b []byte) (*htmlData, error) {
	n, err := html.Parse(bytes.NewReader(b))
	if err != nil {
		return nil, errors.Wrap(err, "parsing HTML")
	}

	lang := findHtmlLang(n)
	return &htmlData{Language: lang, Body: b}, nil
}

func findHtmlLang(n *html.Node) string {

	if n.DataAtom == atom.Html {
		q.Q(n.Attr)
		var lang string
		for _, a := range n.Attr {
			if a.Key == "lang" {
				lang = a.Val
			}
		}
		return lang
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if lang := findHtmlLang(child); lang != "" {
			return lang
		}
	}
	return ""
}
