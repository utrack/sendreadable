package htmllatex

import (
	"bytes"
	"context"
	"net/url"
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
	id  ImageDownloader
	uri string
}

type ImageDownloader interface {
	Download(context.Context, string) (string, error)
}

func New(dwn ImageDownloader, uri string) *Converter {
	return &Converter{
		id:  dwn,
		uri: uri,
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
	case atom.Hr:
		return c.walkHr(ctx, n, buf)
	case atom.Dd:
		return c.walkDd(ctx, n, buf)
	case atom.Dt:
		return c.walkDt(ctx, n, buf)
	case atom.Dl:
		return c.walkDl(ctx, n, buf)
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

func (c *Converter) walkHr(ctx context.Context, n *html.Node, buf *bytes.Buffer) error {
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
	buf.WriteString("\n\\hrule\n")
	buf.WriteString(str)
	return nil
}

func (c *Converter) walkDl(ctx context.Context, n *html.Node, buf *bytes.Buffer) error {
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
	buf.WriteString("\n" + `\begin{description}[style=unboxed, labelwidth=\linewidth, font =\sffamily\itshape\bfseries, listparindent =0pt, before =\sffamily]`)
	buf.WriteString(str)
	buf.WriteString("\n\\end{description}\n")
	return nil
}
func (c *Converter) walkDt(ctx context.Context, n *html.Node, buf *bytes.Buffer) error {
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
	buf.WriteString("\n\\item[") // TODO escape bracket
	buf.WriteString(str)
	buf.WriteString("]")
	return nil
}
func (c *Converter) walkDd(ctx context.Context, n *html.Node, buf *bytes.Buffer) error {
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
	buf.WriteByte('\n')
	buf.WriteByte('\n')
	buf.WriteString(str)
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

	var uri string
	for _, a := range n.Attr {
		if a.Key == "href" {
			uri = a.Val
			break
		}
	}

	uri = strings.ReplaceAll(uri, `\`, `\\`)

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		err := c.walker(ctx, child, nbuf)
		if err != nil {
			return err
		}
	}
	str := nbuf.String()

	if u, err := url.Parse(uri); err == nil && !u.IsAbs() {
		r, err := url.Parse(c.uri)
		if err == nil {
			uri = r.ResolveReference(u).String()
		}
	}

	if uri == "#" {
		buf.WriteString(str)
		return nil
	}

	buf.WriteString("\\href{" + uri + "}{")
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
