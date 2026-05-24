package ast

import (
	"fmt"

	"github.com/yetsing/xian/token"
)

type Helper interface {
	Node
	helperNode()
}

type Keyword struct {
	BaseNode

	Key   string
	Value Expression
}

func NewKeyword(token token.Token, key string, value Expression) *Keyword {
	return &Keyword{
		BaseNode: NewBaseNode(token),
		Key:      key,
		Value:    value,
	}
}

func (k *Keyword) helperNode() {}

func (k *Keyword) String() string {
	return fmt.Sprintf("nodes.Keyword(key=%s, value=%s)", reprString(k.Key), k.Value.String())
}

func (k *Keyword) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Keyword(")
	fmt.Fprintf(sb, "  key=%s,\n", reprString(k.Key))
	sb.WriteIndent()
	fmt.Fprintf(sb, "  value=%s,\n", k.Value.Dumps(indent+2))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (k *Keyword) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "value", Value: k.Value},
	}
}

type Pair struct {
	BaseNode

	Key   Expression
	Value Expression
}

func NewPair(token token.Token, key Expression, value Expression) *Pair {
	return &Pair{
		BaseNode: NewBaseNode(token),
		Key:      key,
		Value:    value,
	}
}

func (p *Pair) helperNode() {}

func (p *Pair) String() string {
	return fmt.Sprintf("nodes.Pair(key=%s, value=%s)", p.Key.String(), p.Value.String())
}

func (p *Pair) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Pair(")
	fmt.Fprintf(sb, "  key=%s,\n", p.Key.Dumps(indent+2))
	sb.WriteIndent()
	fmt.Fprintf(sb, "  value=%s,\n", p.Value.Dumps(indent+2))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (p *Pair) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "key", Value: p.Key},
		{Name: "value", Value: p.Value},
	}
}

type Operand struct {
	BaseNode

	Op   string
	Expr Expression
}

func NewOperand(token token.Token, op string, expr Expression) *Operand {
	return &Operand{
		BaseNode: NewBaseNode(token),
		Op:       op,
		Expr:     expr,
	}
}

func (o *Operand) helperNode() {}

func (o *Operand) String() string {
	return fmt.Sprintf("nodes.Operand(op=%s, expr=%s)", reprString(o.Op), o.Expr.String())
}

func (o *Operand) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Operand(")
	fmt.Fprintf(sb, "  op=%s,\n", reprString(o.Op))
	sb.WriteIndent()
	fmt.Fprintf(sb, "  expr=%s,\n", o.Expr.Dumps(indent+2))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (o *Operand) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "expr", Value: o.Expr},
	}
}
