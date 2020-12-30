package main

import (
	"bytes"
	"strings"

	"github.com/ryboe/q"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func htmlToTex(ht string) (string, error) {
	r := strings.NewReader(ht)

	n, err := html.ParseWithOptions(r, html.ParseOptionEnableScripting(false))
	if err != nil {
		return "", err
	}

	return walkNode(n), nil
}

func latexEscape(text string) string {
	for c, r := range latexSpecialSym {
		text = strings.ReplaceAll(text, c, r)
	}
	return text
}

func walkNode(n *html.Node) string {
	buffer := bytes.NewBuffer(nil)

	walker(n, buffer)
	return buffer.String()
}

func walker(n *html.Node, buf *bytes.Buffer) {

	switch n.DataAtom {
	case atom.Article, atom.Html, atom.Head, atom.Body:
	case atom.Span:
	case atom.Div:
	case atom.B:
		walkB(n, buf)
		return
	case atom.I:
		walkI(n, buf)
		return
	case atom.P:
		walkP(n, buf)
		return
	case atom.A:
		walkA(n, buf)
		return
	default:
		if n.DataAtom.String() != "" {
			q.Q("unknown type", n.DataAtom, n.DataAtom.String(), n.Attr)
		}
	}

	if n.Type == html.TextNode {
		buf.WriteString(latexEscape(n.Data))
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		walker(child, buf)
	}
}

func walkB(n *html.Node, buf *bytes.Buffer) {
	nbuf := bytes.NewBuffer(nil)

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		walker(child, nbuf)
	}
	str := nbuf.String()
	if str == "" {
		return
	}
	buf.WriteString("\\textbf{")
	buf.WriteString(str)
	buf.WriteByte('}')
}

func walkP(n *html.Node, buf *bytes.Buffer) {
	nbuf := bytes.NewBuffer(nil)

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		walker(child, nbuf)
	}
	str := nbuf.String()
	if str == "" {
		return
	}
	buf.WriteString("\n\\par ")
	buf.WriteString(str)
}

func walkI(n *html.Node, buf *bytes.Buffer) {
	nbuf := bytes.NewBuffer(nil)

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		walker(child, nbuf)
	}
	str := nbuf.String()
	if str == "" {
		return
	}
	buf.WriteString("\\emph{")
	buf.WriteString(str)
	buf.WriteByte('}')
}

func walkA(n *html.Node, buf *bytes.Buffer) {
	nbuf := bytes.NewBuffer(nil)
	q.Q("a", n.Attr)

	var url string
	for _, a := range n.Attr {
		if a.Key == "href" {
			url = a.Val
			break
		}
	}

	url = strings.ReplaceAll(url, `\`, `\\`)

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		walker(child, nbuf)
	}
	str := nbuf.String()
	buf.WriteString("\\href{" + url + "}{")
	buf.WriteString(str)
	buf.WriteByte('}')
}
