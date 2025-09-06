package subscript

import (
	"unicode"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// KindSubscript is a NodeKind of the Subscript node.
var KindSubscript = ast.NewNodeKind("Subscript")

// Node struct represents a subscript node.
type Node struct {
	ast.BaseInline
}

// Kind implements ast.Node.Kind.
func (*Node) Kind() ast.NodeKind {
	return KindSubscript
}

// Dump implements ast.Node.Dump.
func (n *Node) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

// NewSubscriptNode returns a new Subscript node.
func NewSubscriptNode() *Node {
	return &Node{}
}

// subscriptParser is an inline parser for subscript.
type subscriptParser struct {
}

var defaultSubscriptParser = &subscriptParser{}

// NewSubscriptParser returns a new InlineParser that parses subscript expressions.
func NewSubscriptParser() parser.InlineParser {
	return defaultSubscriptParser
}

// Trigger implements parser.InlineParser.Trigger.
func (s *subscriptParser) Trigger() []byte {
	return []byte{'~'}
}

// Parse implements parser.InlineParser.Parse.
func (s *subscriptParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	before := block.PrecendingCharacter()
	line, segment := block.PeekLine()

	// Check if we have at least one character after the tilde
	if len(line) < 2 {
		return nil
	}

	// If preceded by whitespace or is first character of line, not a subscript
	if unicode.IsSpace(before) || before == -1 {
		return nil
	}

	// If we have two tildes in sequence, this should be handled by strikethrough
	if len(line) >= 2 && line[1] == '~' {
		return nil
	}

	// Find the content between tildes
	start := 1 // Skip the opening tilde
	end := -1

	// Look for the closing tilde
	for i := start; i < len(line); i++ {
		if line[i] == '~' {
			end = i
			break
		}
	}

	// If no closing tilde found on this line, not a subscript
	if end == -1 {
		return nil
	}

	// Check if there's any content between tildes
	if end <= start {
		return nil
	}

	content := line[start:end]

	// Check if content has any whitespace (not allowed in subscript)
	for _, b := range content {
		if unicode.IsSpace(rune(b)) {
			return nil
		}
	}

	// Check first character requirements: allow any non-whitespace character except tilde
	firstChar := rune(content[0])
	if firstChar == '~' {
		return nil
	}

	// All subsequent characters are allowed except tilde (handled by finding closing tilde above)
	// No additional character validation needed since whitespace is already checked above

	// Create the subscript node
	node := NewSubscriptNode()

	// Advance past the opening tilde
	block.Advance(1)

	// Parse the content inside - create a text segment for the content
	tempSegment := segment.WithStart(segment.Start + start)
	contentSegment := tempSegment.WithStop(segment.Start + end)
	node.AppendChild(node, ast.NewTextSegment(contentSegment))

	// Advance past the content and closing tilde
	block.Advance(end)

	return node
}

// CloseBlock implements parser.InlineParser.CloseBlock.
func (s *subscriptParser) CloseBlock(parent ast.Node, pc parser.Context) {
	// nothing to do
}

// SubscriptHTMLRenderer is a renderer.NodeRenderer implementation that renders Subscript nodes.
type SubscriptHTMLRenderer struct {
	html.Config
}

// NewSubscriptHTMLRenderer returns a new SubscriptHTMLRenderer.
func NewSubscriptHTMLRenderer(opts ...html.Option) renderer.NodeRenderer {
	r := &SubscriptHTMLRenderer{
		Config: html.NewConfig(),
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

// RegisterFuncs implements renderer.NodeRenderer.RegisterFuncs.
func (r *SubscriptHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindSubscript, r.renderSubscript)
}

// SubscriptAttributeFilter defines attribute names which subscript elements can have.
var SubscriptAttributeFilter = html.GlobalAttributeFilter

func (r *SubscriptHTMLRenderer) renderSubscript(
	w util.BufWriter, source []byte, n ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		if n.Attributes() != nil {
			_, _ = w.WriteString("<sub")
			html.RenderAttributes(w, n, SubscriptAttributeFilter)
			_ = w.WriteByte('>')
		} else {
			_, _ = w.WriteString("<sub>")
		}
	} else {
		_, _ = w.WriteString("</sub>")
	}
	return ast.WalkContinue, nil
}

type subscript struct{}

type SubscriptOption func(*subscript)

var Subscript = NewSubscript()

func NewSubscript(opts ...SubscriptOption) *subscript {
	s := &subscript{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *subscript) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewSubscriptParser(), 100),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewSubscriptHTMLRenderer(), 100),
	))
}
