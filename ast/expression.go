package ast

import (
	"fmt"

	"github.com/yetsing/xian/token"
)

// All expression nodes implement this
type Expression interface {
	Node
	expressionNode()
	CanAssign() bool
}

type BinExpr struct {
	BaseNode

	Left     Expression
	Operator string
	Right    Expression
}

func (be *BinExpr) expressionNode() {}

func (be *BinExpr) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "left", Value: be.Left},
		{Name: "right", Value: be.Right},
	}
}

type UnaryExpr struct {
	BaseNode

	Operator string
	Node     Expression
}

func (ue *UnaryExpr) expressionNode() {}

func (ue *UnaryExpr) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "node", Value: ue.Node},
	}
}

type Name struct {
	BaseNode

	Name string
	Ctx  string
}

func NewName(token token.Token, name string, ctx string) *Name {
	n := &Name{
		BaseNode: NewBaseNode(token),
		Name:     name,
		Ctx:      ctx,
	}
	return n
}

func (n *Name) expressionNode() {}

func (n *Name) CanAssign() bool {
	switch n.Name {
	case "true", "false", "True", "False", "none", "None":
		return false
	default:
		return true
	}
}

func (n *Name) String() string {
	return fmt.Sprintf("nodes.Name(name=%s, ctx=%s)", reprString(n.Name), reprString(n.Ctx))
}

func (n *Name) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Name(")
	fmt.Fprintf(sb, "  name=%s,\n", reprString(n.Name))
	sb.WriteIndent()
	fmt.Fprintf(sb, "  ctx=%s,\n", reprString(n.Ctx))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (n *Name) ChildNodes() []ChildNode {
	return nil
}

type NSRef struct {
	BaseNode

	Name string
	Attr string
}

func NewNSRef(token token.Token, name string, attr string) *NSRef {
	ns := &NSRef{
		BaseNode: NewBaseNode(token),
		Name:     name,
		Attr:     attr,
	}
	return ns
}

func (ns *NSRef) expressionNode() {}

func (ns *NSRef) String() string {
	return fmt.Sprintf("nodes.NSRef(name=%s, attr=%s)", reprString(ns.Name), reprString(ns.Attr))
}

func (ns *NSRef) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.NSRef(")
	fmt.Fprintf(sb, "  name=%s,\n", reprString(ns.Name))
	sb.WriteIndent()
	fmt.Fprintf(sb, "  attr=%s,\n", reprString(ns.Attr))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (ns *NSRef) ChildNodes() []ChildNode {
	return nil
}

func (ns *NSRef) CanAssign() bool {
	return true
}

type Literal interface {
	Expression
	literalNode()
}

type Const struct {
	BaseNode

	Value any
}

func NewConst(token token.Token, value any) *Const {
	return &Const{
		BaseNode: NewBaseNode(token),
		Value:    value,
	}
}

func (c *Const) expressionNode() {}
func (c *Const) literalNode()    {}

func (c *Const) String() string {
	return fmt.Sprintf("nodes.Const(value=%s)", c.Value)
}

func (c *Const) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Const(")
	fmt.Fprintf(sb, "  value=%s,\n", c.Value)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

type TemplateData struct {
	BaseNode

	Data string
}

func NewTemplateData(token token.Token, data string) *TemplateData {
	td := &TemplateData{
		Data: data,
	}
	td.BaseNode = NewBaseNode(token)
	return td
}

func (td *TemplateData) expressionNode() {}
func (td *TemplateData) literalNode()    {}

func (td *TemplateData) String() string {
	return fmt.Sprintf("nodes.TemplateData(data=%s)", reprString(td.Data))
}

func (td *TemplateData) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.TemplateData(")
	fmt.Fprintf(sb, "  data=%s,\n", reprString(td.Data))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

type Tuple struct {
	BaseNode

	Items []Expression
	Ctx   string
}

func NewTuple(token token.Token, items []Expression, ctx string) *Tuple {
	t := &Tuple{
		Items: items,
		Ctx:   ctx,
	}
	t.BaseNode = NewBaseNode(token)
	return t
}

func (t *Tuple) expressionNode() {}
func (t *Tuple) literalNode()    {}

func (t *Tuple) CanAssign() bool {
	for _, item := range t.Items {
		if !item.CanAssign() {
			return false
		}
	}
	return true
}

func (t *Tuple) String() string {
	return fmt.Sprintf("nodes.Tuple(items=%s, ctx=%s)", reprExpressionList(t.Items), reprString(t.Ctx))
}

func (t *Tuple) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Tuple(")
	sb.WriteExpressionList("item", t.Items)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  ctx=%s,\n", reprString(t.Ctx))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (t *Tuple) ChildNodes() []ChildNode {
	items := make([]Node, len(t.Items))
	for i := range t.Items {
		items[i] = t.Items[i]
	}
	return []ChildNode{
		{Name: "items", Values: items},
	}
}

type List struct {
	BaseNode

	Items []Expression
}

func NewList(token token.Token, items []Expression) *List {
	l := &List{
		Items: items,
	}
	l.BaseNode = NewBaseNode(token)
	return l
}

func (l *List) expressionNode() {}
func (l *List) literalNode()    {}

func (l *List) String() string {
	return fmt.Sprintf("nodes.List(items=%s)", reprExpressionList(l.Items))
}

func (l *List) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.List(")
	sb.WriteExpressionList("item", l.Items)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (l *List) ChildNodes() []ChildNode {
	items := make([]Node, len(l.Items))
	for i := range l.Items {
		items[i] = l.Items[i]
	}
	return []ChildNode{
		{Name: "items", Values: items},
	}
}

type Dict struct {
	BaseNode

	Items []*Pair
}

func NewDict(token token.Token, items []*Pair) *Dict {
	d := &Dict{
		Items: items,
	}
	d.BaseNode = NewBaseNode(token)
	return d
}

func (d *Dict) expressionNode() {}
func (d *Dict) literalNode()    {}

func (d *Dict) String() string {
	items := make([]Node, len(d.Items))
	for i, item := range d.Items {
		items[i] = item
	}
	return fmt.Sprintf("nodes.Dict(items=%s)", reprNodeList(items))
}

func (d *Dict) Dumps(indent int) string {
	items := make([]Node, len(d.Items))
	for i, item := range d.Items {
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
	items := make([]Node, len(d.Items))
	for i := range d.Items {
		items[i] = d.Items[i]
	}
	return []ChildNode{
		{Name: "items", Values: items},
	}
}

type CondExpr struct {
	BaseNode

	Test  Expression
	Expr1 Expression
	Expr2 Expression
}

func NewCondExpr(token token.Token, test Expression, expr1 Expression, expr2 Expression) *CondExpr {
	ce := &CondExpr{
		Test:  test,
		Expr1: expr1,
		Expr2: expr2,
	}
	ce.BaseNode = NewBaseNode(token)
	return ce
}

func (ce *CondExpr) expressionNode() {}

func (ce *CondExpr) String() string {
	return fmt.Sprintf("nodes.CondExpr(test=%s, expr1=%s, expr2=%s)", ce.Test, ce.Expr1, ce.Expr2)
}

func (ce *CondExpr) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.CondExpr(")
	fmt.Fprintf(sb, "  test=%s,\n", ce.Test)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  expr1=%s,\n", ce.Expr1)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  expr2=%s,\n", ce.Expr2)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (ce *CondExpr) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "test", Value: ce.Test},
		{Name: "expr1", Value: ce.Expr1},
		{Name: "expr2", Value: ce.Expr2},
	}
}

type FilterTestCommon struct {
	BaseNode

	Node      Expression
	Name      string
	Args      []Expression
	Kwargs    []*Keyword
	DynArgs   Expression
	DynKwargs Expression
	IsFilter  bool
}

func (ftc *FilterTestCommon) expressionNode() {}

func (ftc *FilterTestCommon) ChildNodes() []ChildNode {
	items := make([]Node, len(ftc.Args))
	for i := range ftc.Args {
		items[i] = ftc.Args[i]
	}
	kwargs := make([]Node, len(ftc.Kwargs))
	for i := range ftc.Kwargs {
		kwargs[i] = ftc.Kwargs[i]
	}
	return []ChildNode{
		{Name: "node", Value: ftc.Node},
		{Name: "args", Values: items},
		{Name: "kwargs", Values: kwargs},
		{Name: "dynArgs", Value: ftc.DynArgs},
		{Name: "dynKwargs", Value: ftc.DynKwargs},
	}
}

type Filter struct {
	FilterTestCommon
}

func NewFilter(token token.Token, node Expression, name string, args []Expression, kwargs []*Keyword, dynArgs Expression, dynKwargs Expression) *Filter {
	return &Filter{
		FilterTestCommon: FilterTestCommon{
			BaseNode:  NewBaseNode(token),
			Node:      node,
			Name:      name,
			Args:      args,
			Kwargs:    kwargs,
			DynArgs:   dynArgs,
			DynKwargs: dynKwargs,
			IsFilter:  true,
		},
	}
}

func (f *Filter) String() string {
	kwargs := make([]Node, len(f.Kwargs))
	for i, kwarg := range f.Kwargs {
		kwargs[i] = kwarg
	}
	return fmt.Sprintf("nodes.Filter(node=%s, name=%s, args=%s, kwargs=%s, dynArgs=%s, dynKwargs=%s)", f.Node, reprString(f.Name), reprExpressionList(f.Args), reprNodeList(kwargs), f.DynArgs, f.DynKwargs)
}

func (f *Filter) Dumps(indent int) string {
	kwargs := make([]Node, len(f.Kwargs))
	for i, kwarg := range f.Kwargs {
		kwargs[i] = kwarg
	}

	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Filter(")
	fmt.Fprintf(sb, "  node=%s,\n", f.Node)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  name=%s,\n", reprString(f.Name))
	sb.WriteIndent()
	sb.WriteExpressionList("args", f.Args)
	sb.WriteIndent()
	sb.WriteNodeList("kwargs", kwargs)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  dyn_args=%s,\n", f.DynArgs)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  dyn_kwargs=%s,\n", f.DynKwargs)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

type Test struct {
	FilterTestCommon
}

func NewTest(token token.Token, node Expression, name string, args []Expression, kwargs []*Keyword, dynArgs Expression, dynKwargs Expression) *Test {
	return &Test{
		FilterTestCommon: FilterTestCommon{
			BaseNode:  NewBaseNode(token),
			Node:      node,
			Name:      name,
			Args:      args,
			Kwargs:    kwargs,
			DynArgs:   dynArgs,
			DynKwargs: dynKwargs,
			IsFilter:  false,
		},
	}
}

func (t *Test) String() string {
	kwargs := make([]Node, len(t.Kwargs))
	for i, kwarg := range t.Kwargs {
		kwargs[i] = kwarg
	}
	return fmt.Sprintf("nodes.Test(node=%s, name=%s, args=%s, kwargs=%s, dynArgs=%s, dynKwargs=%s)", t.Node, reprString(t.Name), reprExpressionList(t.Args), reprNodeList(kwargs), t.DynArgs, t.DynKwargs)
}

func (t *Test) Dumps(indent int) string {
	kwargs := make([]Node, len(t.Kwargs))
	for i, kwarg := range t.Kwargs {
		kwargs[i] = kwarg
	}

	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Test(")
	fmt.Fprintf(sb, "  node=%s,\n", t.Node)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  name=%s,\n", reprString(t.Name))
	sb.WriteIndent()
	sb.WriteExpressionList("args", t.Args)
	sb.WriteIndent()
	sb.WriteNodeList("kwargs", kwargs)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  dyn_args=%s,\n", t.DynArgs)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  dyn_kwargs=%s,\n", t.DynKwargs)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

type Call struct {
	BaseNode

	Node      Expression
	Args      []Expression
	Kwargs    []*Keyword
	DynArgs   Expression
	DynKwargs Expression
}

func NewCall(token token.Token, node Expression, args []Expression, kwargs []*Keyword, dynArgs Expression, dynKwargs Expression) *Call {
	return &Call{
		BaseNode:  NewBaseNode(token),
		Node:      node,
		Args:      args,
		Kwargs:    kwargs,
		DynArgs:   dynArgs,
		DynKwargs: dynKwargs,
	}
}

func (c *Call) expressionNode() {}

func (c *Call) String() string {
	kwargs := make([]Node, len(c.Kwargs))
	for i, kwarg := range c.Kwargs {
		kwargs[i] = kwarg
	}
	return fmt.Sprintf("nodes.Call(node=%s, args=%s, kwargs=%s, dynArgs=%s, dynKwargs=%s)", c.Node, reprExpressionList(c.Args), reprNodeList(kwargs), c.DynArgs, c.DynKwargs)
}

func (c *Call) Dumps(indent int) string {
	kwargs := make([]Node, len(c.Kwargs))
	for i, kwarg := range c.Kwargs {
		kwargs[i] = kwarg
	}

	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Call(")
	fmt.Fprintf(sb, "  node=%s,\n", c.Node)
	sb.WriteIndent()
	sb.WriteExpressionList("args", c.Args)
	sb.WriteIndent()
	sb.WriteNodeList("kwargs", kwargs)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  dyn_args=%s,\n", c.DynArgs)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  dyn_kwargs=%s,\n", c.DynKwargs)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (c *Call) ChildNodes() []ChildNode {
	args := make([]Node, len(c.Args))
	for i := range c.Args {
		args[i] = c.Args[i]
	}
	kwargs := make([]Node, len(c.Kwargs))
	for i := range c.Kwargs {
		kwargs[i] = c.Kwargs[i]
	}
	return []ChildNode{
		{Name: "node", Value: c.Node},
		{Name: "args", Values: args},
		{Name: "kwargs", Values: kwargs},
		{Name: "dynArgs", Value: c.DynArgs},
		{Name: "dynKwargs", Value: c.DynKwargs},
	}
}

type Getitem struct {
	BaseNode

	Node Expression
	Arg  Expression
	Ctx  string
}

func NewGetitem(token token.Token, node Expression, arg Expression, ctx string) *Getitem {
	gi := &Getitem{
		Node: node,
		Arg:  arg,
		Ctx:  ctx,
	}
	gi.BaseNode = NewBaseNode(token)
	return gi
}

func (gi *Getitem) expressionNode() {}

func (gi *Getitem) String() string {
	return fmt.Sprintf("nodes.Getitem(node=%s, arg=%s, ctx=%s)", gi.Node, gi.Arg, reprString(gi.Ctx))
}

func (gi *Getitem) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Getitem(")
	fmt.Fprintf(sb, "  node=%s,\n", gi.Node)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  arg=%s,\n", gi.Arg)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  ctx=%s,\n", reprString(gi.Ctx))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (gi *Getitem) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "node", Value: gi.Node},
		{Name: "arg", Value: gi.Arg},
	}
}

type Getattr struct {
	BaseNode

	Node Expression
	Attr string
	Ctx  string
}

func NewGetattr(token token.Token, node Expression, attr string, ctx string) *Getattr {
	ga := &Getattr{
		Node: node,
		Attr: attr,
		Ctx:  ctx,
	}
	ga.BaseNode = NewBaseNode(token)
	return ga
}

func (ga *Getattr) expressionNode() {}

func (ga *Getattr) String() string {
	return fmt.Sprintf("nodes.Getattr(node=%s, attr=%s, ctx=%s)", ga.Node, reprString(ga.Attr), reprString(ga.Ctx))
}

func (ga *Getattr) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Getattr(")
	fmt.Fprintf(sb, "  node=%s,\n", ga.Node)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  attr=%s,\n", reprString(ga.Attr))
	sb.WriteIndent()
	fmt.Fprintf(sb, "  ctx=%s,\n", reprString(ga.Ctx))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (ga *Getattr) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "node", Value: ga.Node},
	}
}

type Slice struct {
	BaseNode

	Start Expression
	Stop  Expression
	Step  Expression
}

func NewSlice(token token.Token, start Expression, stop Expression, step Expression) *Slice {
	s := &Slice{
		Start: start,
		Stop:  stop,
		Step:  step,
	}
	s.BaseNode = NewBaseNode(token)
	return s
}

func (s *Slice) expressionNode() {}

func (s *Slice) String() string {
	return fmt.Sprintf("nodes.Slice(start=%s, stop=%s, step=%s)", s.Start, s.Stop, s.Step)
}

func (s *Slice) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Slice(")
	fmt.Fprintf(sb, "  start=%s,\n", s.Start)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  stop=%s,\n", s.Stop)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  step=%s,\n", s.Step)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (s *Slice) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "start", Value: s.Start},
		{Name: "stop", Value: s.Stop},
		{Name: "step", Value: s.Step},
	}
}

type Concat struct {
	BaseNode

	Items []Expression
}

func NewConcat(token token.Token, items []Expression) *Concat {
	c := &Concat{
		Items: items,
	}
	c.BaseNode = NewBaseNode(token)
	return c
}

func (c *Concat) expressionNode() {}

func (c *Concat) String() string {
	return fmt.Sprintf("nodes.Concat(items=%s)", reprExpressionList(c.Items))
}

func (c *Concat) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Concat(")
	sb.WriteExpressionList("items", c.Items)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (c *Concat) ChildNodes() []ChildNode {
	items := make([]Node, len(c.Items))
	for i := range c.Items {
		items[i] = c.Items[i]
	}
	return []ChildNode{
		{Name: "items", Values: items},
	}
}

type Compare struct {
	BaseNode

	Expr Expression
	Ops  []*Operand
}

func NewCompare(token token.Token, expr Expression, ops []*Operand) *Compare {
	return &Compare{
		BaseNode: NewBaseNode(token),
		Expr:     expr,
		Ops:      ops,
	}
}

func (c *Compare) expressionNode() {}

func (c *Compare) String() string {
	ops := make([]Node, len(c.Ops))
	for i := range c.Ops {
		ops[i] = c.Ops[i]
	}
	return fmt.Sprintf("nodes.Compare(expr=%s, ops=%s)", c.Expr, reprNodeList(ops))
}

func (c *Compare) Dumps(indent int) string {
	ops := make([]Node, len(c.Ops))
	for i := range c.Ops {
		ops[i] = c.Ops[i]
	}

	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Compare(")
	fmt.Fprintf(sb, "  expr=%s,\n", c.Expr)
	sb.WriteIndent()
	sb.WriteNodeList("ops", ops)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (c *Compare) ChildNodes() []ChildNode {
	ops := make([]Node, len(c.Ops))
	for i := range c.Ops {
		ops[i] = c.Ops[i]
	}
	return []ChildNode{
		{Name: "expr", Value: c.Expr},
		{Name: "ops", Values: ops},
	}
}

func stringBinExpr(name string, be *BinExpr) string {
	return fmt.Sprintf("nodes.%s(left=%s, right=%s)", name, be.Left, be.Right)
}

func dumpsBinExpr(name string, be *BinExpr, indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent(fmt.Sprintf("nodes.%s(", name))
	fmt.Fprintf(sb, "  left=%s,\n", be.Left)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  right=%s,\n", be.Right)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func stringUnaryExpr(name string, ue *UnaryExpr) string {
	return fmt.Sprintf("nodes.%s(node=%s)", name, ue.Node)
}

func dumpsUnaryExpr(name string, ue *UnaryExpr, indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent(fmt.Sprintf("nodes.%s(", name))
	fmt.Fprintf(sb, "  node=%s,\n", ue.Node)
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
			Left:     left,
			Operator: "*",
			Right:    right,
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
			Left:     left,
			Operator: "/",
			Right:    right,
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
			Left:     left,
			Operator: "//",
			Right:    right,
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
			Left:     left,
			Operator: "+",
			Right:    right,
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
			Left:     left,
			Operator: "-",
			Right:    right,
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
			Left:     left,
			Operator: "%",
			Right:    right,
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
			Left:     left,
			Operator: "**",
			Right:    right,
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
			Left:     left,
			Operator: "and",
			Right:    right,
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
			Left:     left,
			Operator: "or",
			Right:    right,
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
			Operator: "not",
			Node:     node,
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
			Operator: "-",
			Node:     node,
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
			Operator: "+",
			Node:     node,
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

	Expr Expression
}

func NewMarkSafe(token token.Token, expr Expression) *MarkSafe {
	return &MarkSafe{
		BaseNode: NewBaseNode(token),
		Expr:     expr,
	}
}

func (ms *MarkSafe) expressionNode() {}

func (ms *MarkSafe) String() string {
	return fmt.Sprintf("nodes.MarkSafe(expr=%s)", ms.Expr)
}

func (ms *MarkSafe) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.MarkSafe(")
	fmt.Fprintf(sb, "  expr=%s,\n", ms.Expr)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (ms *MarkSafe) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "expr", Value: ms.Expr},
	}
}

type MarkSafeIfAutoescape struct {
	BaseNode

	Expr Expression
}

func NewMarkSafeIfAutoescape(token token.Token, expr Expression) *MarkSafeIfAutoescape {
	return &MarkSafeIfAutoescape{
		BaseNode: NewBaseNode(token),
		Expr:     expr,
	}
}

func (ms *MarkSafeIfAutoescape) expressionNode() {}

func (ms *MarkSafeIfAutoescape) String() string {
	return fmt.Sprintf("nodes.MarkSafeIfAutoescape(expr=%s)", ms.Expr)
}

func (ms *MarkSafeIfAutoescape) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.MarkSafeIfAutoescape(")
	fmt.Fprintf(sb, "  expr=%s,\n", ms.Expr)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (ms *MarkSafeIfAutoescape) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "expr", Value: ms.Expr},
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
