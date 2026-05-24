package lexer

import (
	"fmt"
	"slices"

	"github.com/yetsing/xian/templateerror"
	"github.com/yetsing/xian/token"
)

type LexerInterface interface {
	NextToken() token.Token
}

type TokenStream struct {
	lexer    LexerInterface
	name     string
	filename string

	current token.Token

	tokenIndex int
	tokens     []token.Token
}

func NewTokenStream(lexer LexerInterface, name string, filename string) *TokenStream {
	ts := &TokenStream{
		lexer:    lexer,
		name:     name,
		filename: filename,
		current: token.Token{
			Type:    token.TOKEN_INITIAL,
			Literal: "",
			Start:   token.Position{Line: 1, Column: 1},
			End:     token.Position{Line: 1, Column: 1},
		},
		tokenIndex: 0,
		tokens:     []token.Token{},
	}
	ts.Next() // 预加载第一个 token
	return ts
}

// Next advances the token stream and returns the current token before advancing.
func (ts *TokenStream) Next() token.Token {
	old := ts.current
	ts.current = ts.get(ts.tokenIndex)
	ts.tokenIndex++
	return old
}

func (ts *TokenStream) Current() token.Token {
	return ts.current
}

func (ts *TokenStream) Peek() token.Token {
	return ts.get(ts.tokenIndex)
}

func (ts *TokenStream) Skip() {
	ts.Next()
}

func (ts *TokenStream) SkipN(n int) {
	for range n {
		ts.Next()
	}
}

func (ts *TokenStream) NextIf(tok token.Token) (token.Token, bool) {
	if ts.current.Test(tok) {
		return ts.Next(), true
	}
	return token.Token{}, false
}

func (ts *TokenStream) SkipIf1(ttype token.TokenType) bool {
	_, ok := ts.NextIf(token.Token{Type: ttype})
	return ok
}

func (ts *TokenStream) SkipIf2(ttype token.TokenType, literal string) bool {
	_, ok := ts.NextIf(token.Token{Type: ttype, Literal: literal})
	return ok
}

func (ts *TokenStream) Expect1(ttype token.TokenType) (token.Token, error) {
	return ts.Expect2(ttype, "")
}

func (ts *TokenStream) Expect2(ttype token.TokenType, literal string) (token.Token, error) {
	if ts.current.Test(token.Token{Type: ttype, Literal: literal}) {
		return ts.Next(), nil
	}
	if ts.current.Type == token.TOKEN_ILLEGAL {
		return token.Token{}, templateerror.NewTemplateSyntaxError(
			ts.current.Literal,
			ts.current.Start.Line,
			ts.name,
			ts.filename,
		)
	}
	if ts.current.Type == token.TOKEN_EOF {
		return token.Token{}, templateerror.NewTemplateSyntaxError(
			fmt.Sprintf("unexpected end of template, expected %q", ttype),
			ts.current.Start.Line,
			ts.name,
			ts.filename,
		)
	}
	return token.Token{}, templateerror.NewTemplateSyntaxError(
		fmt.Sprintf("expected token %q, got %q", ttype, ts.current.Type),
		ts.current.Start.Line,
		ts.name,
		ts.filename,
	)
}

func (ts *TokenStream) CurrentEqual1(ttype token.TokenType) bool {
	return ts.current.Type == ttype
}

func (ts *TokenStream) CurrentNotEqual1(ttype token.TokenType) bool {
	return ts.current.Type != ttype
}

func (ts *TokenStream) CurrentEqual2(ttype token.TokenType, literal string) bool {
	return ts.current.Type == ttype && ts.current.Literal == literal
}

func (ts *TokenStream) CurrentIn1(ttypes ...token.TokenType) bool {
	return slices.Contains(ttypes, ts.current.Type)
}

func (ts *TokenStream) CurrentNotIn1(ttypes ...token.TokenType) bool {
	return !slices.Contains(ttypes, ts.current.Type)
}

func (ts *TokenStream) CurrentLiteralIn(literals ...string) bool {
	return slices.Contains(literals, ts.current.Literal)
}

func (ts *TokenStream) CurrentTestAny(toks ...token.Token) bool {
	return ts.current.TestAny(toks...)
}

func (ts *TokenStream) PeekEqual1(ttype token.TokenType) bool {
	return ts.Peek().Type == ttype
}

func (ts *TokenStream) PeekEqual2(ttype token.TokenType, literal string) bool {
	peeked := ts.Peek()
	return peeked.Type == ttype && peeked.Literal == literal
}

func (ts *TokenStream) Eos() bool {
	return ts.current.Type == token.TOKEN_EOF
}

func (ts *TokenStream) get(index int) token.Token {
	for index >= len(ts.tokens) {
		tk := ts.lexer.NextToken()
		ts.tokens = append(ts.tokens, tk)
	}
	return ts.tokens[index]
}
