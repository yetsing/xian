package lexer

import (
	"fmt"
	"strings"

	"github.com/yetsing/xian/internal/token"
)

type CodeTag struct {
	start     string
	end       string
	startType token.TokenType
	endType   token.TokenType
}

type TemplateLexerConfig struct {
	BlockStart    string
	BlockEnd      string
	VariableStart string
	VariableEnd   string
	CommentStart  string
	CommentEnd    string

	LineStatementPrefix string
	LineCommentPrefix   string

	TrimBlocks   bool
	LstripBlocks bool

	NewlineSequence     string
	KeepTrailingNewline bool
}

type TemplateLexer struct {
	input          string
	index          int
	code           string
	codeLexer      *Lexer
	codeTags       []CodeTag
	codeStartToken token.Token
	codeEndToken   token.Token

	config TemplateLexerConfig

	segmentIndex int
	segments     []token.Token

	positioner *token.LineColumnIndex

	// 下一个 token 的索引
	tokenIndex int
	// 已经读取过的 token 数组
	tokens []token.Token
}

func NewTemplateLexer(input string, config TemplateLexerConfig) *TemplateLexer {
	tl := &TemplateLexer{
		input:        input,
		index:        0,
		config:       config,
		segmentIndex: 0,
		segments:     []token.Token{},
		tokenIndex:   0,
		tokens:       []token.Token{},
		positioner:   nil,
	}

	tl.initCodeTag()
	tl.processInput()
	tl.splitSegments()

	return tl
}

// NextToken 获取下一个 token ，同时增加 token 索引
func (tl *TemplateLexer) NextToken() token.Token {
	tk := tl.getToken(tl.tokenIndex)
	tl.tokenIndex++
	return tk
}

// PeekToken 查看下一个 token ，不增加 token 索引
func (tl *TemplateLexer) PeekToken() token.Token {
	return tl.getToken(tl.tokenIndex)
}

func (tl *TemplateLexer) getToken(index int) token.Token {
	if index >= len(tl.tokens) {
		// 索引超出已读取范围，再次进行读取
		tk := tl.readToken()
		tl.tokens = append(tl.tokens, tk)
	}
	return tl.tokens[index]
}

func (tl *TemplateLexer) readToken() token.Token {
BEGIN:
	segment := tl.segments[tl.segmentIndex]
	switch segment.Type {
	case token.TOKEN_DATA:
		tl.segmentIndex++
		return segment
	case token.TOKEN_BLOCK_BEGIN,
		token.TOKEN_BLOCK_END,
		token.TOKEN_VARIABLE_BEGIN,
		token.TOKEN_VARIABLE_END,
		token.TOKEN_COMMENT_BEGIN,
		token.TOKEN_COMMENT_END:
		tl.segmentIndex++
		return segment
	case token.TOKEN_COMMENT:
		tl.segmentIndex++
		return segment
	case token.TOKEN_STRING:
		// 进入代码块，使用 codeLexer 进行词法分析
		if tl.codeLexer == nil {
			tl.codeLexer = NewLexer(segment.Literal)
		}
		tk := tl.codeLexer.NextToken()
		if tk.TypeIs(token.TOKEN_EOF) {
			// 代码块分析完毕，回到主流程
			tl.codeLexer = nil
			tl.segmentIndex++
			goto BEGIN
		}
		return token.Token{
			Type:    tk.Type,
			Literal: tk.Literal,
			Start:   tk.Start.Add(segment.Start.Line, segment.Start.Column),
			End:     tk.End.Add(segment.Start.Line, segment.Start.Column),
		}

	case token.TOKEN_EOF, token.TOKEN_ILLEGAL:
		return segment
	}
	panic(fmt.Sprintf("unexpected segment type: %s", segment.Type))
}

func (tl *TemplateLexer) getPos(index int) token.Position {
	return tl.positioner.MustGetLineColumn(index)
}

func (tl *TemplateLexer) initCodeTag() {
	tl.codeTags = []CodeTag{
		{start: tl.config.BlockStart, end: tl.config.BlockEnd, startType: token.TOKEN_BLOCK_BEGIN, endType: token.TOKEN_BLOCK_END},
		{start: tl.config.VariableStart, end: tl.config.VariableEnd, startType: token.TOKEN_VARIABLE_BEGIN, endType: token.TOKEN_VARIABLE_END},
		{start: tl.config.CommentStart, end: tl.config.CommentEnd, startType: token.TOKEN_COMMENT_BEGIN, endType: token.TOKEN_COMMENT_END},
	}
	if tl.config.LineStatementPrefix != "" {
		tl.codeTags = append(tl.codeTags, CodeTag{start: tl.config.LineStatementPrefix, end: "\n", startType: token.TOKEN_BLOCK_BEGIN, endType: token.TOKEN_BLOCK_END})
	}
	if tl.config.LineCommentPrefix != "" {
		tl.codeTags = append(tl.codeTags, CodeTag{start: tl.config.LineCommentPrefix, end: "\n", startType: token.TOKEN_COMMENT_BEGIN, endType: token.TOKEN_COMMENT_END})
	}
}

func (tl *TemplateLexer) processInput() {
	tl.input = normalizeAllNewlines(tl.input)
	if !tl.config.KeepTrailingNewline && tl.input != "" && strings.HasSuffix(tl.input, "\n") {
		tl.input = tl.input[:len(tl.input)-1]
	}
	tl.positioner = token.NewLineColumnIndex(tl.input)
}

func (tl *TemplateLexer) splitSegments() {
	index := 0
	found := true
	for found {
		matchCodeTag := CodeTag{}
		matchIndex := len(tl.input) + 8
		for _, codeTag := range tl.codeTags {
			startIndex := FindByteIndex(tl.input, codeTag.start, index)
			if startIndex == -1 {
				continue
			}
			if startIndex < matchIndex || (startIndex == matchIndex && len(codeTag.start) > len(matchCodeTag.start)) {
				matchCodeTag = codeTag
				matchIndex = startIndex
			}
		}

		if matchIndex > len(tl.input) {
			// not found any code tag
			break
		}

		startIndex := matchIndex
		endIndex := FindByteIndex(tl.input, matchCodeTag.end, startIndex)
		if endIndex == -1 {
			tl.segments = append(tl.segments, token.Token{
				Type:    token.TOKEN_ILLEGAL,
				Literal: fmt.Sprintf("unclosed code tag %s", matchCodeTag.start),
				Start:   tl.getPos(startIndex),
				End:     tl.getPos(startIndex),
			})
			return
		}

		dataToken := token.Token{
			Type:    token.TOKEN_DATA,
			Literal: tl.input[index:startIndex],
			Start:   tl.getPos(index),
			End:     tl.getPos(startIndex),
		}
		startLen := len(matchCodeTag.start)
		if strings.ContainsRune("+-", rune(tl.input[startIndex+startLen])) {
			startLen++
		}
		startToken := token.Token{
			Type:    matchCodeTag.startType,
			Literal: tl.input[startIndex : startIndex+startLen],
			Start:   tl.getPos(startIndex),
			End:     tl.getPos(startIndex + startLen),
		}
		endLen := len(matchCodeTag.end)
		if strings.ContainsRune("+-", rune(tl.input[endIndex-1])) {
			endLen++
			endIndex--
		}
		endToken := token.Token{
			Type:    matchCodeTag.endType,
			Literal: tl.input[endIndex : endIndex+endLen],
			Start:   tl.getPos(endIndex),
			End:     tl.getPos(endIndex + endLen),
		}
		codeToken := token.Token{
			Type:    token.TOKEN_STRING,
			Literal: tl.input[startIndex+startLen : endIndex],
			Start:   tl.getPos(startIndex + startLen),
			End:     tl.getPos(endIndex),
		}
		if matchCodeTag.startType == token.TOKEN_COMMENT_BEGIN {
			codeToken.Type = token.TOKEN_COMMENT
		}
		if startIndex > index {
			tl.segments = append(tl.segments, dataToken)
		}
		tl.segments = append(tl.segments, startToken, codeToken, endToken)

		index = endIndex + endLen
	}

	if index < len(tl.input) {
		tl.segments = append(tl.segments, token.Token{
			Type:    token.TOKEN_DATA,
			Literal: tl.input[index:],
			Start:   tl.getPos(index),
			End:     tl.getPos(len(tl.input)),
		})
	}

	tl.segments = append(tl.segments, token.Token{
		Type:    token.TOKEN_EOF,
		Literal: "",
		Start:   tl.getPos(len(tl.input)),
		End:     tl.getPos(len(tl.input)),
	})
}

// normalizeAllNewlines handles \r\n and legacy \r, turning both into \n
func normalizeAllNewlines(s string) string {
	replacer := strings.NewReplacer("\r\n", "\n", "\r", "\n")
	return replacer.Replace(s)
}

// FindByteIndex mimics Python's str.find(sub, start) but returns the BYTE index.
// 'start' must be a valid byte index and UTF-8 boundary.
// Returns -1 if the substring is not found.
func FindByteIndex(s, sub string, start int) int {
	// Guard against out-of-bounds start indices
	if start < 0 {
		start = 0
	}
	if start >= len(s) {
		return -1
	}

	// Slice from the start byte and find the substring
	byteIdx := strings.Index(s[start:], sub)
	if byteIdx == -1 {
		return -1
	}

	// The returned index is relative to the slice,
	// so add the 'start' offset to get the absolute byte index.
	return start + byteIdx
}
