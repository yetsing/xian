package ast

import (
	"github.com/yetsing/xian"
	"github.com/yetsing/xian/internal/token"
)

// The Node interface
type Node interface {
	Pos() token.Position
	TokenLiteral() string
	String() string
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
	token       token.Token
	environment *xian.Environment
}

func (bn *BaseNode) Pos() token.Position {
	return bn.token.Start
}

func (bn *BaseNode) TokenLiteral() string {
	return bn.token.Literal
}

func (bn *BaseNode) String() string {
	return bn.token.Literal
}

type Template struct {
	BaseNode
	body []Node
}

// ================= Helper =================

type Helper interface {
	Node
	helperNode()
}

type HelperImpl struct {
	BaseNode
}

func (h *HelperImpl) helperNode() {}

type Keyword struct {
	HelperImpl

	key   string
	value Expression
}

type Pair struct {
	HelperImpl

	key   Expression
	value Expression
}

type Operand struct {
	BaseNode

	op   string
	expr Expression
}
