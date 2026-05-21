package ast

type ExpressionImpl struct {
	BaseNode
}

func (expr *ExpressionImpl) expressionNode() {

}

type BinExpr struct {
	ExpressionImpl

	left     Expression
	operator string
	right    Expression
}

type UnaryExpr struct {
	ExpressionImpl

	operator string
	node     Expression
}

type Name struct {
	ExpressionImpl

	name string
	ctx  string
}

type NSRef struct {
	ExpressionImpl

	name string
	attr string
}

type Literal interface {
	Expression
	literalNode()
}

type LiteralImpl struct {
	ExpressionImpl
}

func (l *LiteralImpl) literalNode() {
}

type Const struct {
	LiteralImpl

	value any
}

type TemplateData struct {
	LiteralImpl

	data string
}

type Tuple struct {
	LiteralImpl

	items []Expression
	ctx   string
}

type List struct {
	LiteralImpl

	items []Expression
}

type Dict struct {
	LiteralImpl

	items []*Pair
}

type CondExpr struct {
	ExpressionImpl

	test  Expression
	expr1 Expression
	expr2 Expression
}

type FilterTestCommon struct {
	ExpressionImpl

	node      Expression
	name      string
	args      []Expression
	kwargs    []*Pair
	dynArgs   Expression
	dynKwargs Expression
}

type Filter struct {
	FilterTestCommon
}

type Test struct {
	FilterTestCommon
}

type Call struct {
	ExpressionImpl

	node      Expression
	args      []Expression
	kwargs    []*Keyword
	dynArgs   Expression
	dynKwargs Expression
}

type Getitem struct {
	ExpressionImpl

	node Expression
	arg  Expression
	ctx  string
}

type Getattr struct {
	ExpressionImpl

	node Expression
	attr string
	ctx  string
}

type Slice struct {
	ExpressionImpl

	start Expression
	stop  Expression
	step  Expression
}

type Concat struct {
	ExpressionImpl

	items []Expression
}

type Compare struct {
	ExpressionImpl

	expr Expression
	ops  []*Operand
}

type Mul struct {
	BinExpr
}

type Div struct {
	BinExpr
}

type FloorDiv struct {
	BinExpr
}

type Add struct {
	BinExpr
}

type Sub struct {
	BinExpr
}

type Mod struct {
	BinExpr
}

type Pow struct {
	BinExpr
}

type And struct {
	BinExpr
}

type Or struct {
	BinExpr
}

type Not struct {
	UnaryExpr
}

type Neg struct {
	UnaryExpr
}

type Pos struct {
	UnaryExpr
}

type EnvironmentAttribute struct {
	ExpressionImpl

	name string
}

type ExtensionAttribute struct {
	ExpressionImpl

	identifier string
	name       string
}

type ImportedName struct {
	ExpressionImpl

	importname string
}

type InternalName struct {
	ExpressionImpl

	name string
}

type MarkSafe struct {
	ExpressionImpl

	expr Expression
}

type MarkSafeIfAutoescape struct {
	ExpressionImpl

	expr Expression
}

type ContextReference struct {
	ExpressionImpl
}

type DerivedContextReference struct {
	ExpressionImpl
}
