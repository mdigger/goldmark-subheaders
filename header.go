// Package subheaders is a extension for the goldmark
// (http://github.com/yuin/goldmark).
//
// This extension adds support for subheaders in markdown.
//  # Header
//  # Sub-header text
package subheaders

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// HeaderBlock is a HTML5 Header block.
type HeaderBlock struct {
	ast.BaseBlock
}

// KindHeaderBlock is a NodeKind of the TextBlock node.
var KindHeaderBlock = ast.NewNodeKind("HeaderBlock")

// Kind implements Node.Kind.
func (n *HeaderBlock) Kind() ast.NodeKind {
	return KindHeaderBlock
}

// Dump implements Node.Dump.
func (n *HeaderBlock) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

// RegisterFuncs implement renderer.NodeRenderer interface.
func (n *HeaderBlock) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindHeaderBlock, n.renderHeader)
}

func (n *HeaderBlock) renderHeader(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// TODO:
	if entering {
		w.WriteString("<header>\n")
	} else {
		w.WriteString("</header>\n")
	}
	return ast.WalkContinue, nil
}

// Transform combines consecutive headings of one level in one header block.
func (n *HeaderBlock) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
	var headers = make([]*ast.Heading, 0, 100)
	// walk all elements
	walk(doc, func(node ast.Node) (status ast.WalkStatus, err error) {
		status = ast.WalkContinue
		if node.Kind() != ast.KindHeading {
			return // it's not header
		}
		next := node.NextSibling()
		if next == nil || next.Kind() != ast.KindHeading || next.HasBlankPreviousLines() {
			return // next node not header or with blank line before
		}
		header := node.(*ast.Heading)
		if header.Level != next.(*ast.Heading).Level {
			return // next node header with different level
		}
		prev := node.PreviousSibling()
		if prev != nil && prev.Kind() == ast.KindHeading &&
			prev.(*ast.Heading).Level == header.Level &&
			!header.HasBlankPreviousLines() {
			return // previous node also header with same level
		}
		// add header to headers list for transformation
		headers = append(headers, header)
		return
	})
	// replace headers
	for _, h := range headers {
		next := h.NextSibling()
		header := new(HeaderBlock)
		h.Parent().ReplaceChild(h.Parent(), h, header)
		header.AppendChild(header, h)
		for {
			if next == nil ||
				next.Kind() != ast.KindHeading ||
				next.HasBlankPreviousLines() ||
				next.(*ast.Heading).Level != h.Level {
				break
			}
			p := ast.NewParagraph()
			for _, attr := range next.Attributes() {
				p.SetAttribute(attr.Name, attr.Value)
			}
			for ch := next.FirstChild(); ch != nil; ch = ch.NextSibling() {
				p.AppendChild(p, ch)
			}
			header.AppendChild(header, p)
			next.Parent().RemoveChild(next.Parent(), next)
			next = next.NextSibling()
		}
	}
}

type walkerFunc = func(ast.Node) (ast.WalkStatus, error)

func walk(node ast.Node, f walkerFunc) error {
	status, err := f(node)
	if err != nil || status == ast.WalkStop || status == ast.WalkSkipChildren {
		return err
	}
	for n := node.FirstChild(); n != nil; n = n.NextSibling() {
		if err = walk(n, f); err != nil {
			return err
		}
	}
	return nil
}

// Extend implement goldmark.Extender interface.
func (n *HeaderBlock) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithASTTransformers(
		util.Prioritized(n, 0),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(n, 0),
	))
}

// Extension is a initialized goldmark extension for sub-headers support.
var Extension = new(HeaderBlock)

// Option is goldmark.Option for sub-headers extension.
var Option = goldmark.WithExtensions(Extension)
