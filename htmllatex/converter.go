package htmllatex

import (
	"bytes"
	"context"
	"net/url"
	"strconv"
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
	text = strings.ReplaceAll(text, `[`, "{[")
	text = strings.ReplaceAll(text, `]`, "]}")
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

func (c *Converter) DoPlain(text string) string {
	return escapeText(text)
}

type Response struct {
	Content string
}

func (c *Converter) Do(ctx context.Context, htext string) (*Response, error) {
	q.Q(htext)
	r := strings.NewReader(htext)

	n, err := html.ParseWithOptions(r, html.ParseOptionEnableScripting(false))
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(nil)

	c.walker(ctx, n, buffer)
	ret := &Response{
		Content: buffer.String(),
	}
	return ret, nil
}

func (c *Converter) walker(ctx context.Context, n *html.Node, buf *bytes.Buffer) error {
	//buf.WriteString(" \\iffalse " + n.DataAtom.String() + " \\fi ")

	switch n.DataAtom {
	case atom.Article, atom.Html, atom.Head, atom.Body:
	case atom.Span:
	case atom.Div:
	case atom.Table:
		return c.walkTable(ctx, n, buf)
	case atom.Tr:
		return c.walkTr(ctx, n, buf)
	case atom.Td, atom.Th: // TODO bold center for Th
		return c.walkTd(ctx, n, buf)
	case atom.Br:
		return c.walkBr(ctx, n, buf)
	case atom.Hr:
		return c.walkHr(ctx, n, buf)
	case atom.Ol:
		return c.walkOl(ctx, n, buf)
	case atom.Ul:
		return c.walkUl(ctx, n, buf)
	case atom.Li:
		return c.walkLi(ctx, n, buf)
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
	case atom.Blockquote:
		return c.walkBlockquote(ctx, n, buf)
	case atom.B, atom.Strong:
		return c.walkB(ctx, n, buf)
	case atom.Abbr:
		return c.walkAbbr(ctx, n, buf)
	case atom.Small:
		return c.walkSmall(ctx, n, buf)
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
		str := escapeText(n.Data)
		strCh := strings.TrimSpace(str)
		strCh = strings.Trim(str, "\r\n\t")
		if strCh != "" {
			buf.WriteString(str)
		}
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

func (c *Converter) walkTable(ctx context.Context, n *html.Node, buf *bytes.Buffer) error {
	nbuf := bytes.NewBuffer(nil)

	cols := countTableCells(n)
	ctx = cntNotOuterPar(ctx)

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
	buf.WriteString("\n\\begin{center}\\begin{tabularx}{\\textwidth}{@{}l|")
	for i := 1; i < cols-1; i++ {
		buf.WriteString("X|")
	}
	buf.WriteString("X@{}")
	buf.WriteString("}\\toprule\n")
	buf.WriteString(str)
	buf.WriteString(" \\bottomrule\\end{tabularx}\\end{center}")

	return nil
}

func (c *Converter) walkTr(ctx context.Context, n *html.Node, buf *bytes.Buffer) error {
	nbuf := bytes.NewBuffer(nil)
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		err := c.walker(ctx, child, nbuf)
		if err != nil {
			return err
		}
	}
	str := nbuf.String()
	str = strings.TrimSuffix(str, " &")
	buf.WriteString(str)
	buf.WriteString(" \\\\\n")
	return nil
}
func (c *Converter) walkTd(ctx context.Context, n *html.Node, buf *bytes.Buffer) error {
	nbuf := bytes.NewBuffer(nil)
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		err := c.walker(ctx, child, nbuf)
		if err != nil {
			return err
		}
	}
	var colspan int
	for _, v := range n.Attr {
		if v.Key != "colspan" {
			continue
		}
		colspan, _ = strconv.Atoi(v.Val)
	}
	str := nbuf.String()
	//str = strings.ReplaceAll(str, "\n", "\\linebreak")
	if colspan > 0 {
		buf.WriteString("\\multicolumn{" + strconv.Itoa(colspan) + "}{c}{")
	}
	buf.WriteString(str)
	if colspan > 0 {
		buf.WriteByte('}')
	}
	buf.WriteString(" &")
	return nil
}

func countTableCells(n *html.Node) int {
	if n.DataAtom == atom.Tr {
		var localCount int
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			if child.DataAtom == atom.Td || child.DataAtom == atom.Th {
				localCount++
				continue
			}
		}
		return localCount
	}

	var max int
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		c := countTableCells(child)
		if c > max {
			max = c
		}
	}

	return max
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
func (c *Converter) walkOl(ctx context.Context, n *html.Node, buf *bytes.Buffer) error {
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
	buf.WriteString("\n\\begin{enumerate}")
	buf.WriteString(str)
	buf.WriteString("\\end{enumerate}")
	return nil
}
func (c *Converter) walkUl(ctx context.Context, n *html.Node, buf *bytes.Buffer) error {
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
	if strings.Trim(str, "\n\r\t ") == "" {
		return nil
	}
	buf.WriteString("\n\\begin{itemize}")
	buf.WriteString(str)
	buf.WriteString("\n\\end{itemize}")
	return nil
}
func (c *Converter) walkLi(ctx context.Context, n *html.Node, buf *bytes.Buffer) error {
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
	buf.WriteString("\n\\item ")
	buf.WriteString(str)
	return nil
}

func (c *Converter) walkAbbr(ctx context.Context, n *html.Node, buf *bytes.Buffer) error {
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
	var title string
	for _, a := range n.Attr {
		if a.Key == "title" {
			title = a.Val
			break
		}
	}
	buf.WriteString(str)
	buf.WriteString("\\footnote{" + escapeText(title) + "}")

	return nil
}
func (c *Converter) walkBr(ctx context.Context, n *html.Node, buf *bytes.Buffer) error {
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
	buf.WriteString("\\linebreak")
	buf.WriteString(str)
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

func (c *Converter) walkBlockquote(ctx context.Context, n *html.Node, buf *bytes.Buffer) error {
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
	buf.WriteString("\n\n\\begin{quotationb}\n")
	buf.WriteString(str)
	buf.WriteString("\n\\end{quotationb}\n")
	return nil
}
func (c *Converter) walkSmall(ctx context.Context, n *html.Node, buf *bytes.Buffer) error {
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
	buf.WriteString("\\small{")
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
	if isParable(ctx) {
		buf.WriteString("\n\\par ")
	}
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

	// uri = strings.ReplaceAll(uri, `\`, `\\`)
	// uri = strings.ReplaceAll(uri, `%`, `%%`)
	ctx = cntNotParable(ctx)

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

	buf.WriteString("\\href{" + escapeText(uri) + "}{")
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

	alt = escapeText(html.UnescapeString(alt))

	if url == "" {
		return nil
	}
	m, err := c.id.Download(ctx, url)
	if err != nil {
		buf.WriteString(`(bad img: ` + err.Error() + `))`)
		return nil
	}
	if isOuterPar(ctx) {
		buf.WriteString(`
\begin{figure}[h]
  \centering
  \includegraphics[max width=\textwidth,keepaspectratio]{` + m + `}
  {\caption*{` + alt + `}}
\end{figure}
`)
		return nil
	}

	buf.WriteString("\\includegraphics[max width=\\textwidth,keepaspectratio]{")
	buf.WriteString(m)
	buf.WriteByte('}')
	if alt != "" {
		buf.WriteString("\\footnote{" + alt + "}")
	}

	return nil
}
