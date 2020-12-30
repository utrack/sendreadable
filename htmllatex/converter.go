package htmllatex

import (
	"bytes"
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/ryboe/q"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

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
}

func escapeText(text string) string {
	text = strings.ReplaceAll(text, `\`, "/")
	for c, r := range latexSpecialSym {
		text = strings.ReplaceAll(text, c, r)
	}
	return text
}

type Converter struct {
	id ImageDownloader
}

type ImageDownloader interface {
	Download(context.Context, string) (string, error)
}

func New(dwn ImageDownloader) *Converter {
	return &Converter{
		id: dwn,
	}
}

func (c *Converter) Do(ctx context.Context, htext string) (string, error) {
	r := strings.NewReader(htext)

	n, err := html.ParseWithOptions(r, html.ParseOptionEnableScripting(false))
	if err != nil {
		return "", err
	}

	buffer := bytes.NewBuffer(nil)

	c.walker(ctx, n, buffer)
	return buffer.String(), nil
}

func (c *Converter) walker(ctx context.Context, n *html.Node, buf *bytes.Buffer) error {

	switch n.DataAtom {
	case atom.Article, atom.Html, atom.Head, atom.Body:
	case atom.Span:
	case atom.Div:
	case atom.Sup:
		return c.walkSup(ctx, n, buf)
	case atom.Center:
		return c.walkCenter(ctx, n, buf)
	case atom.B, atom.Strong:
		return c.walkB(ctx, n, buf)
	case atom.I, atom.Em:
		return c.walkI(ctx, n, buf)
	case atom.P:
		return c.walkP(ctx, n, buf)
	case atom.A:
		return c.walkA(ctx, n, buf)
	case atom.Img:
		return c.walkImg(ctx, n, buf)
	case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6:
		return c.walkHAny(ctx, n, buf)
	default:
		if n.DataAtom.String() != "" {
			q.Q("unknown type", n.DataAtom, n.DataAtom.String(), n.Attr)
		}
	}

	if n.Type == html.TextNode {
		buf.WriteString(escapeText(n.Data))
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		err := c.walker(ctx, child, buf)
		if err != nil {
			return err
		}
	}
	return nil
}

var hdgLevelToKeyword = map[atom.Atom]string{
	atom.H1: "chapter",
	atom.H2: "section",
	atom.H3: "subsection",
	atom.H4: "subsubsection",
	atom.H5: "paragraph",
	atom.H6: "subparagraph",
}

func (c *Converter) walkHAny(ctx context.Context, n *html.Node, buf *bytes.Buffer) error {
	nbuf := bytes.NewBuffer(nil)

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		err := c.walker(ctx, child, nbuf)
		if err != nil {
			return err
		}
	}
	str := nbuf.String()
	if str == "" {
		return nil
	}

	word, ok := hdgLevelToKeyword[n.DataAtom]
	if !ok {
		return errors.Errorf("cannot render heading level '%v'", n.DataAtom.String())
	}
	buf.WriteString("\n\\" + word + "*{")
	buf.WriteString(str)
	buf.WriteString("}\n")
	return nil

}

func (c *Converter) walkCenter(ctx context.Context, n *html.Node, buf *bytes.Buffer) error {
	nbuf := bytes.NewBuffer(nil)

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		err := c.walker(ctx, child, nbuf)
		if err != nil {
			return err
		}
	}
	str := nbuf.String()
	if str == "" {
		return nil
	}
	buf.WriteString("\n\\begin{center}\n")
	buf.WriteString(str)
	buf.WriteString("\n\\end{center}\n")
	return nil
}
func (c *Converter) walkSup(ctx context.Context, n *html.Node, buf *bytes.Buffer) error {
	nbuf := bytes.NewBuffer(nil)

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		err := c.walker(ctx, child, nbuf)
		if err != nil {
			return err
		}
	}
	str := nbuf.String()
	if str == "" {
		return nil
	}
	buf.WriteString("\\textsuperscript{")
	buf.WriteString(str)
	buf.WriteByte('}')
	return nil
}
func (c *Converter) walkB(ctx context.Context, n *html.Node, buf *bytes.Buffer) error {
	nbuf := bytes.NewBuffer(nil)

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		err := c.walker(ctx, child, nbuf)
		if err != nil {
			return err
		}
	}
	str := nbuf.String()
	if str == "" {
		return nil
	}
	buf.WriteString("\\textbf{")
	buf.WriteString(str)
	buf.WriteByte('}')
	return nil
}

func (c *Converter) walkP(ctx context.Context, n *html.Node, buf *bytes.Buffer) error {
	nbuf := bytes.NewBuffer(nil)

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		err := c.walker(ctx, child, nbuf)
		if err != nil {
			return err
		}
	}
	str := nbuf.String()
	if str == "" {
		return nil
	}
	buf.WriteString("\n\\par ")
	buf.WriteString(str)
	return nil
}

func (c *Converter) walkI(ctx context.Context, n *html.Node, buf *bytes.Buffer) error {
	nbuf := bytes.NewBuffer(nil)

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		err := c.walker(ctx, child, nbuf)
		if err != nil {
			return err
		}
	}
	str := nbuf.String()
	if str == "" {
		return nil
	}
	buf.WriteString("\\emph{")
	buf.WriteString(str)
	buf.WriteByte('}')
	return nil
}

func (c *Converter) walkA(ctx context.Context, n *html.Node, buf *bytes.Buffer) error {
	nbuf := bytes.NewBuffer(nil)

	var url string
	for _, a := range n.Attr {
		if a.Key == "href" {
			url = a.Val
			break
		}
	}

	url = strings.ReplaceAll(url, `\`, `\\`)

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		err := c.walker(ctx, child, nbuf)
		if err != nil {
			return err
		}
	}
	str := nbuf.String()
	if url == "#" {
		buf.WriteString(str)
		return nil
	}

	buf.WriteString("\\href{" + url + "}{")
	buf.WriteString(str)
	buf.WriteByte('}')
	return nil
}

func (c *Converter) walkImg(ctx context.Context, n *html.Node, buf *bytes.Buffer) error {
	var alt string
	var url string
	for _, a := range n.Attr {
		if a.Key == "alt" {
			alt = a.Val
		}
		if a.Key == "src" {
			url = a.Val
		}
	}

	if url == "" {
		return nil
	}
	m, err := c.id.Download(ctx, url)
	if err != nil {
		buf.WriteString("\n")
		buf.WriteString(`
\begin{center}
  image not loaded (` + err.Error() + `)
\end{center}
`,
		)
		return nil
	}
	buf.WriteByte('\n')
	buf.WriteString(`
\begin{figure}[h]
  \centering
  \includegraphics[max width=\textwidth,keepaspectratio]{` + m + `}
  {\caption*{` + escapeText(html.UnescapeString(alt)) + `}}
\end{figure}
`)
	return nil
}
