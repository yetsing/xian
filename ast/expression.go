package ast

import (
	"fmt"

	"github.com/yetsing/xian/token"
)

type BinExpr struct {
	BaseNode

	left     Expression
	operator string
	right    Expression
}

func (be *BinExpr) expressionNode() {}

func (be *BinExpr) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "left", Value: be.left},
		{Name: "right", Value: be.right},
	}
}

type UnaryExpr struct {
	BaseNode

	operator string
	node     Expression
}

func (ue *UnaryExpr) expressionNode() {}

func (ue *UnaryExpr) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "node", Value: ue.node},
	}
}

type Name struct {
	BaseNode

	name string
	ctx  string
}

func NewName(token token.Token, name string, ctx string) *Name {
	n := &Name{
		BaseNode: NewBaseNode(token),
		name:     name,
		ctx:      ctx,
	}
	return n
}

func (n *Name) expressionNode() {}

func (n *Name) String() string {
	return fmt.Sprintf("nodes.Name(name=%s, ctx=%s)", reprString(n.name), reprString(n.ctx))
}

func (n *Name) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Name(")
	fmt.Fprintf(sb, "  name=%s,\n", reprString(n.name))
	sb.WriteIndent()
	fmt.Fprintf(sb, "  ctx=%s,\n", reprString(n.ctx))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (n *Name) ChildNodes() []ChildNode {
	return nil
}

type NSRef struct {
	BaseNode

	name string
	attr string
}

func NewNSRef(token token.Token, name string, attr string) *NSRef {
	ns := &NSRef{
		BaseNode: NewBaseNode(token),
		name:     name,
		attr:     attr,
	}
	return ns
}

func (ns *NSRef) expressionNode() {}

func (ns *NSRef) String() string {
	return fmt.Sprintf("nodes.NSRef(name=%s, attr=%s)", reprString(ns.name), reprString(ns.attr))
}

func (ns *NSRef) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.NSRef(")
	fmt.Fprintf(sb, "  name=%s,\n", reprString(ns.name))
	sb.WriteIndent()
	fmt.Fprintf(sb, "  attr=%s,\n", reprString(ns.attr))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (ns *NSRef) ChildNodes() []ChildNode {
	return nil
}

type Literal interface {
	Expression
	literalNode()
}

type Const struct {
	BaseNode

	value any
}

func NewConst(token token.Token, value any) *Const {
	return &Const{
		BaseNode: NewBaseNode(token),
		value:    value,
	}
}

func (c *Const) expressionNode() {}
func (c *Const) literalNode()    {}

func (c *Const) String() string {
	return fmt.Sprintf("nodes.Const(value=%s)", c.value)
}

func (c *Const) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Const(")
	fmt.Fprintf(sb, "  value=%s,\n", c.value)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

type TemplateData struct {
	BaseNode

	data string
}

func NewTemplateData(token token.Token, data string) *TemplateData {
	td := &TemplateData{
		data: data,
	}
	td.BaseNode = NewBaseNode(token)
	return td
}

func (td *TemplateData) expressionNode() {}
func (td *TemplateData) literalNode()    {}

func (td *TemplateData) String() string {
	return fmt.Sprintf("nodes.TemplateData(data=%s)", reprString(td.data))
}

func (td *TemplateData) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.TemplateData(")
	fmt.Fprintf(sb, "  data=%s,\n", reprString(td.data))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

type Tuple struct {
	BaseNode

	items []Expression
	ctx   string
}

func NewTuple(token token.Token, items []Expression, ctx string) *Tuple {
	t := &Tuple{
		items: items,
		ctx:   ctx,
	}
	t.BaseNode = NewBaseNode(token)
	return t
}

func (t *Tuple) expressionNode() {}
func (t *Tuple) literalNode()    {}

func (t *Tuple) String() string {
	return fmt.Sprintf("nodes.Tuple(items=%s, ctx=%s)", reprExpressionList(t.items), reprString(t.ctx))
}

func (t *Tuple) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Tuple(")
	sb.WriteExpressionList("item", t.items)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  ctx=%s,\n", reprString(t.ctx))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (t *Tuple) ChildNodes() []ChildNode {
	items := make([]Node, len(t.items))
	for i := range t.items {
		items[i] = t.items[i]
	}
	return []ChildNode{
		{Name: "items", Values: items},
	}
}

type List struct {
	BaseNode

	items []Expression
}

func NewList(token token.Token, items []Expression) *List {
	l := &List{
		items: items,
	}
	l.BaseNode = NewBaseNode(token)
	return l
}

func (l *List) expressionNode() {}
func (l *List) literalNode()    {}

func (l *List) String() string {
	return fmt.Sprintf("nodes.List(items=%s)", reprExpressionList(l.items))
}

func (l *List) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.List(")
	sb.WriteExpressionList("item", l.items)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (l *List) ChildNodes() []ChildNode {
	items := make([]Node, len(l.items))
	for i := range l.items {
		items[i] = l.items[i]
	}
	return []ChildNode{
		{Name: "items", Values: items},
	}
}

type Dict struct {
	BaseNode

	items []*Pair
}

func NewDict(token token.Token, items []*Pair) *Dict {
	d := &Dict{
		items: items,
	}
	d.BaseNode = NewBaseNode(token)
	return d
}

func (d *Dict) expressionNode() {}
func (d *Dict) literalNode()    {}

func (d *Dict) String() string {
	items := make([]Node, len(d.items))
	for i, item := range d.items {
		items[i] = item
	}
	return fmt.Sprintf("nodes.Dict(items=%s)", reprNodeList(items))
}

func (d *Dict) Dumps(indent int) string {
	items := make([]Node, len(d.items))
	for i, item := range d.items {
		items[i] = item
	}

	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Dict(")
	sb.WriteNodeList("items", items)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (d *Dict) ChildNodes() []ChildNode {
	items := make([]Node, len(d.items))
	for i := range d.items {
		items[i] = d.items[i]
	}
	return []ChildNode{
		{Name: "items", Values: items},
	}
}

type CondExpr struct {
	BaseNode

	test  Expression
	expr1 Expression
	expr2 Expression
}

func NewCondExpr(token token.Token, test Expression, expr1 Expression, expr2 Expression) *CondExpr {
	ce := &CondExpr{
		test:  test,
		expr1: expr1,
		expr2: expr2,
	}
	ce.BaseNode = NewBaseNode(token)
	return ce
}

func (ce *CondExpr) expressionNode() {}

func (ce *CondExpr) String() string {
	return fmt.Sprintf("nodes.CondExpr(test=%s, expr1=%s, expr2=%s)", ce.test, ce.expr1, ce.expr2)
}

func (ce *CondExpr) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.CondExpr(")
	fmt.Fprintf(sb, "  test=%s,\n", ce.test)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  expr1=%s,\n", ce.expr1)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  expr2=%s,\n", ce.expr2)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (ce *CondExpr) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "test", Value: ce.test},
		{Name: "expr1", Value: ce.expr1},
		{Name: "expr2", Value: ce.expr2},
	}
}

type FilterTestCommon struct {
	BaseNode

	node      Expression
	name      string
	args      []Expression
	kwargs    []*Pair
	dynArgs   Expression
	dynKwargs Expression
	isFilter  bool
}

func (ftc *FilterTestCommon) expressionNode() {}

func (ftc *FilterTestCommon) ChildNodes() []ChildNode {
	items := make([]Node, len(ftc.args))
	for i := range ftc.args {
		items[i] = ftc.args[i]
	}
	kwargs := make([]Node, len(ftc.kwargs))
	for i := range ftc.kwargs {
		kwargs[i] = ftc.kwargs[i]
	}
	return []ChildNode{
		{Name: "node", Value: ftc.node},
		{Name: "args", Values: items},
		{Name: "kwargs", Values: kwargs},
		{Name: "dynArgs", Value: ftc.dynArgs},
		{Name: "dynKwargs", Value: ftc.dynKwargs},
	}
}

type Filter struct {
	FilterTestCommon
}

func NewFilter(token token.Token, node Expression, name string, args []Expression, kwargs []*Pair, dynArgs Expression, dynKwargs Expression) *Filter {
	return &Filter{
		FilterTestCommon: FilterTestCommon{
			BaseNode:  NewBaseNode(token),
			node:      node,
			name:      name,
			args:      args,
			kwargs:    kwargs,
			dynArgs:   dynArgs,
			dynKwargs: dynKwargs,
			isFilter:  true,
		},
	}
}

func (f *Filter) String() string {
	kwargs := make([]Node, len(f.kwargs))
	for i, kwarg := range f.kwargs {
		kwargs[i] = kwarg
	}
	return fmt.Sprintf("nodes.Filter(node=%s, name=%s, args=%s, kwargs=%s, dynArgs=%s, dynKwargs=%s)", f.node, reprString(f.name), reprExpressionList(f.args), reprNodeList(kwargs), f.dynArgs, f.dynKwargs)
}

func (f *Filter) Dumps(indent int) string {
	kwargs := make([]Node, len(f.kwargs))
	for i, kwarg := range f.kwargs {
		kwargs[i] = kwarg
	}

	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Filter(")
	fmt.Fprintf(sb, "  node=%s,\n", f.node)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  name=%s,\n", reprString(f.name))
	sb.WriteIndent()
	sb.WriteExpressionList("args", f.args)
	sb.WriteIndent()
	sb.WriteNodeList("kwargs", kwargs)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  dyn_args=%s,\n", f.dynArgs)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  dyn_kwargs=%s,\n", f.dynKwargs)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

type Test struct {
	FilterTestCommon
}

func NewTest(token token.Token, node Expression, name string, args []Expression, kwargs []*Pair, dynArgs Expression, dynKwargs Expression) *Test {
	return &Test{
		FilterTestCommon: FilterTestCommon{
			BaseNode:  NewBaseNode(token),
			node:      node,
			name:      name,
			args:      args,
			kwargs:    kwargs,
			dynArgs:   dynArgs,
			dynKwargs: dynKwargs,
			isFilter:  false,
		},
	}
}

func (t *Test) String() string {
	kwargs := make([]Node, len(t.kwargs))
	for i, kwarg := range t.kwargs {
		kwargs[i] = kwarg
	}
	return fmt.Sprintf("nodes.Test(node=%s, name=%s, args=%s, kwargs=%s, dynArgs=%s, dynKwargs=%s)", t.node, reprString(t.name), reprExpressionList(t.args), reprNodeList(kwargs), t.dynArgs, t.dynKwargs)
}

func (t *Test) Dumps(indent int) string {
	kwargs := make([]Node, len(t.kwargs))
	for i, kwarg := range t.kwargs {
		kwargs[i] = kwarg
	}

	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Test(")
	fmt.Fprintf(sb, "  node=%s,\n", t.node)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  name=%s,\n", reprString(t.name))
	sb.WriteIndent()
	sb.WriteExpressionList("args", t.args)
	sb.WriteIndent()
	sb.WriteNodeList("kwargs", kwargs)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  dyn_args=%s,\n", t.dynArgs)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  dyn_kwargs=%s,\n", t.dynKwargs)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

type Call struct {
	BaseNode

	node      Expression
	args      []Expression
	kwargs    []*Keyword
	dynArgs   Expression
	dynKwargs Expression
}

func NewCall(token token.Token, node Expression, args []Expression, kwargs []*Keyword, dynArgs Expression, dynKwargs Expression) *Call {
	return &Call{
		BaseNode:  NewBaseNode(token),
		node:      node,
		args:      args,
		kwargs:    kwargs,
		dynArgs:   dynArgs,
		dynKwargs: dynKwargs,
	}
}

func (c *Call) expressionNode() {}

func (c *Call) String() string {
	kwargs := make([]Node, len(c.kwargs))
	for i, kwarg := range c.kwargs {
		kwargs[i] = kwarg
	}
	return fmt.Sprintf("nodes.Call(node=%s, args=%s, kwargs=%s, dynArgs=%s, dynKwargs=%s)", c.node, reprExpressionList(c.args), reprNodeList(kwargs), c.dynArgs, c.dynKwargs)
}

func (c *Call) Dumps(indent int) string {
	kwargs := make([]Node, len(c.kwargs))
	for i, kwarg := range c.kwargs {
		kwargs[i] = kwarg
	}

	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Call(")
	fmt.Fprintf(sb, "  node=%s,\n", c.node)
	sb.WriteIndent()
	sb.WriteExpressionList("args", c.args)
	sb.WriteIndent()
	sb.WriteNodeList("kwargs", kwargs)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  dyn_args=%s,\n", c.dynArgs)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  dyn_kwargs=%s,\n", c.dynKwargs)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (c *Call) ChildNodes() []ChildNode {
	args := make([]Node, len(c.args))
	for i := range c.args {
		args[i] = c.args[i]
	}
	kwargs := make([]Node, len(c.kwargs))
	for i := range c.kwargs {
		kwargs[i] = c.kwargs[i]
	}
	return []ChildNode{
		{Name: "node", Value: c.node},
		{Name: "args", Values: args},
		{Name: "kwargs", Values: kwargs},
		{Name: "dynArgs", Value: c.dynArgs},
		{Name: "dynKwargs", Value: c.dynKwargs},
	}
}

type Getitem struct {
	BaseNode

	node Expression
	arg  Expression
	ctx  string
}

func NewGetitem(token token.Token, node Expression, arg Expression, ctx string) *Getitem {
	gi := &Getitem{
		node: node,
		arg:  arg,
		ctx:  ctx,
	}
	gi.BaseNode = NewBaseNode(token)
	return gi
}

func (gi *Getitem) expressionNode() {}

func (gi *Getitem) String() string {
	return fmt.Sprintf("nodes.Getitem(node=%s, arg=%s, ctx=%s)", gi.node, gi.arg, reprString(gi.ctx))
}

func (gi *Getitem) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Getitem(")
	fmt.Fprintf(sb, "  node=%s,\n", gi.node)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  arg=%s,\n", gi.arg)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  ctx=%s,\n", reprString(gi.ctx))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (gi *Getitem) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "node", Value: gi.node},
		{Name: "arg", Value: gi.arg},
	}
}

type Getattr struct {
	BaseNode

	node Expression
	attr string
	ctx  string
}

func NewGetattr(token token.Token, node Expression, attr string, ctx string) *Getattr {
	ga := &Getattr{
		node: node,
		attr: attr,
		ctx:  ctx,
	}
	ga.BaseNode = NewBaseNode(token)
	return ga
}

func (ga *Getattr) expressionNode() {}

func (ga *Getattr) String() string {
	return fmt.Sprintf("nodes.Getattr(node=%s, attr=%s, ctx=%s)", ga.node, reprString(ga.attr), reprString(ga.ctx))
}

func (ga *Getattr) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Getattr(")
	fmt.Fprintf(sb, "  node=%s,\n", ga.node)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  attr=%s,\n", reprString(ga.attr))
	sb.WriteIndent()
	fmt.Fprintf(sb, "  ctx=%s,\n", reprString(ga.ctx))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (ga *Getattr) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "node", Value: ga.node},
	}
}

type Slice struct {
	BaseNode

	start Expression
	stop  Expression
	step  Expression
}

func NewSlice(token token.Token, start Expression, stop Expression, step Expression) *Slice {
	s := &Slice{
		start: start,
		stop:  stop,
		step:  step,
	}
	s.BaseNode = NewBaseNode(token)
	return s
}

func (s *Slice) expressionNode() {}

func (s *Slice) String() string {
	return fmt.Sprintf("nodes.Slice(start=%s, stop=%s, step=%s)", s.start, s.stop, s.step)
}

func (s *Slice) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Slice(")
	fmt.Fprintf(sb, "  start=%s,\n", s.start)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  stop=%s,\n", s.stop)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  step=%s,\n", s.step)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (s *Slice) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "start", Value: s.start},
		{Name: "stop", Value: s.stop},
		{Name: "step", Value: s.step},
	}
}

type Concat struct {
	BaseNode

	items []Expression
}

func NewConcat(token token.Token, items []Expression) *Concat {
	c := &Concat{
		items: items,
	}
	c.BaseNode = NewBaseNode(token)
	return c
}

func (c *Concat) expressionNode() {}

func (c *Concat) String() string {
	return fmt.Sprintf("nodes.Concat(items=%s)", reprExpressionList(c.items))
}

func (c *Concat) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Concat(")
	sb.WriteExpressionList("items", c.items)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (c *Concat) ChildNodes() []ChildNode {
	items := make([]Node, len(c.items))
	for i := range c.items {
		items[i] = c.items[i]
	}
	return []ChildNode{
		{Name: "items", Values: items},
	}
}

type Compare struct {
	BaseNode

	expr Expression
	ops  []*Operand
}

func NewCompare(token token.Token, expr Expression, ops []*Operand) *Compare {
	return &Compare{
		BaseNode: NewBaseNode(token),
		expr:     expr,
		ops:      ops,
	}
}

func (c *Compare) expressionNode() {}

func (c *Compare) String() string {
	ops := make([]Node, len(c.ops))
	for i := range c.ops {
		ops[i] = c.ops[i]
	}
	return fmt.Sprintf("nodes.Compare(expr=%s, ops=%s)", c.expr, reprNodeList(ops))
}

func (c *Compare) Dumps(indent int) string {
	ops := make([]Node, len(c.ops))
	for i := range c.ops {
		ops[i] = c.ops[i]
	}

	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Compare(")
	fmt.Fprintf(sb, "  expr=%s,\n", c.expr)
	sb.WriteIndent()
	sb.WriteNodeList("ops", ops)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (c *Compare) ChildNodes() []ChildNode {
	ops := make([]Node, len(c.ops))
	for i := range c.ops {
		ops[i] = c.ops[i]
	}
	return []ChildNode{
		{Name: "expr", Value: c.expr},
		{Name: "ops", Values: ops},
	}
}

func stringBinExpr(name string, be *BinExpr) string {
	return fmt.Sprintf("nodes.%s(left=%s, right=%s)", name, be.left, be.right)
}

func dumpsBinExpr(name string, be *BinExpr, indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent(fmt.Sprintf("nodes.%s(", name))
	fmt.Fprintf(sb, "  left=%s,\n", be.left)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  right=%s,\n", be.right)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func stringUnaryExpr(name string, ue *UnaryExpr) string {
	return fmt.Sprintf("nodes.%s(node=%s)", name, ue.node)
}

func dumpsUnaryExpr(name string, ue *UnaryExpr, indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent(fmt.Sprintf("nodes.%s(", name))
	fmt.Fprintf(sb, "  node=%s,\n", ue.node)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

type Mul struct {
	BinExpr
}

func NewMul(token token.Token, left Expression, right Expression) *Mul {
	return &Mul{
		BinExpr: BinExpr{
			BaseNode: NewBaseNode(token),
			left:     left,
			operator: "*",
			right:    right,
		},
	}
}

func (m *Mul) String() string {
	return stringBinExpr("Mul", &m.BinExpr)
}

func (m *Mul) Dumps(indent int) string {
	return dumpsBinExpr("Mul", &m.BinExpr, indent)
}

type Div struct {
	BinExpr
}

func NewDiv(token token.Token, left Expression, right Expression) *Div {
	return &Div{
		BinExpr: BinExpr{
			BaseNode: NewBaseNode(token),
			left:     left,
			operator: "/",
			right:    right,
		},
	}
}

func (d *Div) String() string {
	return stringBinExpr("Div", &d.BinExpr)
}

func (d *Div) Dumps(indent int) string {
	return dumpsBinExpr("Div", &d.BinExpr, indent)
}

type FloorDiv struct {
	BinExpr
}

func NewFloorDiv(token token.Token, left Expression, right Expression) *FloorDiv {
	return &FloorDiv{
		BinExpr: BinExpr{
			BaseNode: NewBaseNode(token),
			left:     left,
			operator: "//",
			right:    right,
		},
	}
}

func (fd *FloorDiv) String() string {
	return stringBinExpr("FloorDiv", &fd.BinExpr)
}

func (fd *FloorDiv) Dumps(indent int) string {
	return dumpsBinExpr("FloorDiv", &fd.BinExpr, indent)
}

type Add struct {
	BinExpr
}

func NewAdd(token token.Token, left Expression, right Expression) *Add {
	return &Add{
		BinExpr: BinExpr{
			BaseNode: NewBaseNode(token),
			left:     left,
			operator: "+",
			right:    right,
		},
	}
}

func (a *Add) String() string {
	return stringBinExpr("Add", &a.BinExpr)
}

func (a *Add) Dumps(indent int) string {
	return dumpsBinExpr("Add", &a.BinExpr, indent)
}

type Sub struct {
	BinExpr
}

func NewSub(token token.Token, left Expression, right Expression) *Sub {
	return &Sub{
		BinExpr: BinExpr{
			BaseNode: NewBaseNode(token),
			left:     left,
			operator: "-",
			right:    right,
		},
	}
}

func (s *Sub) String() string {
	return stringBinExpr("Sub", &s.BinExpr)
}

func (s *Sub) Dumps(indent int) string {
	return dumpsBinExpr("Sub", &s.BinExpr, indent)
}

type Mod struct {
	BinExpr
}

func NewMod(token token.Token, left Expression, right Expression) *Mod {
	return &Mod{
		BinExpr: BinExpr{
			BaseNode: NewBaseNode(token),
			left:     left,
			operator: "%",
			right:    right,
		},
	}
}

func (m *Mod) String() string {
	return stringBinExpr("Mod", &m.BinExpr)
}

func (m *Mod) Dumps(indent int) string {
	return dumpsBinExpr("Mod", &m.BinExpr, indent)
}

type Pow struct {
	BinExpr
}

func NewPow(token token.Token, left Expression, right Expression) *Pow {
	return &Pow{
		BinExpr: BinExpr{
			BaseNode: NewBaseNode(token),
			left:     left,
			operator: "**",
			right:    right,
		},
	}
}

func (p *Pow) String() string {
	return stringBinExpr("Pow", &p.BinExpr)
}

func (p *Pow) Dumps(indent int) string {
	return dumpsBinExpr("Pow", &p.BinExpr, indent)
}

type And struct {
	BinExpr
}

func NewAnd(token token.Token, left Expression, right Expression) *And {
	return &And{
		BinExpr: BinExpr{
			BaseNode: NewBaseNode(token),
			left:     left,
			operator: "and",
			right:    right,
		},
	}
}

func (a *And) String() string {
	return stringBinExpr("And", &a.BinExpr)
}

func (a *And) Dumps(indent int) string {
	return dumpsBinExpr("And", &a.BinExpr, indent)
}

type Or struct {
	BinExpr
}

func NewOr(token token.Token, left Expression, right Expression) *Or {
	return &Or{
		BinExpr: BinExpr{
			BaseNode: NewBaseNode(token),
			left:     left,
			operator: "or",
			right:    right,
		},
	}
}

func (o *Or) String() string {
	return stringBinExpr("Or", &o.BinExpr)
}

func (o *Or) Dumps(indent int) string {
	return dumpsBinExpr("Or", &o.BinExpr, indent)
}

type Not struct {
	UnaryExpr
}

func NewNot(token token.Token, node Expression) *Not {
	return &Not{
		UnaryExpr: UnaryExpr{
			BaseNode: NewBaseNode(token),
			operator: "not",
			node:     node,
		},
	}
}

func (n *Not) String() string {
	return stringUnaryExpr("Not", &n.UnaryExpr)
}

func (n *Not) Dumps(indent int) string {
	return dumpsUnaryExpr("Not", &n.UnaryExpr, indent)
}

type Neg struct {
	UnaryExpr
}

func NewNeg(token token.Token, node Expression) *Neg {
	return &Neg{
		UnaryExpr: UnaryExpr{
			BaseNode: NewBaseNode(token),
			operator: "-",
			node:     node,
		},
	}
}

func (n *Neg) String() string {
	return stringUnaryExpr("Neg", &n.UnaryExpr)
}

func (n *Neg) Dumps(indent int) string {
	return dumpsUnaryExpr("Neg", &n.UnaryExpr, indent)
}

type Pos struct {
	UnaryExpr
}

func NewPos(token token.Token, node Expression) *Pos {
	return &Pos{
		UnaryExpr: UnaryExpr{
			BaseNode: NewBaseNode(token),
			operator: "+",
			node:     node,
		},
	}
}

func (p *Pos) String() string {
	return stringUnaryExpr("Pos", &p.UnaryExpr)
}

func (p *Pos) Dumps(indent int) string {
	return dumpsUnaryExpr("Pos", &p.UnaryExpr, indent)
}

type EnvironmentAttribute struct {
	BaseNode

	name string
}

func NewEnvironmentAttribute(token token.Token, name string) *EnvironmentAttribute {
	return &EnvironmentAttribute{
		BaseNode: NewBaseNode(token),
		name:     name,
	}
}

func (ea *EnvironmentAttribute) expressionNode() {}

func (ea *EnvironmentAttribute) String() string {
	return fmt.Sprintf("nodes.EnvironmentAttribute(name=%s)", reprString(ea.name))
}

func (ea *EnvironmentAttribute) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.EnvironmentAttribute(")
	fmt.Fprintf(sb, "  name=%s,\n", reprString(ea.name))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

type ExtensionAttribute struct {
	BaseNode

	identifier string
	name       string
}

func NewExtensionAttribute(token token.Token, identifier string, name string) *ExtensionAttribute {
	return &ExtensionAttribute{
		BaseNode:   NewBaseNode(token),
		identifier: identifier,
		name:       name,
	}
}

func (ea *ExtensionAttribute) expressionNode() {}

func (ea *ExtensionAttribute) String() string {
	return fmt.Sprintf("nodes.ExtensionAttribute(identifier=%s, name=%s)", reprString(ea.identifier), reprString(ea.name))
}

func (ea *ExtensionAttribute) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.ExtensionAttribute(")
	fmt.Fprintf(sb, "  identifier=%s,\n", reprString(ea.identifier))
	sb.WriteIndent()
	fmt.Fprintf(sb, "  name=%s,\n", reprString(ea.name))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

type ImportedName struct {
	BaseNode

	importname string
}

func NewImportedName(token token.Token, importname string) *ImportedName {
	return &ImportedName{
		BaseNode:   NewBaseNode(token),
		importname: importname,
	}
}

func (in *ImportedName) expressionNode() {}

func (in *ImportedName) String() string {
	return fmt.Sprintf("nodes.ImportedName(importname=%s)", reprString(in.importname))
}

func (in *ImportedName) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.ImportedName(")
	fmt.Fprintf(sb, "  importname=%s,\n", reprString(in.importname))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

type InternalName struct {
	BaseNode

	name string
}

func NewInternalName(token token.Token, name string) *InternalName {
	return &InternalName{
		BaseNode: NewBaseNode(token),
		name:     name,
	}
}

func (in *InternalName) expressionNode() {}

func (in *InternalName) String() string {
	return fmt.Sprintf("nodes.InternalName(name=%s)", reprString(in.name))
}

func (in *InternalName) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.InternalName(")
	fmt.Fprintf(sb, "  name=%s,\n", reprString(in.name))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

type MarkSafe struct {
	BaseNode

	expr Expression
}

func NewMarkSafe(token token.Token, expr Expression) *MarkSafe {
	return &MarkSafe{
		BaseNode: NewBaseNode(token),
		expr:     expr,
	}
}

func (ms *MarkSafe) expressionNode() {}

func (ms *MarkSafe) String() string {
	return fmt.Sprintf("nodes.MarkSafe(expr=%s)", ms.expr)
}

func (ms *MarkSafe) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.MarkSafe(")
	fmt.Fprintf(sb, "  expr=%s,\n", ms.expr)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (ms *MarkSafe) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "expr", Value: ms.expr},
	}
}

type MarkSafeIfAutoescape struct {
	BaseNode

	expr Expression
}

func NewMarkSafeIfAutoescape(token token.Token, expr Expression) *MarkSafeIfAutoescape {
	return &MarkSafeIfAutoescape{
		BaseNode: NewBaseNode(token),
		expr:     expr,
	}
}

func (ms *MarkSafeIfAutoescape) expressionNode() {}

func (ms *MarkSafeIfAutoescape) String() string {
	return fmt.Sprintf("nodes.MarkSafeIfAutoescape(expr=%s)", ms.expr)
}

func (ms *MarkSafeIfAutoescape) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.MarkSafeIfAutoescape(")
	fmt.Fprintf(sb, "  expr=%s,\n", ms.expr)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (ms *MarkSafeIfAutoescape) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "expr", Value: ms.expr},
	}
}

type ContextReference struct {
	BaseNode
}

func NewContextReference(token token.Token) *ContextReference {
	return &ContextReference{
		BaseNode: NewBaseNode(token),
	}
}

func (cr *ContextReference) expressionNode() {}

func (cr *ContextReference) String() string {
	return "nodes.ContextReference()"
}

func (cr *ContextReference) Dumps(indent int) string {
	return "nodes.ContextReference()"
}

type DerivedContextReference struct {
	BaseNode
}

func NewDerivedContextReference(token token.Token) *DerivedContextReference {
	return &DerivedContextReference{
		BaseNode: NewBaseNode(token),
	}
}

func (dcr *DerivedContextReference) expressionNode() {}

func (dcr *DerivedContextReference) String() string {
	return "nodes.DerivedContextReference()"
}

func (dcr *DerivedContextReference) Dumps(indent int) string {
	return "nodes.DerivedContextReference()"
}
