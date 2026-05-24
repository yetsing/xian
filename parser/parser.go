package parser

import (
	"fmt"
	"slices"
	"strings"

	"github.com/yetsing/xian/ast"
	"github.com/yetsing/xian/lexer"
	"github.com/yetsing/xian/templateerror"
	"github.com/yetsing/xian/token"
)

func _t1(ttype token.TokenType) token.Token {
	return token.Token{Type: ttype}
}

func _t2(ttype token.TokenType, literal string) token.Token {
	return token.Token{Type: ttype, Literal: literal}
}

var statementKeywords = map[string]struct{}{
	"for":        {},
	"if":         {},
	"block":      {},
	"extends":    {},
	"print":      {},
	"macro":      {},
	"include":    {},
	"from":       {},
	"import":     {},
	"set":        {},
	"with":       {},
	"autoescape": {},
}

type Parser struct {
	stream   *lexer.TokenStream
	name     string
	filename string

	tagStack      []string
	endTokenStack [][]token.Token
}

func NewParser(stream *lexer.TokenStream, name string, filename string) *Parser {
	return &Parser{
		stream:   stream,
		name:     name,
		filename: filename,
		tagStack: []string{},
	}
}

func (p *Parser) Parse() (ast.Template, error) {
	tok := p.stream.Current()
	body, err := p.subparse(nil)
	if err != nil {
		return ast.Template{}, err
	}
	if !p.stream.Eos() {
		return ast.Template{}, p.fail("expected end of template", p.stream.Current())
	}
	return ast.Template{
		BaseNode: ast.NewBaseNode(tok),
		Body:     body,
	}, nil
}

func (p *Parser) fail(msg string, tok token.Token) error {
	return templateerror.NewTemplateSyntaxError(msg, tok.Start.Line, p.name, p.filename)
}

func (p *Parser) failf(format string, args ...any) error {
	tok := p.stream.Current()
	return p.fail(fmt.Sprintf(format, args...), tok)
}

func (p *Parser) failtf(tok token.Token, format string, args ...any) error {
	return p.fail(fmt.Sprintf(format, args...), tok)
}

func (p *Parser) isTupleEnd(extraEndRules []token.Token) bool {
	if p.stream.CurrentIn1(token.TOKEN_VARIABLE_END, token.TOKEN_BLOCK_END, token.TOKEN_RPAREN) {
		return true
	}
	if len(extraEndRules) > 0 {
		return p.stream.CurrentTestAny(extraEndRules...)
	}
	return false
}

func (p *Parser) parseStatement() (ast.Node, error) {
	tok := p.stream.Current()
	if tok.Type != token.TOKEN_NAME {
		return nil, p.fail("tag name expected", tok)
	}

	p.tagStack = append(p.tagStack, tok.Literal)
	defer func() {
		if len(p.tagStack) > 0 {
			p.tagStack = p.tagStack[:len(p.tagStack)-1]
		}
	}()

	switch tok.Literal {
	case "for":
		return p.parseFor()
	case "if":
		return p.parseIf()
	case "block":
		return p.parseBlock()
	case "extends":
		return p.parseExtends()
	case "print":
		return p.parsePrint()
	case "macro":
		return p.parseMacro()
	case "include":
		return p.parseInclude()
	case "from":
		return p.parseFrom()
	case "import":
		return p.parseImport()
	case "set":
		return p.parseSet()
	case "with":
		return p.parseWith()
	case "autoescape":
		return p.parseAutoescape()

	case "call":
		return p.parseCallBlock()
	case "filter":
		return p.parseFilterBlock()
	}

	return nil, p.failf("unknown tag %q", tok.Literal)
}

func (p *Parser) parseStatements(endTokens []token.Token, dropNeedle bool) ([]ast.Node, error) {
	p.stream.SkipIf1(token.TOKEN_COLON)

	if _, err := p.stream.Expect1(token.TOKEN_BLOCK_END); err != nil {
		return nil, err
	}

	result, err := p.subparse(endTokens)
	if err != nil {
		return nil, err
	}

	if p.stream.CurrentEqual1(token.TOKEN_EOF) {
		return nil, p.fail("unexpected end of template", p.stream.Current())
	}

	if dropNeedle {
		p.stream.Next()
	}

	return result, nil
}

func (p *Parser) parseSet() (ast.Node, error) {
	tok := p.stream.Next()
	target, err := p.parseAssignTarget(true, false, nil, true)
	if err != nil {
		return nil, err
	}
	if p.stream.SkipIf1(token.TOKEN_ASSIGN) {
		expr, err := p.parseTuple(false, true, nil, false, false)
		if err != nil {
			return nil, err
		}
		return ast.NewAssign(tok, target, expr), nil
	}
	filterNodeExpr, err := p.parseFilter(nil, false)
	if err != nil {
		return nil, err
	}
	filterNode, ok := filterNodeExpr.(*ast.Filter)
	if !ok {
		return nil, p.fail("expected filter expression", filterNodeExpr.Token())
	}
	body, err := p.parseStatements([]token.Token{_t2(token.TOKEN_NAME, "endset")}, true)
	if err != nil {
		return nil, err
	}
	return ast.NewAssignBlock(tok, target, filterNode, body), nil
}

func (p *Parser) parseFor() (ast.Node, error) {
	tok, err := p.stream.Expect2(token.TOKEN_NAME, "for")
	if err != nil {
		return nil, err
	}
	target, err := p.parseAssignTarget(true, false, []token.Token{_t2(token.TOKEN_NAME, "in")}, false)
	if err != nil {
		return nil, err
	}
	if _, err := p.stream.Expect2(token.TOKEN_NAME, "in"); err != nil {
		return nil, err
	}
	iter, err := p.parseTuple(false, false, []token.Token{_t2(token.TOKEN_NAME, "recursive")}, false, false)
	if err != nil {
		return nil, err
	}
	var test ast.Expression
	if p.stream.SkipIf2(token.TOKEN_NAME, "if") {
		test, err = p.parseExpression(true)
		if err != nil {
			return nil, err
		}
	}
	recursive := p.stream.SkipIf2(token.TOKEN_NAME, "recursive")
	body, err := p.parseStatements([]token.Token{_t2(token.TOKEN_NAME, "endfor"), _t2(token.TOKEN_NAME, "else")}, false)
	if err != nil {
		return nil, err
	}
	var elseBody []ast.Node
	if endTk := p.stream.Next(); endTk.Literal == "endfor" {
		elseBody = []ast.Node{}
	} else {
		elseBody, err = p.parseStatements([]token.Token{_t2(token.TOKEN_NAME, "endfor")}, true)
		if err != nil {
			return nil, err
		}
	}
	return ast.NewFor(tok, target, iter, body, elseBody, test, recursive), nil
}

func (p *Parser) parseIf() (ast.Node, error) {
	tok, err := p.stream.Expect2(token.TOKEN_NAME, "if")
	if err != nil {
		return nil, err
	}
	node := ast.NewIf(tok, nil, nil, nil, nil)
	result := node
	for {
		node.Test, err = p.parseTuple(false, false, nil, false, false)
		if err != nil {
			return nil, err
		}
		node.Body, err = p.parseStatements([]token.Token{_t2(token.TOKEN_NAME, "elif"), _t2(token.TOKEN_NAME, "else"), _t2(token.TOKEN_NAME, "endif")}, false)
		if err != nil {
			return nil, err
		}
		tok := p.stream.Next()
		if tok.Test(_t2(token.TOKEN_NAME, "elif")) {
			node = ast.NewIf(p.stream.Current(), nil, nil, nil, nil)
			result.Elif_ = append(result.Elif_, node)
			continue
		} else if tok.Test(_t2(token.TOKEN_NAME, "else")) {
			elseBody, err := p.parseStatements([]token.Token{_t2(token.TOKEN_NAME, "endif")}, true)
			if err != nil {
				return nil, err
			}
			result.Else_ = elseBody
		}
		break
	}
	return result, nil
}

func (p *Parser) parseWith() (ast.Node, error) {
	tok := p.stream.Next()
	targets := []ast.Expression{}
	values := []ast.Expression{}
	for p.stream.CurrentNotEqual1(token.TOKEN_BLOCK_END) {
		if len(targets) > 0 {
			if _, err := p.stream.Expect1(token.TOKEN_COMMA); err != nil {
				return nil, err
			}
		}
		target, err := p.parseAssignTarget(true, false, nil, false)
		if err != nil {
			return nil, err
		}
		ast.SetCtx(target, ast.ExprContextParam)
		targets = append(targets, target)
		if _, err := p.stream.Expect1(token.TOKEN_ASSIGN); err != nil {
			return nil, err
		}
		value, err := p.parseExpression(true)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	body, err := p.parseStatements([]token.Token{_t2(token.TOKEN_NAME, "endwith")}, true)
	if err != nil {
		return nil, err
	}
	return ast.NewWith(tok, targets, values, body), nil
}

func (p *Parser) parseAutoescape() (ast.Node, error) {
	node := ast.NewScopedEvalContextModifier(p.stream.Next(), nil, nil)
	value, err := p.parseExpression(true)
	if err != nil {
		return nil, err
	}
	node.Options = []*ast.Keyword{ast.NewKeyword(value.Token(), "autoescape", value)}
	body, err := p.parseStatements([]token.Token{_t2(token.TOKEN_NAME, "endautoescape")}, true)
	if err != nil {
		return nil, err
	}
	node.Body = body
	return ast.NewScope(node.Token(), []ast.Node{node}), nil
}

func (p *Parser) parseBlock() (ast.Node, error) {
	node := ast.NewBlock(p.stream.Next(), "", nil, false, false)
	t, err := p.stream.Expect1(token.TOKEN_NAME)
	if err != nil {
		return nil, err
	}
	node.Name = t.Literal
	node.Scoped = p.stream.SkipIf2(token.TOKEN_NAME, "scoped")
	node.Required = p.stream.SkipIf2(token.TOKEN_NAME, "required")

	if p.stream.CurrentEqual1(token.TOKEN_SUB) {
		return nil, p.fail("block names can not contain hyphens", p.stream.Current())
	}

	node.Body, err = p.parseStatements([]token.Token{_t2(token.TOKEN_NAME, "endblock")}, true)
	if err != nil {
		return nil, err
	}

	if node.Required {
		for _, bodyNode := range node.Body {
			out, ok := bodyNode.(*ast.Output)
			if !ok {
				return nil, p.fail("Required blocks can only contain comments or whitespace", bodyNode.Token())
			}
			if slices.ContainsFunc(out.Nodes, func(outputNode ast.Expression) bool {
				tpl, ok := outputNode.(*ast.TemplateData)
				if !ok {
					return true
				}
				if len(strings.TrimSpace(tpl.Data)) > 0 {
					return true
				}
				return false
			}) {
				return nil, p.fail("Required blocks can only contain comments or whitespace", bodyNode.Token())
			}
		}
	}

	p.stream.SkipIf2(token.TOKEN_NAME, node.Name)
	return node, nil
}

func (p *Parser) parseExtends() (ast.Node, error) {
	tok := p.stream.Next()
	tpl, err := p.parseExpression(true)
	if err != nil {
		return nil, err
	}
	return ast.NewExtends(tok, tpl), nil
}

func (p *Parser) parseImportContext(defaultVal bool) (bool, error) {
	if (p.stream.CurrentTestAny(_t2(token.TOKEN_NAME, "with"), _t2(token.TOKEN_NAME, "without"))) && p.stream.PeekEqual2(token.TOKEN_NAME, "context") {
		withToken := p.stream.Next()
		withCtx := withToken.Literal == "with"
		p.stream.Skip()
		return withCtx, nil
	}
	return defaultVal, nil
}

func (p *Parser) parseInclude() (ast.Node, error) {
	tok := p.stream.Next()
	tpl, err := p.parseExpression(true)
	if err != nil {
		return nil, err
	}
	ignoreMissing := false
	if p.stream.CurrentEqual2(token.TOKEN_NAME, "ignore") && p.stream.PeekEqual2(token.TOKEN_NAME, "missing") {
		ignoreMissing = true
		p.stream.SkipN(2)
	}
	withCtx, err := p.parseImportContext(true)
	if err != nil {
		return nil, err
	}
	return ast.NewInclude(tok, tpl, withCtx, ignoreMissing), nil
}

func (p *Parser) parseImport() (ast.Node, error) {
	tok := p.stream.Next()
	tpl, err := p.parseExpression(true)
	if err != nil {
		return nil, err
	}
	if _, err := p.stream.Expect2(token.TOKEN_NAME, "as"); err != nil {
		return nil, err
	}
	targetNode, err := p.parseAssignTarget(true, true, nil, false)
	if err != nil {
		return nil, err
	}
	name, ok := targetNode.(*ast.Name)
	if !ok {
		return nil, p.fail("import target must be a name", targetNode.Token())
	}
	withCtx, err := p.parseImportContext(false)
	if err != nil {
		return nil, err
	}
	return ast.NewImport(tok, tpl, name.Name, withCtx), nil
}

func (p *Parser) parseFrom() (ast.Node, error) {
	tok := p.stream.Next()
	tpl, err := p.parseExpression(true)
	if err != nil {
		return nil, err
	}
	if _, err := p.stream.Expect2(token.TOKEN_NAME, "import"); err != nil {
		return nil, err
	}
	names := []string{}
	aliases := []string{}
	withContext := false

	parseContext := func() bool {
		if (p.stream.CurrentLiteralIn("with", "without")) && p.stream.PeekEqual2(token.TOKEN_NAME, "context") {
			withToken := p.stream.Next()
			withContext = withToken.Literal == "with"
			p.stream.Skip()
			return true
		}
		return false
	}

	for {
		if len(names) > 0 {
			if _, err := p.stream.Expect1(token.TOKEN_COMMA); err != nil {
				return nil, err
			}
		}
		if p.stream.CurrentEqual1(token.TOKEN_NAME) {
			if parseContext() {
				break
			}
			targetNode, err := p.parseAssignTarget(true, true, nil, false)
			if err != nil {
				return nil, err
			}
			target, ok := targetNode.(*ast.Name)
			if !ok {
				return nil, p.fail("import target must be a name", targetNode.Token())
			}
			if strings.HasPrefix(target.Name, "_") {
				return nil, p.fail(
					"names starting with an underline can not be imported",
					target.Token(),
				)
			}
			names = append(names, target.Name)
			if p.stream.SkipIf2(token.TOKEN_NAME, "as") {
				asTargetNode, err := p.parseAssignTarget(true, true, nil, false)
				if err != nil {
					return nil, err
				}
				asTarget, ok := asTargetNode.(*ast.Name)
				if !ok {
					return nil, p.fail("import alias must be a name", asTargetNode.Token())
				}
				aliases = append(aliases, asTarget.Name)
			} else {
				aliases = append(aliases, "") // 占位符，表示使用原名
			}
			if parseContext() || p.stream.CurrentNotEqual1(token.TOKEN_COMMA) {
				break
			}
		} else {
			if _, err := p.stream.Expect1(token.TOKEN_NAME); err != nil {
				return nil, err
			}
		}
	}

	return ast.NewFromImport(tok, tpl, names, aliases, withContext), nil
}

func (p *Parser) parseSignature() ([]*ast.Name, []ast.Expression, error) {
	if _, err := p.stream.Expect1(token.TOKEN_LPAREN); err != nil {
		return nil, nil, err
	}
	args := []*ast.Name{}
	defaults := []ast.Expression{}
	for p.stream.CurrentNotEqual1(token.TOKEN_RPAREN) {
		if len(args) > 0 {
			if _, err := p.stream.Expect1(token.TOKEN_COMMA); err != nil {
				return nil, nil, err
			}
		}
		argExpr, err := p.parseAssignTarget(true, true, nil, false)
		if err != nil {
			return nil, nil, err
		}
		arg, ok := argExpr.(*ast.Name)
		if !ok {
			return nil, nil, p.fail("expected parameter name", argExpr.Token())
		}
		ast.SetCtx(arg, ast.ExprContextParam)
		if p.stream.SkipIf1(token.TOKEN_ASSIGN) {
			defv, err := p.parseExpression(true)
			if err != nil {
				return nil, nil, err
			}
			defaults = append(defaults, defv)
		} else if len(defaults) > 0 {
			return nil, nil, p.fail("non-default argument follows default argument", p.stream.Current())
		}
		args = append(args, arg)
	}
	if _, err := p.stream.Expect1(token.TOKEN_RPAREN); err != nil {
		return nil, nil, err
	}
	return args, defaults, nil
}

func (p *Parser) parseCallBlock() (ast.Node, error) {
	tok := p.stream.Next()
	args := []*ast.Name{}
	defaults := []ast.Expression{}
	if p.stream.CurrentEqual1(token.TOKEN_LPAREN) {
		var err error
		args, defaults, err = p.parseSignature()
		if err != nil {
			return nil, err
		}
	}
	callExpr, err := p.parseExpression(true)
	if err != nil {
		return nil, err
	}
	call, ok := callExpr.(*ast.Call)
	if !ok {
		return nil, p.fail("expected call", tok)
	}
	body, err := p.parseStatements([]token.Token{_t2(token.TOKEN_NAME, "endcall")}, true)
	if err != nil {
		return nil, err
	}
	return ast.NewCallBlock(tok, call, args, defaults, body), nil
}

func (p *Parser) parseFilterBlock() (ast.Node, error) {
	tok := p.stream.Next()
	expr, err := p.parseFilter(nil, true)
	if err != nil {
		return nil, err
	}
	filter, ok := expr.(*ast.Filter)
	if !ok {
		return nil, p.fail("expected filter expression", expr.Token())
	}
	body, err := p.parseStatements([]token.Token{_t2(token.TOKEN_NAME, "endfilter")}, true)
	if err != nil {
		return nil, err
	}
	return ast.NewFilterBlock(tok, body, filter), nil
}

func (p *Parser) parseMacro() (ast.Node, error) {
	tok := p.stream.Next()
	rv, err := p.parseAssignTarget(true, true, nil, false)
	if err != nil {
		return nil, err
	}
	target, ok := rv.(*ast.Name)
	if !ok {
		return nil, p.fail("expected macro name", rv.Token())
	}
	args, defaults, err := p.parseSignature()
	if err != nil {
		return nil, err
	}
	body, err := p.parseStatements(
		[]token.Token{_t2(token.TOKEN_NAME, "endmacro")},
		true,
	)
	if err != nil {
		return nil, err
	}

	return ast.NewMacro(tok, target.Name, args, defaults, body), nil
}

func (p *Parser) parsePrint() (ast.Node, error) {
	tok := p.stream.Next()
	nodes := []ast.Expression{}
	for p.stream.CurrentNotEqual1(token.TOKEN_BLOCK_END) {
		if len(nodes) > 0 {
			_, err := p.stream.Expect1(token.TOKEN_COMMA)
			if err != nil {
				return nil, err
			}
		}
		expr, err := p.parseExpression(true)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, expr)
	}
	return ast.NewOutput(tok, nodes), nil
}

func (p *Parser) parseAssignTarget(withTuple bool, nameOnly bool, extraEndRules []token.Token, withNamespace bool) (ast.Expression, error) {
	var target ast.Expression
	var err error
	if nameOnly {
		tok, err := p.stream.Expect1(token.TOKEN_NAME)
		if err != nil {
			return nil, err
		}
		target = ast.NewName(tok, tok.Literal, ast.ExprContextStore)
	} else {
		if withTuple {
			target, err = p.parseTuple(true, true, extraEndRules, false, withNamespace)
		} else {
			target, err = p.parsePrimary(withNamespace)
		}
		if err != nil {
			return nil, err
		}
		ast.SetCtx(target, ast.ExprContextStore)
	}
	if !target.CanAssign() {
		return nil, p.fail(fmt.Sprintf("can't assign to %T", target), target.Token())
	}
	return target, nil
}

func (p *Parser) parseExpression(withCondexpr bool) (ast.Expression, error) {
	if withCondexpr {
		return p.parseCondexpr()
	}
	return p.parseOr()
}

func (p *Parser) parseCondexpr() (ast.Expression, error) {
	tok := p.stream.Current()
	expr1, err := p.parseOr()
	if err != nil {
		return nil, err
	}
	for p.stream.SkipIf2(token.TOKEN_NAME, "if") {
		expr2, err := p.parseOr()
		if err != nil {
			return nil, err
		}
		var expr3 ast.Expression
		if p.stream.SkipIf2(token.TOKEN_NAME, "else") {
			expr3, err = p.parseCondexpr()
			if err != nil {
				return nil, err
			}
		}
		expr1 = ast.NewCondExpr(tok, expr2, expr1, expr3)
		tok = p.stream.Current()
	}
	return expr1, nil
}

func (p *Parser) parseOr() (ast.Expression, error) {
	tok := p.stream.Current()
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	for p.stream.SkipIf2(token.TOKEN_NAME, "or") {
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = ast.NewOr(tok, left, right)
		tok = p.stream.Current()
	}
	return left, nil
}

func (p *Parser) parseAnd() (ast.Expression, error) {
	tok := p.stream.Current()
	left, err := p.parseNot()
	if err != nil {
		return nil, err
	}
	for p.stream.SkipIf2(token.TOKEN_NAME, "and") {
		right, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		left = ast.NewAnd(tok, left, right)
		tok = p.stream.Current()
	}
	return left, nil
}

func (p *Parser) parseNot() (ast.Expression, error) {
	if p.stream.CurrentEqual2(token.TOKEN_NAME, "not") {
		tok := p.stream.Next()
		expr, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		return ast.NewNot(tok, expr), nil
	}
	return p.parseCompare()
}

func (p *Parser) parseCompare() (ast.Expression, error) {
	tok := p.stream.Current()
	expr, err := p.parseMath1()
	if err != nil {
		return nil, err
	}
	ops := []*ast.Operand{}
	for {
		ct := p.stream.Current()
		if p.stream.CurrentIn1(token.TOKEN_EQ, token.TOKEN_NE, token.TOKEN_LT, token.TOKEN_LTEQ, token.TOKEN_GT, token.TOKEN_GTEQ) {
			op := ct.Type
			p.stream.Next()
			rhs, err := p.parseMath1()
			if err != nil {
				return nil, err
			}
			ops = append(ops, ast.NewOperand(ct, string(op), rhs))
		} else if p.stream.SkipIf2(token.TOKEN_NAME, "in") {
			rhs, err := p.parseMath1()
			if err != nil {
				return nil, err
			}
			ops = append(ops, ast.NewOperand(ct, "in", rhs))
		} else if p.stream.CurrentEqual2(token.TOKEN_NAME, "not") && p.stream.PeekEqual2(token.TOKEN_NAME, "in") {
			p.stream.SkipN(2)
			rhs, err := p.parseMath1()
			if err != nil {
				return nil, err
			}
			ops = append(ops, ast.NewOperand(ct, "notin", rhs))
		} else {
			break
		}
	}
	if len(ops) == 0 {
		return expr, nil
	}
	return ast.NewCompare(tok, expr, ops), nil
}

func (p *Parser) parseMath1() (ast.Expression, error) {
	tok := p.stream.Current()
	left, err := p.parseConcat()
	if err != nil {
		return nil, err
	}
	for p.stream.CurrentIn1(token.TOKEN_ADD, token.TOKEN_SUB) {
		op := p.stream.Current()
		p.stream.Next()
		right, err := p.parseConcat()
		if err != nil {
			return nil, err
		}
		if op.Type == token.TOKEN_ADD {
			left = ast.NewAdd(tok, left, right)
		} else {
			left = ast.NewSub(tok, left, right)
		}
		tok = p.stream.Current()
	}
	return left, nil
}

func (p *Parser) parseConcat() (ast.Expression, error) {
	tok := p.stream.Current()
	first, err := p.parseMath2()
	if err != nil {
		return nil, err
	}
	items := []ast.Expression{first}
	for p.stream.CurrentEqual1(token.TOKEN_TILDE) {
		p.stream.Next()
		nextItem, err := p.parseMath2()
		if err != nil {
			return nil, err
		}
		items = append(items, nextItem)
	}
	if len(items) == 1 {
		return items[0], nil
	}
	return ast.NewConcat(tok, items), nil
}

func (p *Parser) parseMath2() (ast.Expression, error) {
	tok := p.stream.Current()
	left, err := p.parsePow()
	if err != nil {
		return nil, err
	}
	for p.stream.CurrentIn1(token.TOKEN_MUL, token.TOKEN_DIV, token.TOKEN_FLOORDIV, token.TOKEN_MOD) {
		op := p.stream.Next()
		right, err := p.parsePow()
		if err != nil {
			return nil, err
		}
		switch op.Type {
		case token.TOKEN_MUL:
			left = ast.NewMul(tok, left, right)
		case token.TOKEN_DIV:
			left = ast.NewDiv(tok, left, right)
		case token.TOKEN_FLOORDIV:
			left = ast.NewFloorDiv(tok, left, right)
		case token.TOKEN_MOD:
			left = ast.NewMod(tok, left, right)
		}
		tok = p.stream.Current()
	}
	return left, nil
}

func (p *Parser) parsePow() (ast.Expression, error) {
	tok := p.stream.Current()
	left, err := p.parseUnary(true)
	if err != nil {
		return nil, err
	}
	for p.stream.CurrentEqual1(token.TOKEN_POW) {
		p.stream.Next()
		right, err := p.parseUnary(true)
		if err != nil {
			return nil, err
		}
		left = ast.NewPow(tok, left, right)
		tok = p.stream.Current()
	}
	return left, nil
}

func (p *Parser) parseUnary(withFilter bool) (ast.Expression, error) {
	tok := p.stream.Current()
	var node ast.Expression
	var err error
	if tok.Type == token.TOKEN_SUB {
		p.stream.Next()
		node, err = p.parseUnary(false)
		if err != nil {
			return nil, err
		}
		node = ast.NewNeg(tok, node)
	} else if tok.Type == token.TOKEN_ADD {
		p.stream.Next()
		node, err = p.parseUnary(false)
		if err != nil {
			return nil, err
		}
		node = ast.NewPos(tok, node)
	} else {
		node, err = p.parsePrimary(false)
		if err != nil {
			return nil, err
		}
	}
	node, err = p.parsePostfix(node)
	if err != nil {
		return nil, err
	}
	if withFilter {
		node, err = p.parseFilterExpr(node)
		if err != nil {
			return nil, err
		}
	}
	return node, nil
}

func (p *Parser) parsePrimary(withNamespace bool) (ast.Expression, error) {
	tok := p.stream.Current()
	if tok.Type == token.TOKEN_NAME {
		p.stream.Next()
		switch tok.Literal {
		case "true", "True":
			return ast.NewConst(tok, true), nil
		case "false", "False":
			return ast.NewConst(tok, false), nil
		case "none", "None":
			return ast.NewConst(tok, nil), nil
		default:
			if withNamespace && p.stream.CurrentEqual1(token.TOKEN_DOT) {
				p.stream.Next()
				attr, err := p.stream.Expect1(token.TOKEN_NAME)
				if err != nil {
					return nil, err
				}
				return ast.NewNSRef(tok, tok.Literal, attr.Literal), nil
			}
			return ast.NewName(tok, tok.Literal, ast.ExprContextLoad), nil
		}
	}
	if tok.Type == token.TOKEN_STRING {
		p.stream.Next()
		buf := tok.Literal
		for p.stream.CurrentEqual1(token.TOKEN_STRING) {
			buf += p.stream.Current().Literal
			p.stream.Next()
		}
		return ast.NewConst(tok, buf), nil
	}
	if tok.Type == (token.TOKEN_INTEGER) {
		p.stream.Next()
		v, err := toint64(tok.Literal)
		if err != nil {
			return nil, p.fail("invalid integer literal", tok)
		}
		return ast.NewConst(tok, v), nil
	}
	if tok.Type == token.TOKEN_FLOAT {
		p.stream.Next()
		v, err := tofloat64(tok.Literal)
		if err != nil {
			return nil, p.fail("invalid float literal", tok)
		}
		return ast.NewConst(tok, v), nil
	}
	if tok.Type == token.TOKEN_LPAREN {
		p.stream.Next()
		node, err := p.parseTuple(false, true, nil, true, false)
		if err != nil {
			return nil, err
		}
		if _, err := p.stream.Expect1(token.TOKEN_RPAREN); err != nil {
			return nil, err
		}
		return node, nil
	}
	if tok.Type == token.TOKEN_LBRACKET {
		return p.parseList()
	}
	if tok.Type == token.TOKEN_LBRACE {
		return p.parseDict()
	}
	return nil, p.failtf(tok, "unexpected token %q", tok.Type)
}

func (p *Parser) parseTuple(simplified bool, withCondexpr bool, extraEndRules []token.Token, explicitParentheses bool, withNamespace bool) (ast.Expression, error) {
	tok := p.stream.Current()
	var parseOne func() (ast.Expression, error)
	if simplified {
		parseOne = func() (ast.Expression, error) {
			return p.parsePrimary(withNamespace)
		}
	} else {
		parseOne = func() (ast.Expression, error) {
			return p.parseExpression(withCondexpr)
		}
	}

	args := []ast.Expression{}
	isTuple := false

	for {
		if len(args) > 0 {
			if _, err := p.stream.Expect1(token.TOKEN_COMMA); err != nil {
				return nil, err
			}
		}
		if p.isTupleEnd(extraEndRules) {
			break
		}
		node, err := parseOne()
		if err != nil {
			return nil, err
		}
		args = append(args, node)
		if p.stream.CurrentEqual1(token.TOKEN_COMMA) {
			isTuple = true
		} else {
			break
		}
	}

	if !isTuple {
		if len(args) > 0 {
			return args[0], nil
		}
		if !explicitParentheses {
			return nil, p.failf("Expected an expression, got %q", p.stream.Current().Type)
		}
	}

	return ast.NewTuple(tok, args, ast.ExprContextLoad), nil
}

func (p *Parser) parseList() (ast.Expression, error) {
	tok, err := p.stream.Expect1(token.TOKEN_LBRACKET)
	if err != nil {
		return nil, err
	}
	items := []ast.Expression{}
	for p.stream.CurrentNotEqual1(token.TOKEN_RBRACKET) {
		if len(items) > 0 {
			if _, err := p.stream.Expect1(token.TOKEN_COMMA); err != nil {
				return nil, err
			}
		}
		if p.stream.CurrentEqual1(token.TOKEN_RBRACKET) {
			break
		}
		item, err := p.parseExpression(true)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if _, err := p.stream.Expect1(token.TOKEN_RBRACKET); err != nil {
		return nil, err
	}
	return ast.NewList(tok, items), nil
}

func (p *Parser) parseDict() (ast.Expression, error) {
	tok, err := p.stream.Expect1(token.TOKEN_LBRACE)
	if err != nil {
		return nil, err
	}
	items := []*ast.Pair{}
	for p.stream.CurrentNotEqual1(token.TOKEN_RBRACE) {
		if len(items) > 0 {
			if _, err := p.stream.Expect1(token.TOKEN_COMMA); err != nil {
				return nil, err
			}
		}
		if p.stream.CurrentEqual1(token.TOKEN_RBRACE) {
			break
		}
		key, err := p.parseExpression(true)
		if err != nil {
			return nil, err
		}
		if _, err := p.stream.Expect1(token.TOKEN_COLON); err != nil {
			return nil, err
		}
		value, err := p.parseExpression(true)
		if err != nil {
			return nil, err
		}
		items = append(items, ast.NewPair(key.Token(), key, value))
	}
	if _, err := p.stream.Expect1(token.TOKEN_RBRACE); err != nil {
		return nil, err
	}
	return ast.NewDict(tok, items), nil
}

func (p *Parser) parsePostfix(node ast.Expression) (ast.Expression, error) {
	var err error
	for {
		current := p.stream.Current()
		if current.TypeIn(token.TOKEN_DOT, token.TOKEN_LBRACKET) {
			node, err = p.parseSubscript(node)
			if err != nil {
				return nil, err
			}
		} else if current.Type == token.TOKEN_LPAREN {
			node, err = p.parseCall(node)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	return node, nil
}

func (p *Parser) parseFilterExpr(node ast.Expression) (ast.Expression, error) {
	var err error
	for {
		current := p.stream.Current()
		if current.Type == token.TOKEN_PIPE {
			node, err = p.parseFilter(node, false)
			if err != nil {
				return nil, err
			}
		} else if current.Equal2(token.TOKEN_NAME, "is") {
			node, err = p.parseTest(node)
			if err != nil {
				return nil, err
			}
		} else if current.Type == token.TOKEN_LPAREN {
			node, err = p.parseCall(node)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	return node, nil
}

func (p *Parser) parseSubscript(node ast.Expression) (ast.Expression, error) {
	tok := p.stream.Next()
	var arg ast.Expression

	if tok.Type == token.TOKEN_DOT {
		attrToken := p.stream.Current()
		p.stream.Next()
		if attrToken.Type == token.TOKEN_NAME {
			return ast.NewGetattr(tok, node, attrToken.Literal, ast.ExprContextLoad), nil
		} else if attrToken.Type != token.TOKEN_INTEGER {
			return nil, p.fail("expected name or number", attrToken)
		}
		n, err := toint64(attrToken.Literal)
		if err != nil {
			return nil, p.fail("invalid integer literal", attrToken)
		}
		arg = ast.NewConst(attrToken, n)
		return ast.NewGetitem(tok, node, arg, ast.ExprContextLoad), nil
	}

	if tok.Type == token.TOKEN_LBRACKET {
		args := []ast.Expression{}
		for p.stream.CurrentNotEqual1(token.TOKEN_RBRACKET) {
			if len(args) > 0 {
				if _, err := p.stream.Expect1(token.TOKEN_COMMA); err != nil {
					return nil, err
				}
			}
			val, err := p.parseSubscribed()
			if err != nil {
				return nil, err
			}
			args = append(args, val)
		}
		if len(args) == 1 {
			arg = args[0]
		} else {
			arg = ast.NewTuple(tok, args, ast.ExprContextLoad)
		}
		return ast.NewGetitem(tok, node, arg, ast.ExprContextLoad), nil
	}
	return nil, p.fail("expected subscript expression", tok)
}

func (p *Parser) parseSubscribed() (ast.Expression, error) {
	tok := p.stream.Current()
	args := []ast.Expression{}

	if p.stream.CurrentEqual1(token.TOKEN_COLON) {
		p.stream.Next()
		args = []ast.Expression{nil}
	} else {
		node, err := p.parseExpression(true)
		if err != nil {
			return nil, err
		}
		if p.stream.CurrentNotEqual1(token.TOKEN_COLON) {
			return node, nil
		}
		p.stream.Next()
		args = []ast.Expression{node}
	}

	if p.stream.CurrentEqual1(token.TOKEN_COLON) {
		args = append(args, nil)
	} else if p.stream.CurrentNotIn1(token.TOKEN_RBRACKET, token.TOKEN_COMMA) {
		node, err := p.parseExpression(true)
		if err != nil {
			return nil, err
		}
		args = append(args, node)
	} else {
		args = append(args, nil)
	}

	if p.stream.CurrentEqual1(token.TOKEN_COLON) {
		p.stream.Next()
		if p.stream.CurrentNotIn1(token.TOKEN_RBRACKET, token.TOKEN_COMMA) {
			val, err := p.parseExpression(true)
			if err != nil {
				return nil, err
			}
			args = append(args, val)
		} else {
			args = append(args, nil)
		}
	} else {
		args = append(args, nil)
	}

	return ast.NewSlice(tok, args[0], args[1], args[2]), nil
}

func (p *Parser) parseCallArgs() ([]ast.Expression, []*ast.Keyword, ast.Expression, ast.Expression, error) {
	tok, err := p.stream.Expect1(token.TOKEN_LPAREN)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	args := []ast.Expression{}
	kwargs := []*ast.Keyword{}
	var dynArgs ast.Expression
	var dynKwargs ast.Expression
	requireComma := false

	ensure := func(ok bool) error {
		if ok {
			return nil
		}
		return p.fail("invalid syntax for function call expression", tok)
	}

	for p.stream.CurrentNotEqual1(token.TOKEN_RPAREN) {
		if requireComma {
			if _, err := p.stream.Expect1(token.TOKEN_COMMA); err != nil {
				return nil, nil, nil, nil, err
			}

			// support for trailing comma
			if p.stream.CurrentEqual1(token.TOKEN_RPAREN) {
				break
			}
		}

		if p.stream.CurrentEqual1(token.TOKEN_MUL) {
			if err := ensure(dynArgs == nil && dynKwargs == nil); err != nil {
				return nil, nil, nil, nil, err
			}
			p.stream.Next()
			dynArgs, err = p.parseExpression(true)
			if err != nil {
				return nil, nil, nil, nil, err
			}
		} else if p.stream.CurrentEqual1(token.TOKEN_POW) {
			if err := ensure(dynKwargs == nil); err != nil {
				return nil, nil, nil, nil, err
			}
			p.stream.Next()
			dynKwargs, err = p.parseExpression(true)
			if err != nil {
				return nil, nil, nil, nil, err
			}
		} else {
			if p.stream.CurrentEqual1(token.TOKEN_NAME) && p.stream.PeekEqual1(token.TOKEN_ASSIGN) {
				// Parsing a kwarg
				if err := ensure(dynKwargs == nil); err != nil {
					return nil, nil, nil, nil, err
				}
				key := p.stream.Current().Literal
				p.stream.SkipN(2)
				value, err := p.parseExpression(true)
				if err != nil {
					return nil, nil, nil, nil, err
				}
				kwargs = append(kwargs, ast.NewKeyword(value.Token(), key, value))
			} else {
				// Parsing an arg
				if err := ensure(dynArgs == nil && dynKwargs == nil && len(kwargs) == 0); err != nil {
					return nil, nil, nil, nil, err
				}
				val, err := p.parseExpression(true)
				if err != nil {
					return nil, nil, nil, nil, err
				}
				args = append(args, val)
			}
		}

		requireComma = true
	}

	if _, err := p.stream.Expect1(token.TOKEN_RPAREN); err != nil {
		return nil, nil, nil, nil, err
	}
	return args, kwargs, dynArgs, dynKwargs, nil
}

func (p *Parser) parseCall(node ast.Expression) (ast.Expression, error) {
	tok := p.stream.Current()
	args, kwargs, dynArgs, dynKwargs, err := p.parseCallArgs()
	if err != nil {
		return nil, err
	}
	return ast.NewCall(tok, node, args, kwargs, dynArgs, dynKwargs), nil
}

func (p *Parser) parseFilter(node ast.Expression, startInline bool) (ast.Expression, error) {
	for p.stream.CurrentEqual1(token.TOKEN_PIPE) || startInline {
		if !startInline {
			p.stream.Next()
		}
		tok, err := p.stream.Expect1(token.TOKEN_NAME)
		if err != nil {
			return nil, err
		}
		name := tok.Literal
		for p.stream.CurrentEqual1(token.TOKEN_DOT) {
			p.stream.Next()
			cur, err := p.stream.Expect1(token.TOKEN_NAME)
			if err != nil {
				return nil, err
			}
			name += "." + cur.Literal
		}

		var (
			args      []ast.Expression
			kwargs    []*ast.Keyword
			dynArgs   ast.Expression
			dynKwargs ast.Expression
		)
		if p.stream.CurrentEqual1(token.TOKEN_LPAREN) {
			args, kwargs, dynArgs, dynKwargs, err = p.parseCallArgs()
			if err != nil {
				return nil, err
			}
		}
		node = ast.NewFilter(tok, node, name, args, kwargs, dynArgs, dynKwargs)
		startInline = false
	}
	return node, nil
}

func (p *Parser) parseTest(node ast.Expression) (ast.Expression, error) {
	var err error
	tok := p.stream.Next()
	negated := false
	if p.stream.CurrentEqual2(token.TOKEN_NAME, "not") {
		p.stream.Next()
		negated = true
	}
	cur, err := p.stream.Expect1(token.TOKEN_NAME)
	if err != nil {
		return nil, err
	}
	name := cur.Literal
	for p.stream.CurrentEqual1(token.TOKEN_DOT) {
		p.stream.Next()
		cur, err = p.stream.Expect1(token.TOKEN_NAME)
		if err != nil {
			return nil, err
		}
		name += "." + cur.Literal
	}
	var (
		args      []ast.Expression
		kwargs    []*ast.Keyword
		dynArgs   ast.Expression
		dynKwargs ast.Expression
	)
	if p.stream.CurrentEqual1(token.TOKEN_LPAREN) {
		args, kwargs, dynArgs, dynKwargs, err = p.parseCallArgs()
		if err != nil {
			return nil, err
		}
	} else if p.stream.CurrentIn1(
		token.TOKEN_NAME,
		token.TOKEN_STRING,
		token.TOKEN_INTEGER,
		token.TOKEN_FLOAT,
		token.TOKEN_LPAREN,
		token.TOKEN_LBRACKET,
		token.TOKEN_LBRACE,
	) && p.stream.CurrentTestAny(_t2(token.TOKEN_NAME, "else"), _t2(token.TOKEN_PIPE, "or"), _t2(token.TOKEN_PIPE, "and")) {
		if p.stream.CurrentEqual2(token.TOKEN_NAME, "is") {
			return nil, p.fail("You cannot chain multiple tests with is", p.stream.Current())
		}
		argNode, err := p.parsePrimary(false)
		if err != nil {
			return nil, err
		}
		argNode, err = p.parsePostfix(argNode)
		if err != nil {
			return nil, err
		}
		args = []ast.Expression{argNode}
	}
	node = ast.NewTest(tok, node, name, args, kwargs, dynArgs, dynKwargs)
	if negated {
		node = ast.NewNot(tok, node)
	}
	return node, nil
}

func (p *Parser) subparse(endTokens []token.Token) ([]ast.Node, error) {
	body := []ast.Node{}
	dataBuffer := []ast.Expression{}
	addData := func(n ast.Expression) {
		dataBuffer = append(dataBuffer, n)
	}
	flushData := func() {
		if len(dataBuffer) == 0 {
			return
		}
		output := ast.NewOutput(dataBuffer[0].Token(), dataBuffer)
		body = append(body, output)
		dataBuffer = []ast.Expression{}
	}

	if endTokens != nil {
		p.endTokenStack = append(p.endTokenStack, endTokens)
		defer func() {
			p.endTokenStack = p.endTokenStack[:len(p.endTokenStack)-1]
		}()
	}

	for !p.stream.Eos() {
		tok := p.stream.Current()
		switch tok.Type {
		case token.TOKEN_DATA:
			if tok.Literal != "" {
				addData(ast.NewTemplateData(tok, tok.Literal))
			}
			p.stream.Next()
		case token.TOKEN_VARIABLE_BEGIN:
			p.stream.Next()
			expr, err := p.parseTuple(false, true, nil, false, false)
			if err != nil {
				return nil, err
			}
			addData(expr)
			if _, err := p.stream.Expect1(token.TOKEN_VARIABLE_END); err != nil {
				return nil, err
			}
		case token.TOKEN_BLOCK_BEGIN:
			flushData()
			p.stream.Next()
			if len(endTokens) > 0 && p.stream.CurrentTestAny(endTokens...) {
				return body, nil
			}
			rv, err := p.parseStatement()
			if err != nil {
				return nil, err
			}
			body = append(body, rv)
			if _, err := p.stream.Expect1(token.TOKEN_BLOCK_END); err != nil {
				return nil, err
			}
		default:
			panic("internal parsing error")
		}
	}

	flushData()

	return body, nil
}
