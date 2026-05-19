package lexer

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/yetsing/xian/internal/token"
)

type LexerState struct {
	index    int
	ch       rune
	position token.Position
}

type Lexer struct {
	// 输入文本
	input string
	// unicode 列表
	ucodes []rune

	// ch 在 ucodes 里面的下标
	index int
	// current char
	ch rune
	// ch 所在的行列(从 0 开始)
	position token.Position

	// 标记索引和位置，方便计算 Token 的 start end
	markIndex    int
	markPosition token.Position

	// token 数组和索引，用来支持回溯
	// 下一个 token 的索引
	tokenIndex int
	tokens     []token.Token
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, index: -1}
	l.init()
	l.readChar()
	return l
}

func (l *Lexer) GetLines() []string {
	return strings.Split(l.input, "\n")
}

// Dump 读取当前 token 索引
func (l *Lexer) Dump() int {
	return l.tokenIndex
}

// Restore 恢复至指定 token 索引
// 与 Dump 配合使用，可以在读取 token 之后进行撤销
// 具体使用例子可看 parser/parser.go
func (l *Lexer) Restore(index int) {
	l.tokenIndex = index
}

// NextToken 获取下一个 token ，同时增加 token 索引
func (l *Lexer) NextToken() token.Token {
	tk := l.getToken(l.tokenIndex)
	l.tokenIndex++
	return tk
}

// PeekToken 查看下一个 token ，不增加 token 索引
func (l *Lexer) PeekToken() token.Token {
	return l.getToken(l.tokenIndex)
}

func (l *Lexer) getToken(index int) token.Token {
	//    索引超出已读取范围，再次进行读取
	if index >= len(l.tokens) {
		tk := l.readToken()
		l.tokens = append(l.tokens, tk)
	}
	//    索引没有超出已读取范围，直接取之前已经读取过的
	return l.tokens[index]
}

func (l *Lexer) readToken() token.Token {
	var ttype token.TokenType

	l.skipWhitespace()

	l.mark()

	switch l.ch {
	case '=':
		// "==" 相等符号
		if l.peekCharIs('=') {
			l.advance(1)
			ttype = token.TOKEN_EQ
		} else {
			ttype = token.TOKEN_ASSIGN
		}
	case '+':
		ttype = token.TOKEN_ADD
	case '-':
		ttype = token.TOKEN_SUB
	case '!':
		// "!=" 不相等符号
		if l.peekCharIs('=') {
			l.advance(1)
			ttype = token.TOKEN_NE
		} else {
			l.advance(1)
			tok := l.buildToken(token.TOKEN_ILLEGAL)
			tok.Literal = "invalid char !"
			return tok
		}
	case '*':
		if l.peekCharIs('*') {
			l.readChar()
			ttype = token.TOKEN_POW
		} else {
			ttype = token.TOKEN_MUL
		}
	case '/':
		if l.peekCharIs('/') {
			l.readChar()
			ttype = token.TOKEN_FLOORDIV
		} else {
			ttype = token.TOKEN_DIV
		}
	case '%':
		ttype = token.TOKEN_MOD
	case '<':
		if l.peekCharIs('=') {
			l.readChar()
			ttype = token.TOKEN_LTEQ
		} else {
			ttype = token.TOKEN_LT
		}
	case '>':
		if l.peekCharIs('=') {
			l.readChar()
			ttype = token.TOKEN_GTEQ
		} else {
			ttype = token.TOKEN_GT
		}
	case '~':
		ttype = token.TOKEN_TILDE
	case '|':
		ttype = token.TOKEN_PIPE
	case ';':
		ttype = token.TOKEN_SEMICOLON
	case ':':
		ttype = token.TOKEN_COLON
	case ',':
		ttype = token.TOKEN_COMMA
	case '{':
		ttype = token.TOKEN_LBRACE
	case '}':
		ttype = token.TOKEN_RBRACE
	case '(':
		ttype = token.TOKEN_LPAREN
	case ')':
		ttype = token.TOKEN_RPAREN
	case '"':
		return l.readString(l.ch)
	case '\'':
		return l.readString(l.ch)
	case '`':
		return l.readRawString()
	case '[':
		ttype = token.TOKEN_LBRACKET
	case ']':
		ttype = token.TOKEN_RBRACKET
	case '.':
		ttype = token.TOKEN_DOT
	case 0:
		ttype = token.TOKEN_EOF
	default:
		if isIdentifierStart(l.ch) {
			return l.readIdentifier()
		} else if isDigit(l.ch) {
			return l.readNumber()
		} else {
			ttype = token.TOKEN_ILLEGAL
			ch := l.ch
			l.readChar()
			tok := l.buildToken(token.TOKEN_ILLEGAL)
			tok.Literal = "invalid char " + string(ch)
			return tok
		}
	}

	l.readChar()
	return l.buildToken(ttype)
}

func (l *Lexer) init() {
	l.ucodes = []rune(l.input)
	l.position.Line = 0
	l.position.Column = -1
	l.tokenIndex = 0
}

func (l *Lexer) dumpState() LexerState {
	return LexerState{
		index:    l.index,
		ch:       l.ch,
		position: l.position,
	}
}

func (l *Lexer) restoreState(state LexerState) {
	l.index = state.index
	l.ch = state.ch
	l.position = state.position
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readChar()
	}
	//for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
	//	l.readChar()
	//}
}

func (l *Lexer) readChar() {
	l.index++
	if l.index >= len(l.ucodes) {
		l.ch = 0
	} else {
		if l.ch == '\n' {
			l.position.Line++
			l.position.Column = -1
		}
		l.ch = l.ucodes[l.index]
		l.position.Column++
	}
}

func (l *Lexer) advance(n int) {
	for range n {
		l.readChar()
	}
}

func (l *Lexer) getString(n int) (string, error) {
	if l.index+n > len(l.ucodes) {
		return string(l.ucodes[l.index:len(l.ucodes)]), errors.New("not enough char")
	}
	s := string(l.ucodes[l.index : l.index+n])
	return s, nil
}

func (l *Lexer) peekCharIs(ch rune) bool {
	nextIndex := l.index + 1
	if nextIndex >= len(l.ucodes) {
		return 0 == ch
	} else {
		return l.ucodes[nextIndex] == ch
	}
}

// 标记一个位置
func (l *Lexer) mark() {
	l.markIndex = l.index
	l.markPosition.Line = l.position.Line
	l.markPosition.Column = l.position.Column
}

func (l *Lexer) buildToken(ttype token.TokenType) token.Token {
	start := l.markPosition
	end := l.position
	startIndex := l.markIndex
	endIndex := l.index
	switch ttype {
	// case token.TOKEN_STRING:
	// 移除首尾的引号
	// startIndex++
	// endIndex--
	case token.TOKEN_EOF:
		//    EOF 时， endIndex 已经超出范围
		//    确保 EOF token 的 Literal 为空字符串
		return token.Token{
			Type:    token.TOKEN_EOF,
			Literal: "",
			Start:   start,
			End:     end,
		}
	}
	tok := token.Token{Type: ttype, Literal: string(l.ucodes[startIndex:endIndex]), Start: start, End: end}
	return tok
}

func (l *Lexer) readIdentifier() token.Token {
	for isIdentifierContinue(l.ch) {
		l.readChar()
	}
	tok := l.buildToken(token.TOKEN_NAME)
	tok.Type = token.LookupIdent(tok.Literal)
	return tok
}

func (l *Lexer) readNumber() token.Token {
	state := l.dumpState()
	if tok, err := l.tryFloatnumber1(); err == nil {
		return tok
	}
	l.restoreState(state)
	if tok, err := l.tryFloatnumber2(); err == nil {
		return tok
	}
	l.restoreState(state)
	if tok, err := l.tryDecinteger(); err == nil {
		return tok
	}
	l.restoreState(state)
	if tok, err := l.tryBininteger(); err == nil {
		return tok
	}
	l.restoreState(state)
	if tok, err := l.tryOctinteger(); err == nil {
		return tok
	}
	l.restoreState(state)
	if tok, err := l.tryHexinteger(); err == nil {
		return tok
	}
	l.restoreState(state)
	if tok, err := l.tryZerointeger(); err == nil {
		return tok
	}
	l.restoreState(state)

	tok := l.buildToken(token.TOKEN_ILLEGAL)
	tok.Literal = "invalid number"
	return tok
}

func (l *Lexer) passDigitpart() error {
	if !isDigit(l.ch) {
		return errors.New("digit part expected")
	}
	for {
		if l.ch == '_' {
			l.readChar()
			continue
		}
		if !isDigit(l.ch) {
			break
		}
		l.readChar()
	}
	return nil
}

// floatnumber1: digitpart "." [digitpart] [exponent]
func (l *Lexer) tryFloatnumber1() (token.Token, error) {
	var err error

	err = l.passDigitpart()
	if err != nil {
		return token.Token{}, err
	}
	if l.ch != '.' {
		return token.Token{}, errors.New("dot expected")
	}
	l.readChar()
	if isDigit(l.ch) {
		err = l.passDigitpart()
		if err != nil {
			return token.Token{}, err
		}
	}
	if l.ch == 'e' || l.ch == 'E' {
		l.readChar()
		if l.ch == '+' || l.ch == '-' {
			l.readChar()
		}
		err = l.passDigitpart()
		if err != nil {
			return token.Token{}, err
		}
	}
	tok := l.buildToken(token.TOKEN_FLOAT)
	return tok, nil
}

// floatnumber2: digitpart exponent
func (l *Lexer) tryFloatnumber2() (token.Token, error) {
	var err error

	err = l.passDigitpart()
	if err != nil {
		return token.Token{}, err
	}
	if l.ch != 'e' && l.ch != 'E' {
		return token.Token{}, errors.New("exponent expected")
	}
	l.readChar()
	if l.ch == '+' || l.ch == '-' {
		l.readChar()
	}
	err = l.passDigitpart()
	if err != nil {
		return token.Token{}, err
	}
	tok := l.buildToken(token.TOKEN_FLOAT)
	return tok, nil
}

// decinteger: nonzerodigit (["_"] digit)*
func (l *Lexer) tryDecinteger() (token.Token, error) {
	if !(isDigit(l.ch) && l.ch != '0') {
		return token.Token{}, errors.New("non-zero digit expected")
	}

	for {
		if l.ch == '_' {
			l.readChar()
			continue
		}
		if !isDigit(l.ch) {
			break
		}
		l.readChar()
	}
	tok := l.buildToken(token.TOKEN_INTEGER)
	return tok, nil
}

// bininteger:   "0" ("b" | "B") (["_"] bindigit)+
// octinteger:   "0" ("o" | "O") (["_"] octdigit)+
// hexinteger:   "0" ("x" | "X") (["_"] hexdigit)+
func (l *Lexer) tryNonDecinteger(letters string, check func(rune) bool) (token.Token, error) {
	if l.ch != '0' {
		return token.Token{}, errors.New("invalid prefix")
	}
	l.readChar()
	if !strings.ContainsRune(letters, l.ch) {
		return token.Token{}, errors.New("invalid prefix")
	}
	l.readChar()
	for {
		if l.ch == '_' {
			l.readChar()
			continue
		}
		if !check(l.ch) {
			break
		}
		l.readChar()
	}
	tok := l.buildToken(token.TOKEN_INTEGER)
	return tok, nil
}

// bininteger: "0" ("b" | "B") (["_"] bindigit)+
func (l *Lexer) tryBininteger() (token.Token, error) {
	return l.tryNonDecinteger("bB", isBindigit)
}

// octinteger: "0" ("o" | "O") (["_"] octdigit)+
func (l *Lexer) tryOctinteger() (token.Token, error) {
	return l.tryNonDecinteger("oO", isOctDigit)
}

// hexinteger: "0" ("x" | "X") (["_"] hexdigit)+
func (l *Lexer) tryHexinteger() (token.Token, error) {
	return l.tryNonDecinteger("xX", isHexdigit)
}

// zerointeger:  "0"+ (["_"] "0")*
func (l *Lexer) tryZerointeger() (token.Token, error) {
	if l.ch != '0' {
		return token.Token{}, errors.New("zero expected")
	}
	for {
		if l.ch == '_' {
			l.readChar()
			continue
		}
		if l.ch != '0' {
			break
		}
		l.readChar()
	}
	tok := l.buildToken(token.TOKEN_INTEGER)
	return tok, nil
}

var escapeMap = map[rune]rune{
	'\\': '\\',
	'\'': '\'',
	'"':  '"',
	'a':  '\a',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
	'v':  '\v',
}

func parseRune(s string, base int, bitSize int) (rune, error) {
	n, err := strconv.ParseUint(s, base, bitSize)
	if err != nil {
		return 0, err
	}
	r := rune(n)
	return r, nil
}

func (l *Lexer) readString(end rune) token.Token {
	var sb strings.Builder
	// 跳过开始的引号
	l.readChar()
	for {
		// 参考 Python 的转义字符 https://docs.python.org/3/reference/lexical_analysis.html#string-and-bytes-literals
		// 处理转义字符
		ch := l.ch
		if ch == '\\' {
			l.readChar()
			if actual, ok := escapeMap[l.ch]; ok {
				sb.WriteRune(actual)
				l.readChar()
				continue
			}
			// 解析 Unicode 转义字符
			var ucode rune
			var codeLen int
			switch l.ch {
			case 'x':
				// 格式为 "\xhh" h 代表十六进制字符
				// 跳过 'x' 字符
				l.readChar()
				codeLen = 2
				s, err := l.getString(codeLen)
				if err != nil {
					tok := l.buildToken(token.TOKEN_ILLEGAL)
					tok.Literal = "illegal escape sequence"
					return tok
				}
				ucode, err = parseRune(s, 16, 8)
				if err != nil {
					tok := l.buildToken(token.TOKEN_ILLEGAL)
					tok.Literal = "illegal escape sequence"
					return tok
				}
			case 'u':
				// 格式为 "\uhhhh" h 代表十六进制字符
				// 跳过 'u' 字符
				l.readChar()
				codeLen = 4
				s, err := l.getString(codeLen)
				if err != nil {
					tok := l.buildToken(token.TOKEN_ILLEGAL)
					tok.Literal = "illegal escape sequence"
					return tok
				}
				ucode, err = parseRune(s, 16, 16)
				if err != nil {
					tok := l.buildToken(token.TOKEN_ILLEGAL)
					tok.Literal = "illegal escape sequence"
					return tok
				}
			case 'U':
				// 格式为 "\Uhhhhhhhh" h 代表十六进制字符
				// 跳过 'u' 字符
				l.readChar()
				codeLen = 8
				s, err := l.getString(codeLen)
				if err != nil {
					tok := l.buildToken(token.TOKEN_ILLEGAL)
					tok.Literal = "illegal escape sequence"
					return tok
				}
				ucode, err = parseRune(s, 16, 32)
				if err != nil {
					tok := l.buildToken(token.TOKEN_ILLEGAL)
					tok.Literal = "illegal escape sequence"
					return tok
				}
			case '0', '1', '2', '3', '4', '5', '6', '7':
				// 格式为 "\ooo" o 代表八进制字符，最大为 "\377" (255)
				codeLen = 3
				s, err := l.getString(codeLen)
				if err != nil {
					tok := l.buildToken(token.TOKEN_ILLEGAL)
					tok.Literal = "illegal escape sequence"
					return tok
				}
				ucode, err = parseRune(s, 8, 8)
				if err != nil {
					tok := l.buildToken(token.TOKEN_ILLEGAL)
					tok.Literal = "illegal escape sequence"
					return tok
				}
			default:
				// 非法转义字符
				tok := l.buildToken(token.TOKEN_ILLEGAL)
				tok.Literal = "illegal escape sequence"
				return tok
			}
			if !utf8.ValidRune(ucode) {
				tok := l.buildToken(token.TOKEN_ILLEGAL)
				tok.Literal = "escape sequence is invalid Unicode code point"
				return tok
			}
			sb.WriteRune(ucode)
			l.advance(codeLen)
			continue
		}

		if l.ch == end {
			break
		}
		if l.ch == 0 || l.ch == '\n' {
			tok := l.buildToken(token.TOKEN_ILLEGAL)
			tok.Literal = "string literal not terminated"
			return tok
		}
		sb.WriteRune(l.ch)
		l.readChar()
	}
	// 跳过末尾的引号
	l.readChar()
	tok := l.buildToken(token.TOKEN_STRING)
	tok.Literal = sb.String()
	return tok
}

func (l *Lexer) readRawString() token.Token {
	// 跳过开始的引号
	l.readChar()
	for {
		if l.ch == 0 {
			tok := l.buildToken(token.TOKEN_ILLEGAL)
			tok.Literal = "string literal not terminated"
			return tok
		}
		if l.ch == '`' {
			break
		}
		l.readChar()
	}
	// 跳过结尾的引号
	l.readChar()
	return l.buildToken(token.TOKEN_STRING)
}

// 参考 Python 的规则 https://docs.python.org/3/reference/lexical_analysis.html#identifiers
var idStartCategorys = []*unicode.RangeTable{
	unicode.Lu,
	unicode.Ll,
	unicode.Lt,
	unicode.Lm,
	unicode.Lo,
	unicode.Nl,
	unicode.Other_ID_Start,
}
var idContinueCategorys = []*unicode.RangeTable{
	unicode.Lu,
	unicode.Ll,
	unicode.Lt,
	unicode.Lm,
	unicode.Lo,
	unicode.Nl,
	unicode.Other_ID_Start,

	unicode.Nd,
	unicode.Pc,
	unicode.Mn,
	unicode.Mc,
	unicode.Other_ID_Continue,
}

func isIdentifierStart(ch rune) bool {
	if ch == '_' || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
		// fast path for common ASCII characters
		return true
	}
	for _, table := range idStartCategorys {
		if unicode.Is(table, ch) {
			return true
		}
	}
	return false
}

func isIdentifierContinue(ch rune) bool {
	if ch == '_' || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') {
		// fast path for common ASCII characters
		return true
	}
	for _, table := range idContinueCategorys {
		if unicode.Is(table, ch) {
			return true
		}
	}
	return false
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func isBindigit(ch rune) bool {
	return ch == '0' || ch == '1'
}

func isOctDigit(ch rune) bool {
	return '0' <= ch && ch <= '7'
}

func isHexdigit(ch rune) bool {
	return ('0' <= ch && ch <= '9') || ('a' <= ch && ch <= 'f') || ('A' <= ch && ch <= 'F')
}
