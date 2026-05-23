package lexer

import (
	"fmt"

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

func (ts *TokenStream) Next() token.Token {
	ts.current = ts.get(ts.tokenIndex)
	ts.tokenIndex++
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

func (ts *TokenStream) Current() token.Token {
	return ts.current
}

func (ts *TokenStream) Expect(tt token.TokenType) error {
	if ts.current.TypeIs(tt) {
		ts.Next()
		return nil
	}
	if ts.current.TypeIs(token.TOKEN_ILLEGAL) {
		return templateerror.NewTemplateSyntaxError(
			ts.current.Literal,
			ts.current.Start.Line,
			ts.name,
			ts.filename,
		)
	}
	return templateerror.NewTemplateSyntaxError(
		fmt.Sprintf("expected token %s, got %s", tt, ts.current.Type),
		ts.current.Start.Line,
		ts.name,
		ts.filename,
	)
}

func (ts *TokenStream) Eos() bool {
	return ts.current.TypeIs(token.TOKEN_EOF)
}

func (ts *TokenStream) get(index int) token.Token {
	for index >= len(ts.tokens) {
		tk := ts.lexer.NextToken()
		ts.tokens = append(ts.tokens, tk)
	}
	return ts.tokens[index]
}
