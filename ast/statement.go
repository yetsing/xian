package ast

import (
	"fmt"

	"github.com/yetsing/xian/token"
)

// All statement nodes implement this
type Statement interface {
	Node
	statementNode()
}

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

	Target    Node
	Iter      Node
	Body      []Node
	Else_     []Node
	Test      Node
	Recursive bool
}

func NewFor(token token.Token, target Node, iter Node, body []Node, else_ []Node, test Node, recursive bool) *For {
	return &For{
		BaseNode:  NewBaseNode(token),
		Target:    target,
		Iter:      iter,
		Body:      body,
		Else_:     else_,
		Test:      test,
		Recursive: recursive,
	}
}

func (f *For) statementNode() {}

func (f *For) String() string {
	return fmt.Sprintf("nodes.For(target=%s, iter=%s, body=%s, else_=%s, test=%s, recursive=%t)",
		f.Target.String(), f.Iter.String(), reprNodeList(f.Body), reprNodeList(f.Else_), f.Test.String(), f.Recursive)
}

func (f *For) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.For(")
	fmt.Fprintf(sb, "  target=%s,\n", f.Target.Dumps(indent+2))
	sb.WriteIndent()
	fmt.Fprintf(sb, "  iter=%s,\n", f.Iter.Dumps(indent+2))
	sb.WriteIndent()
	sb.WriteNodeList("body", f.Body)
	sb.WriteIndent()
	sb.WriteNodeList("else_", f.Else_)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  test=%s,\n", f.Test.Dumps(indent+2))
	sb.WriteIndent()
	fmt.Fprintf(sb, "  recursive=%t,\n", f.Recursive)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (f *For) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "target", Value: f.Target},
		{Name: "iter", Value: f.Iter},
		{Name: "body", Values: f.Body},
		{Name: "else_", Values: f.Else_},
		{Name: "test", Value: f.Test},
	}
}

type If struct {
	BaseNode

	Test  Node
	Body  []Node
	Elif_ []*If
	Else_ []Node
}

func NewIf(token token.Token, test Node, body []Node, elif_ []*If, else_ []Node) *If {
	return &If{
		BaseNode: NewBaseNode(token),
		Test:     test,
		Body:     body,
		Elif_:    elif_,
		Else_:    else_,
	}
}

func (i *If) statementNode() {}

func (i *If) String() string {
	elifNodes := make([]Node, len(i.Elif_))
	for j := range i.Elif_ {
		elifNodes[j] = i.Elif_[j]
	}
	return fmt.Sprintf("nodes.If(test=%s, body=%s, elif_=%s, else_=%s)",
		i.Test.String(), reprNodeList(i.Body), reprNodeList(elifNodes), reprNodeList(i.Else_))
}

func (i *If) Dumps(indent int) string {
	elifNodes := make([]Node, len(i.Elif_))
	for j := range i.Elif_ {
		elifNodes[j] = i.Elif_[j]
	}

	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.If(")
	fmt.Fprintf(sb, "  test=%s,\n", i.Test.Dumps(indent+2))
	sb.WriteIndent()
	sb.WriteNodeList("body", i.Body)
	sb.WriteIndent()
	sb.WriteNodeList("elif_", elifNodes)
	sb.WriteIndent()
	sb.WriteNodeList("else_", i.Else_)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (i *If) ChildNodes() []ChildNode {
	elifNodes := make([]Node, len(i.Elif_))
	for j := range i.Elif_ {
		elifNodes[j] = i.Elif_[j]
	}

	return []ChildNode{
		{Name: "test", Value: i.Test},
		{Name: "body", Values: i.Body},
		{Name: "elif_", Values: elifNodes},
		{Name: "else_", Values: i.Else_},
	}
}

type Macro struct {
	BaseNode

	Name     string
	Args     []*Name
	Defaults []Expression
	Body     []Node
}

func NewMacro(token token.Token, name string, args []*Name, defaults []Expression, body []Node) *Macro {
	return &Macro{
		BaseNode: NewBaseNode(token),
		Name:     name,
		Args:     args,
		Defaults: defaults,
		Body:     body,
	}
}

func (m *Macro) statementNode() {}

func (m *Macro) String() string {
	args := make([]Node, len(m.Args))
	for i := range m.Args {
		args[i] = m.Args[i]
	}
	return fmt.Sprintf("nodes.Macro(name=%s, args=%s, defaults=%s, body=%s)",
		reprString(m.Name), reprNodeList(args), reprExpressionList(m.Defaults), reprNodeList(m.Body))
}

func (m *Macro) Dumps(indent int) string {
	args := make([]Node, len(m.Args))
	for i := range m.Args {
		args[i] = m.Args[i]
	}

	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Macro(")
	fmt.Fprintf(sb, "  name=%s,\n", reprString(m.Name))
	sb.WriteIndent()
	sb.WriteNodeList("args", args)
	sb.WriteIndent()
	sb.WriteExpressionList("defaults", m.Defaults)
	sb.WriteIndent()
	sb.WriteNodeList("body", m.Body)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (m *Macro) ChildNodes() []ChildNode {
	args := make([]Node, len(m.Args))
	for i := range m.Args {
		args[i] = m.Args[i]
	}
	defaults := make([]Node, len(m.Defaults))
	for i := range m.Defaults {
		defaults[i] = m.Defaults[i]
	}

	return []ChildNode{
		{Name: "args", Values: args},
		{Name: "defaults", Values: defaults},
		{Name: "body", Values: m.Body},
	}
}

type CallBlock struct {
	BaseNode

	Call     *Call
	Args     []*Name
	Defaults []Expression
	Body     []Node
}

func NewCallBlock(token token.Token, call *Call, args []*Name, defaults []Expression, body []Node) *CallBlock {
	return &CallBlock{
		BaseNode: NewBaseNode(token),
		Call:     call,
		Args:     args,
		Defaults: defaults,
		Body:     body,
	}
}

func (cb *CallBlock) statementNode() {}

func (cb *CallBlock) String() string {
	args := make([]Node, len(cb.Args))
	for i := range cb.Args {
		args[i] = cb.Args[i]
	}
	return fmt.Sprintf("nodes.CallBlock(call=%s, args=%s, defaults=%s, body=%s)",
		cb.Call.String(), reprNodeList(args), reprExpressionList(cb.Defaults), reprNodeList(cb.Body))
}

func (cb *CallBlock) Dumps(indent int) string {
	args := make([]Node, len(cb.Args))
	for i := range cb.Args {
		args[i] = cb.Args[i]
	}

	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.CallBlock(")
	fmt.Fprintf(sb, "  call=%s,\n", cb.Call.Dumps(indent+2))
	sb.WriteIndent()
	sb.WriteNodeList("args", args)
	sb.WriteIndent()
	sb.WriteExpressionList("defaults", cb.Defaults)
	sb.WriteIndent()
	sb.WriteNodeList("body", cb.Body)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (cb *CallBlock) ChildNodes() []ChildNode {
	args := make([]Node, len(cb.Args))
	for i := range cb.Args {
		args[i] = cb.Args[i]
	}
	defaults := make([]Node, len(cb.Defaults))
	for i := range cb.Defaults {
		defaults[i] = cb.Defaults[i]
	}

	return []ChildNode{
		{Name: "call", Value: cb.Call},
		{Name: "args", Values: args},
		{Name: "defaults", Values: defaults},
		{Name: "body", Values: cb.Body},
	}
}

type FilterBlock struct {
	BaseNode

	Body   []Node
	Filter *Filter
}

func NewFilterBlock(token token.Token, body []Node, filter *Filter) *FilterBlock {
	return &FilterBlock{
		BaseNode: NewBaseNode(token),
		Body:     body,
		Filter:   filter,
	}
}

func (fb *FilterBlock) statementNode() {}

func (fb *FilterBlock) String() string {
	return fmt.Sprintf("nodes.FilterBlock(body=%s, filter=%s)",
		reprNodeList(fb.Body), fb.Filter.String())
}

func (fb *FilterBlock) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.FilterBlock(")
	sb.WriteIndent()
	sb.WriteNodeList("body", fb.Body)
	fmt.Fprintf(sb, "  filter=%s,\n", fb.Filter.Dumps(indent+2))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (fb *FilterBlock) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "body", Values: fb.Body},
		{Name: "filter", Value: fb.Filter},
	}
}

type With struct {
	BaseNode

	Targets []Expression
	Values  []Expression
	Body    []Node
}

func NewWith(token token.Token, targets []Expression, values []Expression, body []Node) *With {
	return &With{
		BaseNode: NewBaseNode(token),
		Targets:  targets,
		Values:   values,
		Body:     body,
	}
}

func (w *With) statementNode() {}

func (w *With) String() string {
	return fmt.Sprintf("nodes.With(targets=%s, values=%s, body=%s)",
		reprExpressionList(w.Targets), reprExpressionList(w.Values), reprNodeList(w.Body))
}

func (w *With) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.With(")
	sb.WriteExpressionList("targets", w.Targets)
	sb.WriteIndent()
	sb.WriteExpressionList("values", w.Values)
	sb.WriteIndent()
	sb.WriteNodeList("body", w.Body)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (w *With) ChildNodes() []ChildNode {
	targets := make([]Node, len(w.Targets))
	for i := range w.Targets {
		targets[i] = w.Targets[i]
	}
	values := make([]Node, len(w.Values))
	for i := range w.Values {
		values[i] = w.Values[i]
	}

	return []ChildNode{
		{Name: "targets", Values: targets},
		{Name: "values", Values: values},
		{Name: "body", Values: w.Body},
	}
}

type Block struct {
	BaseNode

	Name     string
	Body     []Node
	Scoped   bool
	Required bool
}

func NewBlock(token token.Token, name string, body []Node, scoped bool, required bool) *Block {
	return &Block{
		BaseNode: NewBaseNode(token),
		Name:     name,
		Body:     body,
		Scoped:   scoped,
		Required: required,
	}
}

func (b *Block) statementNode() {}

func (b *Block) String() string {
	return fmt.Sprintf("nodes.Block(name=%s, body=%s, scoped=%t, required=%t)",
		reprString(b.Name), reprNodeList(b.Body), b.Scoped, b.Required)
}

func (b *Block) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Block(")
	fmt.Fprintf(sb, "  name=%s,\n", reprString(b.Name))
	sb.WriteIndent()
	sb.WriteNodeList("body", b.Body)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  scoped=%t,\n", b.Scoped)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  required=%t,\n", b.Required)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (b *Block) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "body", Values: b.Body},
	}
}

type Include struct {
	BaseNode

	Template      Expression
	WithContext   bool
	IgnoreMissing bool
}

func NewInclude(token token.Token, template Expression, withContext bool, ignoreMissing bool) *Include {
	return &Include{
		BaseNode:      NewBaseNode(token),
		Template:      template,
		WithContext:   withContext,
		IgnoreMissing: ignoreMissing,
	}
}

func (i *Include) statementNode() {}

func (i *Include) String() string {
	return fmt.Sprintf("nodes.Include(template=%s, with_context=%t, ignore_missing=%t)",
		i.Template.String(), i.WithContext, i.IgnoreMissing)
}

func (i *Include) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Include(")
	fmt.Fprintf(sb, "  template=%s,\n", i.Template.Dumps(indent+2))
	sb.WriteIndent()
	fmt.Fprintf(sb, "  with_context=%t,\n", i.WithContext)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  ignore_missing=%t,\n", i.IgnoreMissing)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (i *Include) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "template", Value: i.Template},
	}
}

type Import struct {
	BaseNode

	Template    Expression
	Target      string
	WithContext bool
}

func NewImport(token token.Token, template Expression, target string, withContext bool) *Import {
	return &Import{
		BaseNode:    NewBaseNode(token),
		Template:    template,
		Target:      target,
		WithContext: withContext,
	}
}

func (im *Import) statementNode() {}

func (im *Import) String() string {
	return fmt.Sprintf("nodes.Import(template=%s, target=%s, with_context=%t)",
		im.Template.String(), reprString(im.Target), im.WithContext)
}

func (im *Import) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Import(")
	fmt.Fprintf(sb, "  template=%s,\n", im.Template.Dumps(indent+2))
	sb.WriteIndent()
	fmt.Fprintf(sb, "  target=%s,\n", reprString(im.Target))
	sb.WriteIndent()
	fmt.Fprintf(sb, "  with_context=%t,\n", im.WithContext)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (im *Import) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "template", Value: im.Template},
	}
}

type FromImport struct {
	BaseNode

	Template    Expression
	Names       []string
	Aliases     []string
	WithContext bool
}

func NewFromImport(token token.Token, template Expression, names []string, aliases []string, withContext bool) *FromImport {
	return &FromImport{
		BaseNode:    NewBaseNode(token),
		Template:    template,
		Names:       names,
		Aliases:     aliases,
		WithContext: withContext,
	}
}

func (fi *FromImport) statementNode() {}

func (fi *FromImport) String() string {
	return fmt.Sprintf("nodes.FromImport(template=%s, names=%s, with_context=%t)",
		fi.Template.String(), reprStringList(fi.Names), fi.WithContext)
}

func (fi *FromImport) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.FromImport(")
	fmt.Fprintf(sb, "  template=%s,\n", fi.Template.Dumps(indent+2))
	sb.WriteIndent()
	sb.WriteStringList("names", fi.Names)
	sb.WriteIndent()
	fmt.Fprintf(sb, "  with_context=%t,\n", fi.WithContext)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (fi *FromImport) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "template", Value: fi.Template},
	}
}

type ExprStmt struct {
	BaseNode

	Node Node
}

func NewExprStmt(token token.Token, node Node) *ExprStmt {
	return &ExprStmt{
		BaseNode: NewBaseNode(token),
		Node:     node,
	}
}

func (es *ExprStmt) statementNode() {}

func (es *ExprStmt) String() string {
	return fmt.Sprintf("nodes.ExprStmt(node=%s)", es.Node.String())
}

func (es *ExprStmt) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.ExprStmt(")
	fmt.Fprintf(sb, "  node=%s,\n", es.Node.Dumps(indent+2))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (es *ExprStmt) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "node", Value: es.Node},
	}
}

type Assign struct {
	BaseNode

	Target Expression
	Node   Node
}

func NewAssign(token token.Token, target Expression, node Node) *Assign {
	return &Assign{
		BaseNode: NewBaseNode(token),
		Target:   target,
		Node:     node,
	}
}

func (a *Assign) statementNode() {}

func (a *Assign) String() string {
	return fmt.Sprintf("nodes.Assign(target=%s, node=%s)", a.Target.String(), a.Node.String())
}

func (a *Assign) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Assign(")
	fmt.Fprintf(sb, "  target=%s,\n", a.Target.Dumps(indent+2))
	sb.WriteIndent()
	fmt.Fprintf(sb, "  node=%s,\n", a.Node.Dumps(indent+2))
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (a *Assign) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "target", Value: a.Target},
		{Name: "node", Value: a.Node},
	}
}

type AssignBlock struct {
	BaseNode

	Target Expression
	Filter *Filter
	Body   []Node
}

func NewAssignBlock(token token.Token, target Expression, filter *Filter, body []Node) *AssignBlock {
	return &AssignBlock{
		BaseNode: NewBaseNode(token),
		Target:   target,
		Filter:   filter,
		Body:     body,
	}
}

func (ab *AssignBlock) statementNode() {}

func (ab *AssignBlock) String() string {
	return fmt.Sprintf("nodes.AssignBlock(target=%s, filter=%s, body=%s)",
		ab.Target.String(), ab.Filter.String(), reprNodeList(ab.Body))
}

func (ab *AssignBlock) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.AssignBlock(")
	fmt.Fprintf(sb, "  target=%s,\n", ab.Target.Dumps(indent+2))
	sb.WriteIndent()
	fmt.Fprintf(sb, "  filter=%s,\n", ab.Filter.Dumps(indent+2))
	sb.WriteIndent()
	sb.WriteNodeList("body", ab.Body)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (ab *AssignBlock) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "target", Value: ab.Target},
		{Name: "filter", Value: ab.Filter},
		{Name: "body", Values: ab.Body},
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

	Body []Node
}

func NewScope(token token.Token, body []Node) *Scope {
	return &Scope{
		BaseNode: NewBaseNode(token),
		Body:     body,
	}
}

func (s *Scope) statementNode() {}

func (s *Scope) String() string {
	return fmt.Sprintf("nodes.Scope(body=%s)", reprNodeList(s.Body))
}

func (s *Scope) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.Scope(")
	sb.WriteNodeList("body", s.Body)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (s *Scope) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "body", Values: s.Body},
	}
}

type OverlayScopy struct {
	BaseNode

	Context Expression
	Body    []Node
}

func NewOverlayScopy(token token.Token, context Expression, body []Node) *OverlayScopy {
	return &OverlayScopy{
		BaseNode: NewBaseNode(token),
		Context:  context,
		Body:     body,
	}
}

func (os *OverlayScopy) statementNode() {}

func (os *OverlayScopy) String() string {
	return fmt.Sprintf("nodes.OverlayScopy(context=%s, body=%s)", os.Context.String(), reprNodeList(os.Body))
}

func (os *OverlayScopy) Dumps(indent int) string {
	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.OverlayScopy(")
	fmt.Fprintf(sb, "  context=%s,\n", os.Context.Dumps(indent+2))
	sb.WriteIndent()
	sb.WriteNodeList("body", os.Body)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (os *OverlayScopy) ChildNodes() []ChildNode {
	return []ChildNode{
		{Name: "context", Value: os.Context},
		{Name: "body", Values: os.Body},
	}
}

type EvalContextModifier struct {
	BaseNode

	Options []*Keyword
}

func NewEvalContextModifier(token token.Token, options []*Keyword) *EvalContextModifier {
	return &EvalContextModifier{
		BaseNode: NewBaseNode(token),
		Options:  options,
	}
}

func (ecm *EvalContextModifier) statementNode() {}

func (ecm *EvalContextModifier) String() string {
	options := make([]Node, len(ecm.Options))
	for i := range ecm.Options {
		options[i] = ecm.Options[i]
	}
	return fmt.Sprintf("nodes.EvalContextModifier(options=%s)", reprNodeList(options))
}

func (ecm *EvalContextModifier) Dumps(indent int) string {
	options := make([]Node, len(ecm.Options))
	for i := range ecm.Options {
		options[i] = ecm.Options[i]
	}

	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.EvalContextModifier(")
	sb.WriteNodeList("options", options)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (ecm *EvalContextModifier) ChildNodes() []ChildNode {
	options := make([]Node, len(ecm.Options))
	for i := range ecm.Options {
		options[i] = ecm.Options[i]
	}

	return []ChildNode{
		{Name: "options", Values: options},
	}
}

type ScopedEvalContextModifier struct {
	EvalContextModifier

	Body []Node
}

func NewScopedEvalContextModifier(token token.Token, options []*Keyword, body []Node) *ScopedEvalContextModifier {
	return &ScopedEvalContextModifier{
		EvalContextModifier: *NewEvalContextModifier(token, options),
		Body:                body,
	}
}

func (sem *ScopedEvalContextModifier) statementNode() {}

func (sem *ScopedEvalContextModifier) String() string {
	options := make([]Node, len(sem.Options))
	for i := range sem.Options {
		options[i] = sem.Options[i]
	}
	return fmt.Sprintf("nodes.ScopedEvalContextModifier(options=%s, body=%s)", reprNodeList(options), reprNodeList(sem.Body))
}

func (sem *ScopedEvalContextModifier) Dumps(indent int) string {
	options := make([]Node, len(sem.Options))
	for i := range sem.Options {
		options[i] = sem.Options[i]
	}

	sb := newStringBuilder(indent)
	sb.WriteLineIndent("nodes.ScopedEvalContextModifier(")
	sb.WriteNodeList("options", options)
	sb.WriteIndent()
	sb.WriteNodeList("body", sem.Body)
	sb.WriteIndent()
	sb.WriteString(")")
	return sb.String()
}

func (sem *ScopedEvalContextModifier) ChildNodes() []ChildNode {
	childnodes := sem.EvalContextModifier.ChildNodes()
	childnodes = append(childnodes, ChildNode{Name: "body", Values: sem.Body})

	return childnodes
}
