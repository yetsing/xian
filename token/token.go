package token

import "slices"

type TokenType string

const (
	TOKEN_ILLEGAL = "illegal"

	TOKEN_ADD                 = "add"       // "+"
	TOKEN_ASSIGN              = "assign"    // "="
	TOKEN_COLON               = "colon"     // ":"
	TOKEN_COMMA               = "comma"     // ","
	TOKEN_DIV                 = "div"       // "/"
	TOKEN_DOT                 = "dot"       // "."
	TOKEN_EQ                  = "eq"        // "=="
	TOKEN_FLOORDIV            = "floordiv"  // "//"
	TOKEN_GT                  = "gt"        // ">"
	TOKEN_GTEQ                = "gteq"      // ">="
	TOKEN_LBRACE              = "lbrace"    // "{"
	TOKEN_LBRACKET            = "lbracket"  // "["
	TOKEN_LPAREN              = "lparen"    // "("
	TOKEN_LT                  = "lt"        // "<"
	TOKEN_LTEQ                = "lteq"      // "<="
	TOKEN_MOD                 = "mod"       // "%"
	TOKEN_MUL                 = "mul"       // "*"
	TOKEN_NE                  = "ne"        // "!="
	TOKEN_PIPE                = "pipe"      // "|"
	TOKEN_POW                 = "pow"       // "**"
	TOKEN_RBRACE              = "rbrace"    // "}"
	TOKEN_RBRACKET            = "rbracket"  // "]"
	TOKEN_RPAREN              = "rparen"    // ")"
	TOKEN_SEMICOLON           = "semicolon" // ";"
	TOKEN_SUB                 = "sub"       // "-"
	TOKEN_TILDE               = "tilde"     // "~"
	TOKEN_WHITESPACE          = "whitespace"
	TOKEN_FLOAT               = "float"
	TOKEN_INTEGER             = "integer"
	TOKEN_NAME                = "name"
	TOKEN_STRING              = "string"
	TOKEN_OPERATOR            = "operator"
	TOKEN_BLOCK_BEGIN         = "block_begin"
	TOKEN_BLOCK_END           = "block_end"
	TOKEN_VARIABLE_BEGIN      = "variable_begin"
	TOKEN_VARIABLE_END        = "variable_end"
	TOKEN_RAW_BEGIN           = "raw_begin"
	TOKEN_RAW_END             = "raw_end"
	TOKEN_COMMENT_BEGIN       = "comment_begin"
	TOKEN_COMMENT_END         = "comment_end"
	TOKEN_COMMENT             = "comment"
	TOKEN_LINESTATEMENT_BEGIN = "linestatement_begin"
	TOKEN_LINESTATEMENT_END   = "linestatement_end"
	TOKEN_LINECOMMENT_BEGIN   = "linecomment_begin"
	TOKEN_LINECOMMENT_END     = "linecomment_end"
	TOKEN_LINECOMMENT         = "linecomment"
	TOKEN_DATA                = "data"
	TOKEN_INITIAL             = "initial"
	TOKEN_EOF                 = "eof"
)

type Token struct {
	Type    TokenType
	Literal string
	Start   Position
	End     Position
}

func (t *Token) TypeIs(ttype TokenType) bool {
	return t.Type == ttype
}

func (t *Token) TypeNotIs(ttype TokenType) bool {
	return t.Type != ttype
}

func (t *Token) TypeIn(ttypes ...TokenType) bool {
	return slices.ContainsFunc(ttypes, t.TypeIs)
}

func (t *Token) LiteralIs(s string) bool {
	return t.Literal == s
}

func (t *Token) Is2(ttype TokenType, literal string) bool {
	return t.TypeIs(ttype) && t.LiteralIs(literal)
}

func (t *Token) IsOperator() bool {
	_, ok := operators[t.Literal]
	return ok
}

var operators = map[string]TokenType{
	"+":  TOKEN_ADD,
	"-":  TOKEN_SUB,
	"/":  TOKEN_DIV,
	"//": TOKEN_FLOORDIV,
	"*":  TOKEN_MUL,
	"%":  TOKEN_MOD,
	"**": TOKEN_POW,
	"~":  TOKEN_TILDE,
	"[":  TOKEN_LBRACKET,
	"]":  TOKEN_RBRACKET,
	"(":  TOKEN_LPAREN,
	")":  TOKEN_RPAREN,
	"{":  TOKEN_LBRACE,
	"}":  TOKEN_RBRACE,
	"==": TOKEN_EQ,
	"!=": TOKEN_NE,
	">":  TOKEN_GT,
	">=": TOKEN_GTEQ,
	"<":  TOKEN_LT,
	"<=": TOKEN_LTEQ,
	"=":  TOKEN_ASSIGN,
	".":  TOKEN_DOT,
	":":  TOKEN_COLON,
	"|":  TOKEN_PIPE,
	",":  TOKEN_COMMA,
	";":  TOKEN_SEMICOLON,
}
