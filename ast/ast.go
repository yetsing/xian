package ast

import (
	"fmt"

	"github.com/yetsing/xian/token"
)

const (
	ExprContextLoad  = "load"
	ExprContextStore = "store"
	ExprContextParam = "param"
)

type ChildNode struct {
	Name   string
	Value  Node
	Values []Node
}

// The Node interface
type Node interface {
	Pos() token.Position
	Token() token.Token
	String() string
	Dumps(indent int) string
	ChildNodes() []ChildNode
}

type BaseNode struct {
	token token.Token
}

func NewBaseNode(token token.Token) BaseNode {
	return BaseNode{token: token}
}

func (bn *BaseNode) Pos() token.Position {
	return bn.token.Start
}

func (bn *BaseNode) Token() token.Token {
	return bn.token
}

func (bn *BaseNode) String() string {
	return bn.token.Literal
}

func (bn *BaseNode) Dumps(indent int) string {
	return "NotImplemented"
}

func (bn *BaseNode) ChildNodes() []ChildNode {
	return nil
}

func (bn *BaseNode) CanAssign() bool {
	return false
}

type Template struct {
	BaseNode

	Body []Node
}

func NewTemplate(token token.Token, body []Node) *Template {
	return &Template{
		BaseNode: NewBaseNode(token),
		Body:     body,
	}
}

func (t *Template) String() string {
	return fmt.Sprintf("nodes.Template(body=%s)", reprNodeList(t.Body))
}

func (t *Template) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Template(")
	sb.WriteNodeList("body", t.Body)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (t *Template) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "body", Values: t.Body},
	}
}
