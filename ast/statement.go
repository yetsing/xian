package ast

import (
	"fmt"

	"github.com/yetsing/xian/token"
)

type Output struct {
	BaseNode

	Nodes []Expression
}

func NewOutput(tk token.Token, nodes []Expression) *Output {
	o := &Output{Nodes: nodes}
	o.BaseNode = NewBaseNode(tk)
	return o
}

func (o *Output) statementNode() {}

func (o *Output) String() string {
	return fmt.Sprintf("nodes.Output(nodes=%s)", reprExpressionList(o.Nodes))
}

func (o *Output) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Output(")
	sb.WriteExpressionList("nodes", o.Nodes)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (o *Output) ChildNodes() []ChildNode {
	nodes := make([]Node, len(o.Nodes))
	for i := range o.Nodes {
		nodes[i] = o.Nodes[i]
	}

	return []ChildNode{
		{Name: "nodes", Values: nodes},
	}
}

type Extends struct {
	BaseNode

	Template Expression
}

func NewExtends(token token.Token, template Expression) *Extends {
	return &Extends{
		BaseNode: NewBaseNode(token),
		Template: template,
	}
}

func (e *Extends) statementNode() {}

func (e *Extends) String() string {
	return fmt.Sprintf("nodes.Extends(template=%s)", e.Template.String())
}

func (e *Extends) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Extends(")
	fmt.Fprintf(sb, "  template=%s,\n", e.Template.Dumps(indent+2))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (e *Extends) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "template", Value: e.Template},
	}
}

type For struct {
	BaseNode

	target    Node
	iter      Node
	body      []Node
	else_     []Node
	test      Node
	recursive bool
}

func NewFor(token token.Token, target Node, iter Node, body []Node, else_ []Node, test Node, recursive bool) *For {
	return &For{
		BaseNode:  NewBaseNode(token),
		target:    target,
		iter:      iter,
		body:      body,
		else_:     else_,
		test:      test,
		recursive: recursive,
	}
}

func (f *For) statementNode() {}

func (f *For) String() string {
	return fmt.Sprintf("nodes.For(target=%s, iter=%s, body=%s, else_=%s, test=%s, recursive=%t)",
		f.target.String(), f.iter.String(), reprNodeList(f.body), reprNodeList(f.else_), f.test.String(), f.recursive)
}

func (f *For) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.For(")
	fmt.Fprintf(sb, "  target=%s,\n", f.target.Dumps(indent+2))
	sb.WriteIndent()
	fmt.Fprintf(sb, "  iter=%s,\n", f.iter.Dumps(indent+2))
	sb.WriteIndent()
	sb.WriteNodeList("body", f.body)
	sb.WriteIndent()
	sb.WriteNodeList("else_", f.else_)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  test=%s,\n", f.test.Dumps(indent+2))
	sb.WriteIndent()
	fmt.Fprintf(sb, "  recursive=%t,\n", f.recursive)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (f *For) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "target", Value: f.target},
		{Name: "iter", Value: f.iter},
		{Name: "body", Values: f.body},
		{Name: "else_", Values: f.else_},
		{Name: "test", Value: f.test},
	}
}

type If struct {
	BaseNode

	test  Node
	body  []Node
	elif_ []*If
	else_ []Node
}

func NewIf(token token.Token, test Node, body []Node, elif_ []*If, else_ []Node) *If {
	return &If{
		BaseNode: NewBaseNode(token),
		test:     test,
		body:     body,
		elif_:    elif_,
		else_:    else_,
	}
}

func (i *If) statementNode() {}

func (i *If) String() string {
	elifNodes := make([]Node, len(i.elif_))
	for j := range i.elif_ {
		elifNodes[j] = i.elif_[j]
	}
	return fmt.Sprintf("nodes.If(test=%s, body=%s, elif_=%s, else_=%s)",
		i.test.String(), reprNodeList(i.body), reprNodeList(elifNodes), reprNodeList(i.else_))
}

func (i *If) Dumps(indent int) string {
	elifNodes := make([]Node, len(i.elif_))
	for j := range i.elif_ {
		elifNodes[j] = i.elif_[j]
	}

	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.If(")
	fmt.Fprintf(sb, "  test=%s,\n", i.test.Dumps(indent+2))
	sb.WriteIndent()
	sb.WriteNodeList("body", i.body)
	sb.WriteIndent()
	sb.WriteNodeList("elif_", elifNodes)
	sb.WriteIndent()
	sb.WriteNodeList("else_", i.else_)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (i *If) ChildNodes() []ChildNode {
	elifNodes := make([]Node, len(i.elif_))
	for j := range i.elif_ {
		elifNodes[j] = i.elif_[j]
	}

	return []ChildNode{
		{Name: "test", Value: i.test},
		{Name: "body", Values: i.body},
		{Name: "elif_", Values: elifNodes},
		{Name: "else_", Values: i.else_},
	}
}

type Macro struct {
	BaseNode

	name     string
	args     []*Name
	defaults []Expression
	body     []Node
}

func NewMacro(token token.Token, name string, args []*Name, defaults []Expression, body []Node) *Macro {
	return &Macro{
		BaseNode: NewBaseNode(token),
		name:     name,
		args:     args,
		defaults: defaults,
		body:     body,
	}
}

func (m *Macro) statementNode() {}

func (m *Macro) String() string {
	args := make([]Node, len(m.args))
	for i := range m.args {
		args[i] = m.args[i]
	}
	return fmt.Sprintf("nodes.Macro(name=%s, args=%s, defaults=%s, body=%s)",
		reprString(m.name), reprNodeList(args), reprExpressionList(m.defaults), reprNodeList(m.body))
}

func (m *Macro) Dumps(indent int) string {
	args := make([]Node, len(m.args))
	for i := range m.args {
		args[i] = m.args[i]
	}

	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Macro(")
	fmt.Fprintf(sb, "  name=%s,\n", reprString(m.name))
	sb.WriteIndent()
	sb.WriteNodeList("args", args)
	sb.WriteIndent()
	sb.WriteExpressionList("defaults", m.defaults)
	sb.WriteIndent()
	sb.WriteNodeList("body", m.body)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (m *Macro) ChildNodes() []ChildNode {
	args := make([]Node, len(m.args))
	for i := range m.args {
		args[i] = m.args[i]
	}
	defaults := make([]Node, len(m.defaults))
	for i := range m.defaults {
		defaults[i] = m.defaults[i]
	}

	return []ChildNode{
		{Name: "args", Values: args},
		{Name: "defaults", Values: defaults},
		{Name: "body", Values: m.body},
	}
}

type CallBlock struct {
	BaseNode

	call     *Call
	args     []*Name
	defaults []Expression
	body     []Node
}

func NewCallBlock(token token.Token, call *Call, args []*Name, defaults []Expression, body []Node) *CallBlock {
	return &CallBlock{
		BaseNode: NewBaseNode(token),
		call:     call,
		args:     args,
		defaults: defaults,
		body:     body,
	}
}

func (cb *CallBlock) statementNode() {}

func (cb *CallBlock) String() string {
	args := make([]Node, len(cb.args))
	for i := range cb.args {
		args[i] = cb.args[i]
	}
	return fmt.Sprintf("nodes.CallBlock(call=%s, args=%s, defaults=%s, body=%s)",
		cb.call.String(), reprNodeList(args), reprExpressionList(cb.defaults), reprNodeList(cb.body))
}

func (cb *CallBlock) Dumps(indent int) string {
	args := make([]Node, len(cb.args))
	for i := range cb.args {
		args[i] = cb.args[i]
	}

	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.CallBlock(")
	fmt.Fprintf(sb, "  call=%s,\n", cb.call.Dumps(indent+2))
	sb.WriteIndent()
	sb.WriteNodeList("args", args)
	sb.WriteIndent()
	sb.WriteExpressionList("defaults", cb.defaults)
	sb.WriteIndent()
	sb.WriteNodeList("body", cb.body)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (cb *CallBlock) ChildNodes() []ChildNode {
	args := make([]Node, len(cb.args))
	for i := range cb.args {
		args[i] = cb.args[i]
	}
	defaults := make([]Node, len(cb.defaults))
	for i := range cb.defaults {
		defaults[i] = cb.defaults[i]
	}

	return []ChildNode{
		{Name: "call", Value: cb.call},
		{Name: "args", Values: args},
		{Name: "defaults", Values: defaults},
		{Name: "body", Values: cb.body},
	}
}

type FilterBlock struct {
	BaseNode

	body   []Node
	filter *Filter
}

func NewFilterBlock(token token.Token, body []Node, filter *Filter) *FilterBlock {
	return &FilterBlock{
		BaseNode: NewBaseNode(token),
		body:     body,
		filter:   filter,
	}
}

func (fb *FilterBlock) statementNode() {}

func (fb *FilterBlock) String() string {
	return fmt.Sprintf("nodes.FilterBlock(body=%s, filter=%s)",
		reprNodeList(fb.body), fb.filter.String())
}

func (fb *FilterBlock) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.FilterBlock(")
	sb.WriteIndent()
	sb.WriteNodeList("body", fb.body)
	fmt.Fprintf(sb, "  filter=%s,\n", fb.filter.Dumps(indent+2))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (fb *FilterBlock) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "body", Values: fb.body},
		{Name: "filter", Value: fb.filter},
	}
}

type With struct {
	BaseNode

	targets []Expression
	values  []Expression
	body    []Node
}

func NewWith(token token.Token, targets []Expression, values []Expression, body []Node) *With {
	return &With{
		BaseNode: NewBaseNode(token),
		targets:  targets,
		values:   values,
		body:     body,
	}
}

func (w *With) statementNode() {}

func (w *With) String() string {
	return fmt.Sprintf("nodes.With(targets=%s, values=%s, body=%s)",
		reprExpressionList(w.targets), reprExpressionList(w.values), reprNodeList(w.body))
}

func (w *With) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.With(")
	sb.WriteExpressionList("targets", w.targets)
	sb.WriteIndent()
	sb.WriteExpressionList("values", w.values)
	sb.WriteIndent()
	sb.WriteNodeList("body", w.body)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (w *With) ChildNodes() []ChildNode {
	targets := make([]Node, len(w.targets))
	for i := range w.targets {
		targets[i] = w.targets[i]
	}
	values := make([]Node, len(w.values))
	for i := range w.values {
		values[i] = w.values[i]
	}

	return []ChildNode{
		{Name: "targets", Values: targets},
		{Name: "values", Values: values},
		{Name: "body", Values: w.body},
	}
}

type Block struct {
	BaseNode

	name     string
	body     []Node
	scoped   bool
	required bool
}

func NewBlock(token token.Token, name string, body []Node, scoped bool, required bool) *Block {
	return &Block{
		BaseNode: NewBaseNode(token),
		name:     name,
		body:     body,
		scoped:   scoped,
		required: required,
	}
}

func (b *Block) statementNode() {}

func (b *Block) String() string {
	return fmt.Sprintf("nodes.Block(name=%s, body=%s, scoped=%t, required=%t)",
		reprString(b.name), reprNodeList(b.body), b.scoped, b.required)
}

func (b *Block) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Block(")
	fmt.Fprintf(sb, "  name=%s,\n", reprString(b.name))
	sb.WriteIndent()
	sb.WriteNodeList("body", b.body)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  scoped=%t,\n", b.scoped)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  required=%t,\n", b.required)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (b *Block) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "body", Values: b.body},
	}
}

type Include struct {
	BaseNode

	template      Expression
	withContext   bool
	ignoreMissing bool
}

func NewInclude(token token.Token, template Expression, withContext bool, ignoreMissing bool) *Include {
	return &Include{
		BaseNode:      NewBaseNode(token),
		template:      template,
		withContext:   withContext,
		ignoreMissing: ignoreMissing,
	}
}

func (i *Include) statementNode() {}

func (i *Include) String() string {
	return fmt.Sprintf("nodes.Include(template=%s, with_context=%t, ignore_missing=%t)",
		i.template.String(), i.withContext, i.ignoreMissing)
}

func (i *Include) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Include(")
	fmt.Fprintf(sb, "  template=%s,\n", i.template.Dumps(indent+2))
	sb.WriteIndent()
	fmt.Fprintf(sb, "  with_context=%t,\n", i.withContext)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  ignore_missing=%t,\n", i.ignoreMissing)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (i *Include) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "template", Value: i.template},
	}
}

type Import struct {
	BaseNode

	template    Expression
	target      string
	withContext bool
}

func NewImport(token token.Token, template Expression, target string, withContext bool) *Import {
	return &Import{
		BaseNode:    NewBaseNode(token),
		template:    template,
		target:      target,
		withContext: withContext,
	}
}

func (im *Import) statementNode() {}

func (im *Import) String() string {
	return fmt.Sprintf("nodes.Import(template=%s, target=%s, with_context=%t)",
		im.template.String(), reprString(im.target), im.withContext)
}

func (im *Import) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Import(")
	fmt.Fprintf(sb, "  template=%s,\n", im.template.Dumps(indent+2))
	sb.WriteIndent()
	fmt.Fprintf(sb, "  target=%s,\n", reprString(im.target))
	sb.WriteIndent()
	fmt.Fprintf(sb, "  with_context=%t,\n", im.withContext)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (im *Import) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "template", Value: im.template},
	}
}

type FromImport struct {
	BaseNode

	template    Expression
	names       []string
	withContext bool
}

func NewFromImport(token token.Token, template Expression, names []string, withContext bool) *FromImport {
	return &FromImport{
		BaseNode:    NewBaseNode(token),
		template:    template,
		names:       names,
		withContext: withContext,
	}
}

func (fi *FromImport) statementNode() {}

func (fi *FromImport) String() string {
	return fmt.Sprintf("nodes.FromImport(template=%s, names=%s, with_context=%t)",
		fi.template.String(), reprStringList(fi.names), fi.withContext)
}

func (fi *FromImport) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.FromImport(")
	fmt.Fprintf(sb, "  template=%s,\n", fi.template.Dumps(indent+2))
	sb.WriteIndent()
	sb.WriteStringList("names", fi.names)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  with_context=%t,\n", fi.withContext)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (fi *FromImport) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "template", Value: fi.template},
	}
}

type ExprStmt struct {
	BaseNode

	node Node
}

func NewExprStmt(token token.Token, node Node) *ExprStmt {
	return &ExprStmt{
		BaseNode: NewBaseNode(token),
		node:     node,
	}
}

func (es *ExprStmt) statementNode() {}

func (es *ExprStmt) String() string {
	return fmt.Sprintf("nodes.ExprStmt(node=%s)", es.node.String())
}

func (es *ExprStmt) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.ExprStmt(")
	fmt.Fprintf(sb, "  node=%s,\n", es.node.Dumps(indent+2))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (es *ExprStmt) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "node", Value: es.node},
	}
}

type Assign struct {
	BaseNode

	target Expression
	node   Node
}

func NewAssign(token token.Token, target Expression, node Node) *Assign {
	return &Assign{
		BaseNode: NewBaseNode(token),
		target:   target,
		node:     node,
	}
}

func (a *Assign) statementNode() {}

func (a *Assign) String() string {
	return fmt.Sprintf("nodes.Assign(target=%s, node=%s)", a.target.String(), a.node.String())
}

func (a *Assign) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Assign(")
	fmt.Fprintf(sb, "  target=%s,\n", a.target.Dumps(indent+2))
	sb.WriteIndent()
	fmt.Fprintf(sb, "  node=%s,\n", a.node.Dumps(indent+2))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (a *Assign) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "target", Value: a.target},
		{Name: "node", Value: a.node},
	}
}

type AssignBlock struct {
	BaseNode

	target Expression
	filter *Filter
	body   []Node
}

func NewAssignBlock(token token.Token, target Expression, filter *Filter, body []Node) *AssignBlock {
	return &AssignBlock{
		BaseNode: NewBaseNode(token),
		target:   target,
		filter:   filter,
		body:     body,
	}
}

func (ab *AssignBlock) statementNode() {}

func (ab *AssignBlock) String() string {
	return fmt.Sprintf("nodes.AssignBlock(target=%s, filter=%s, body=%s)",
		ab.target.String(), ab.filter.String(), reprNodeList(ab.body))
}

func (ab *AssignBlock) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.AssignBlock(")
	fmt.Fprintf(sb, "  target=%s,\n", ab.target.Dumps(indent+2))
	sb.WriteIndent()
	fmt.Fprintf(sb, "  filter=%s,\n", ab.filter.Dumps(indent+2))
	sb.WriteIndent()
	sb.WriteNodeList("body", ab.body)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (ab *AssignBlock) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "target", Value: ab.target},
		{Name: "filter", Value: ab.filter},
		{Name: "body", Values: ab.body},
	}
}

type Continue struct {
	BaseNode
}

func NewContinue(token token.Token) *Continue {
	return &Continue{
		BaseNode: NewBaseNode(token),
	}
}

func (c *Continue) statementNode() {}

func (c *Continue) String() string {
	return "nodes.Continue()"
}

func (c *Continue) Dumps(indent int) string {
	return "nodes.Continue()"
}

func (c *Continue) ChildNodes() []ChildNode {
	return nil
}

type Break struct {
	BaseNode
}

func NewBreak(token token.Token) *Break {
	return &Break{
		BaseNode: NewBaseNode(token),
	}
}

func (b *Break) statementNode() {}

func (b *Break) String() string {
	return "nodes.Break()"
}

func (b *Break) Dumps(indent int) string {
	return "nodes.Break()"
}

func (b *Break) ChildNodes() []ChildNode {
	return nil
}

type Scope struct {
	BaseNode

	body []Node
}

func NewScope(token token.Token, body []Node) *Scope {
	return &Scope{
		BaseNode: NewBaseNode(token),
		body:     body,
	}
}

func (s *Scope) statementNode() {}

func (s *Scope) String() string {
	return fmt.Sprintf("nodes.Scope(body=%s)", reprNodeList(s.body))
}

func (s *Scope) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Scope(")
	sb.WriteNodeList("body", s.body)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (s *Scope) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "body", Values: s.body},
	}
}

type OverlayScopy struct {
	BaseNode

	context Expression
	body    []Node
}

func NewOverlayScopy(token token.Token, context Expression, body []Node) *OverlayScopy {
	return &OverlayScopy{
		BaseNode: NewBaseNode(token),
		context:  context,
		body:     body,
	}
}

func (os *OverlayScopy) statementNode() {}

func (os *OverlayScopy) String() string {
	return fmt.Sprintf("nodes.OverlayScopy(context=%s, body=%s)", os.context.String(), reprNodeList(os.body))
}

func (os *OverlayScopy) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.OverlayScopy(")
	fmt.Fprintf(sb, "  context=%s,\n", os.context.Dumps(indent+2))
	sb.WriteIndent()
	sb.WriteNodeList("body", os.body)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (os *OverlayScopy) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "context", Value: os.context},
		{Name: "body", Values: os.body},
	}
}

type EvalContextModifier struct {
	BaseNode

	options []*Keyword
}

func NewEvalContextModifier(token token.Token, options []*Keyword) *EvalContextModifier {
	return &EvalContextModifier{
		BaseNode: NewBaseNode(token),
		options:  options,
	}
}

func (ecm *EvalContextModifier) statementNode() {}

func (ecm *EvalContextModifier) String() string {
	options := make([]Node, len(ecm.options))
	for i := range ecm.options {
		options[i] = ecm.options[i]
	}
	return fmt.Sprintf("nodes.EvalContextModifier(options=%s)", reprNodeList(options))
}

func (ecm *EvalContextModifier) Dumps(indent int) string {
	options := make([]Node, len(ecm.options))
	for i := range ecm.options {
		options[i] = ecm.options[i]
	}

	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.EvalContextModifier(")
	sb.WriteNodeList("options", options)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (ecm *EvalContextModifier) ChildNodes() []ChildNode {
	options := make([]Node, len(ecm.options))
	for i := range ecm.options {
		options[i] = ecm.options[i]
	}

	return []ChildNode{
		{Name: "options", Values: options},
	}
}

type ScopedEvalContextModifier struct {
	EvalContextModifier

	body []Node
}

func NewScopedEvalContextModifier(token token.Token, options []*Keyword, body []Node) *ScopedEvalContextModifier {
	return &ScopedEvalContextModifier{
		EvalContextModifier: *NewEvalContextModifier(token, options),
		body:                body,
	}
}

func (sem *ScopedEvalContextModifier) statementNode() {}

func (sem *ScopedEvalContextModifier) String() string {
	options := make([]Node, len(sem.options))
	for i := range sem.options {
		options[i] = sem.options[i]
	}
	return fmt.Sprintf("nodes.ScopedEvalContextModifier(options=%s, body=%s)", reprNodeList(options), reprNodeList(sem.body))
}

func (sem *ScopedEvalContextModifier) Dumps(indent int) string {
	options := make([]Node, len(sem.options))
	for i := range sem.options {
		options[i] = sem.options[i]
	}

	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.ScopedEvalContextModifier(")
	sb.WriteNodeList("options", options)
	sb.WriteIndent()
	sb.WriteNodeList("body", sem.body)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (sem *ScopedEvalContextModifier) ChildNodes() []ChildNode {
	childnodes := sem.EvalContextModifier.ChildNodes()
	childnodes = append(childnodes, ChildNode{Name: "body", Values: sem.body})

	return childnodes
}
