package subscript

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	// "github.com/yuin/goldmark"
	// "github.com/yuin/goldmark/ast"
	// "github.com/yuin/goldmark/parser"
	// "github.com/yuin/goldmark/renderer"
	// "github.com/yuin/goldmark/text"
	// "github.com/yuin/goldmark/util"
)

var Kind = ast.NewNodeKind("Subscript")

type Node struct {
	ast.BaseInline
}

func (*Node) Kind() ast.NodeKind {
	return Kind
}

type subscript struct {}

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
	m.Parser().AddOptions()
	m.Renderer().AddOptions()
}
