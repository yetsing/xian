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

// All statement nodes implement this
type Statement interface {
	Node
	statementNode()
}

// All expression nodes implement this
type Expression interface {
	Node
	expressionNode()
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

// ================= Helper =================

type Helper interface {
	Node
	helperNode()
}

type Keyword struct {
	BaseNode

	key   string
	value Expression
}

func NewKeyword(token token.Token, key string, value Expression) *Keyword {
	return &Keyword{
		BaseNode: NewBaseNode(token),
		key:      key,
		value:    value,
	}
}

func (k *Keyword) helperNode() {}

func (k *Keyword) String() string {
	return fmt.Sprintf("nodes.Keyword(key=%s, value=%s)", reprString(k.key), k.value.String())
}

func (k *Keyword) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Keyword(")
	fmt.Fprintf(sb, "  key=%s,\n", reprString(k.key))
	sb.WriteIndent()
	fmt.Fprintf(sb, "  value=%s,\n", k.value.Dumps(indent+2))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (k *Keyword) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "value", Value: k.value},
	}
}

type Pair struct {
	BaseNode

	key   Expression
	value Expression
}

func (p *Pair) helperNode() {}

func (p *Pair) String() string {
	return fmt.Sprintf("nodes.Pair(key=%s, value=%s)", p.key.String(), p.value.String())
}

func (p *Pair) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Pair(")
	fmt.Fprintf(sb, "  key=%s,\n", p.key.Dumps(indent+2))
	sb.WriteIndent()
	fmt.Fprintf(sb, "  value=%s,\n", p.value.Dumps(indent+2))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (p *Pair) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "key", Value: p.key},
		{Name: "value", Value: p.value},
	}
}

type Operand struct {
	BaseNode

	op   string
	expr Expression
}

func (o *Operand) helperNode() {}

func (o *Operand) String() string {
	return fmt.Sprintf("nodes.Operand(op=%s, expr=%s)", reprString(o.op), o.expr.String())
}

func (o *Operand) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Operand(")
	fmt.Fprintf(sb, "  op=%s,\n", reprString(o.op))
	sb.WriteIndent()
	fmt.Fprintf(sb, "  expr=%s,\n", o.expr.Dumps(indent+2))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (o *Operand) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "expr", Value: o.expr},
	}
}
