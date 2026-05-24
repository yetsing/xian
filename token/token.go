package token

import (
	"slices"
)

type TokenType string

const (
	TOKEN_ILLEGAL             TokenType = "illegal"
	TOKEN_ADD                 TokenType = "add"       // "+"
	TOKEN_ASSIGN              TokenType = "assign"    // "="
	TOKEN_COLON               TokenType = "colon"     // ":"
	TOKEN_COMMA               TokenType = "comma"     // ","
	TOKEN_DIV                 TokenType = "div"       // "/"
	TOKEN_DOT                 TokenType = "dot"       // "."
	TOKEN_EQ                  TokenType = "eq"        // "=="
	TOKEN_FLOORDIV            TokenType = "floordiv"  // "//"
	TOKEN_GT                  TokenType = "gt"        // ">"
	TOKEN_GTEQ                TokenType = "gteq"      // ">="
	TOKEN_LBRACE              TokenType = "lbrace"    // "{"
	TOKEN_LBRACKET            TokenType = "lbracket"  // "["
	TOKEN_LPAREN              TokenType = "lparen"    // "("
	TOKEN_LT                  TokenType = "lt"        // "<"
	TOKEN_LTEQ                TokenType = "lteq"      // "<="
	TOKEN_MOD                 TokenType = "mod"       // "%"
	TOKEN_MUL                 TokenType = "mul"       // "*"
	TOKEN_NE                  TokenType = "ne"        // "!="
	TOKEN_PIPE                TokenType = "pipe"      // "|"
	TOKEN_POW                 TokenType = "pow"       // "**"
	TOKEN_RBRACE              TokenType = "rbrace"    // "}"
	TOKEN_RBRACKET            TokenType = "rbracket"  // "]"
	TOKEN_RPAREN              TokenType = "rparen"    // ")"
	TOKEN_SEMICOLON           TokenType = "semicolon" // ";"
	TOKEN_SUB                 TokenType = "sub"       // "-"
	TOKEN_TILDE               TokenType = "tilde"     // "~"
	TOKEN_WHITESPACE          TokenType = "whitespace"
	TOKEN_FLOAT               TokenType = "float"
	TOKEN_INTEGER             TokenType = "integer"
	TOKEN_NAME                TokenType = "name"
	TOKEN_STRING              TokenType = "string"
	TOKEN_OPERATOR            TokenType = "operator"
	TOKEN_BLOCK_BEGIN         TokenType = "block_begin"
	TOKEN_BLOCK_END           TokenType = "block_end"
	TOKEN_VARIABLE_BEGIN      TokenType = "variable_begin"
	TOKEN_VARIABLE_END        TokenType = "variable_end"
	TOKEN_RAW_BEGIN           TokenType = "raw_begin"
	TOKEN_RAW_END             TokenType = "raw_end"
	TOKEN_COMMENT_BEGIN       TokenType = "comment_begin"
	TOKEN_COMMENT_END         TokenType = "comment_end"
	TOKEN_COMMENT             TokenType = "comment"
	TOKEN_LINESTATEMENT_BEGIN TokenType = "linestatement_begin"
	TOKEN_LINESTATEMENT_END   TokenType = "linestatement_end"
	TOKEN_LINECOMMENT_BEGIN   TokenType = "linecomment_begin"
	TOKEN_LINECOMMENT_END     TokenType = "linecomment_end"
	TOKEN_LINECOMMENT         TokenType = "linecomment"
	TOKEN_DATA                TokenType = "data"
	TOKEN_INITIAL             TokenType = "initial"
	TOKEN_EOF                 TokenType = "eof"
)

type Token struct {
	Type    TokenType
	Literal string
	Start   Position
	End     Position
}

func (t *Token) Equal2(ty TokenType, lit string) bool {
	return t.Type == ty && t.Literal == lit
}

func (t *Token) TypeIn(types ...TokenType) bool {
	return slices.Contains(types, t.Type)
}

func (t *Token) TypeNotIn(types ...TokenType) bool {
	return !slices.Contains(types, t.Type)
}

func (t *Token) LiteralIn(literals ...string) bool {
	return slices.Contains(literals, t.Literal)
}

func (t *Token) LiteralNotIn(literals ...string) bool {
	return !slices.Contains(literals, t.Literal)
}

// Test checks if the token matches the given token.
// If the given token's Literal is empty, only the Type is compared.
func (t *Token) Test(tk Token) bool {
	if tk.Literal == "" {
		return t.Type == tk.Type
	}

	return t.Type == tk.Type && t.Literal == tk.Literal
}

func (t *Token) TestAny(tks ...Token) bool {
	return slices.ContainsFunc(tks, t.Test)
}

func (t *Token) Empty() bool {
	return t.Type == "" && t.Literal == ""
}
