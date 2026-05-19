package ast

type StatementImpl struct {
	BaseNode
}

func (s *StatementImpl) statementNode() {

}

type Output struct {
	StatementImpl

	nodes []Expression
}

type Extends struct {
	StatementImpl

	template Expression
}

type For struct {
	StatementImpl

	target    Node
	iter      Node
	body      []Node
	else_     []Node
	test      Node
	recursive bool
}

type If struct {
	StatementImpl

	test  Node
	body  []Node
	elif  []*If
	else_ []Node
}

type Macro struct {
	StatementImpl

	name     string
	args     []*Name
	defaults []Expression
	body     []Node
}

type CallBlock struct {
	StatementImpl

	call     *Call
	args     []*Name
	defaults []Expression
	body     []Node
}

type FilterBlock struct {
	StatementImpl

	body   []Node
	filter *Filter
}

type With struct {
	StatementImpl

	targets []Expression
	values  []Expression
	body    []Node
}

type Block struct {
	StatementImpl

	name string
	body []Node
}

type Include struct {
	StatementImpl

	template      Expression
	withContext   bool
	ignoreMissing bool
}

type Import struct {
	StatementImpl

	template    Expression
	target      string
	withContext bool
}

type FromImport struct {
	StatementImpl

	template    Expression
	names       []string
	withContext bool
}

type ExprStmt struct {
	StatementImpl

	node Node
}

type Assign struct {
	StatementImpl

	target Expression
	node   Node
}

type AssignBlock struct {
	StatementImpl

	target Expression
	filter *Filter
	body   []Node
}

type Continue struct {
	StatementImpl
}

type Break struct {
	StatementImpl
}

type Scope struct {
	StatementImpl

	body []Node
}

type OverlayScopy struct {
	StatementImpl

	context Expression
	body    []Node
}

type EvalContextModifier struct {
	StatementImpl

	options []*Keyword
}

type ScopedEvalContextModifier struct {
	EvalContextModifier

	body []Node
}
