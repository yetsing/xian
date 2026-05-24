package lexer

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/yetsing/xian/config"
	"github.com/yetsing/xian/token"
)

type CodeTag struct {
	start     string
	end       string
	startType token.TokenType
	endType   token.TokenType
}

type TemplateLexer struct {
	input          string
	index          int
	code           string
	codeLexer      *Lexer
	codeTags       []CodeTag
	codeStartToken token.Token
	codeEndToken   token.Token

	config config.Config

	segmentIndex int
	segments     []token.Token

	positioner *token.LineColumnIndex

	ignoredTokens []token.TokenType
	ignoreIfEmpty []token.TokenType
}

func NewTemplateLexer(input string, config *config.Config) *TemplateLexer {
	tl := &TemplateLexer{
		input:        input,
		index:        0,
		config:       *config,
		segmentIndex: 0,
		segments:     []token.Token{},
		positioner:   nil,
		ignoredTokens: []token.TokenType{
			token.TOKEN_COMMENT_BEGIN,
			token.TOKEN_COMMENT,
			token.TOKEN_COMMENT_END,
			token.TOKEN_LINECOMMENT_BEGIN,
			token.TOKEN_LINECOMMENT,
			token.TOKEN_LINECOMMENT_END,
		},
		ignoreIfEmpty: []token.TokenType{
			token.TOKEN_DATA,
			token.TOKEN_COMMENT,
			token.TOKEN_LINECOMMENT,
		},
	}

	tl.initCodeTag()
	tl.processInput()
	tl.splitSegments()

	return tl
}

// NextToken 获取下一个 token ，同时增加 token 索引
func (tl *TemplateLexer) NextToken() token.Token {
	tk := tl.readToken()
	ignored := false
	for tk.Type != token.TOKEN_EOF {
		tk, ignored = tl.convertToken(tk)
		if ignored {
			// 如果是需要忽略的 token ，继续读取下一个
			tk = tl.readToken()
			continue
		}
		break
	}
	return tk
}

func (tl *TemplateLexer) convertToken(tk token.Token) (token.Token, bool) {
	var err error

	if tk.TypeIn(tl.ignoredTokens...) {
		return tk, true
	}
	if tk.TypeIn(tl.ignoreIfEmpty...) && tk.Literal == "" {
		return tk, true
	}

	switch tk.Type {
	case token.TOKEN_LINESTATEMENT_BEGIN:
		tk.Type = token.TOKEN_BLOCK_BEGIN
	case token.TOKEN_LINESTATEMENT_END:
		tk.Type = token.TOKEN_BLOCK_END
	case token.TOKEN_RAW_BEGIN, token.TOKEN_RAW_END:
		return tk, true
	case token.TOKEN_DATA:
		tk.Literal = tl.normalizeNewlines(tk.Literal)
	case token.TOKEN_STRING:
		if tk.Literal == "" {
			break
		}
		// remove the quotes and unescape the string literal
		literal := tk.Literal[1 : len(tk.Literal)-1]
		if tk.Literal[0] == '\'' {
			// 转义双引号
			var sb strings.Builder
			backslash := false
			for _, ch := range literal {
				if ch == '"' && !backslash {
					sb.WriteRune('\\')
				}
				sb.WriteRune(ch)
				backslash = ch == '\\' && !backslash
			}
			literal = sb.String()
		}
		// Unquote 不支持多行双引号字符串，因此分行 Unquote
		lines := strings.Split(literal, "\n")
		for i, line := range lines {
			lines[i], err = strconv.Unquote(fmt.Sprintf(`"%s"`, line))
			if err != nil {
				tk.Type = token.TOKEN_ILLEGAL
				tk.Literal = fmt.Sprintf("string literal unquote error: %v", err)
				return tk, false
			}
		}
		// 跟 jinja2 保持一致，字符串中的换行会被转换成配置中的换行序列
		// 注意：替换的换行是文本本身的换行，而不是字符串字面量中的 \n
		tk.Literal = strings.Join(lines, tl.config.NewlineSequence)
	}
	return tk, false
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
		token.TOKEN_RAW_BEGIN,
		token.TOKEN_RAW_END,
		token.TOKEN_COMMENT_BEGIN,
		token.TOKEN_COMMENT_END,
		token.TOKEN_COMMENT,
		token.TOKEN_LINESTATEMENT_BEGIN,
		token.TOKEN_LINESTATEMENT_END,
		token.TOKEN_LINECOMMENT_BEGIN,
		token.TOKEN_LINECOMMENT_END,
		token.TOKEN_LINECOMMENT:
		tl.segmentIndex++
		return segment
	case token.TOKEN_STRING:
		// 进入代码块，使用 codeLexer 进行词法分析
		if tl.codeLexer == nil {
			tl.codeLexer = NewLexer(segment.Literal, true)
		}
		tk := tl.codeLexer.NextToken()
		if tk.Type == token.TOKEN_EOF {
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
		{start: tl.config.BlockStartString, end: tl.config.BlockEndString, startType: token.TOKEN_BLOCK_BEGIN, endType: token.TOKEN_BLOCK_END},
		{start: tl.config.VariableStartString, end: tl.config.VariableEndString, startType: token.TOKEN_VARIABLE_BEGIN, endType: token.TOKEN_VARIABLE_END},
		{start: tl.config.CommentStartString, end: tl.config.CommentEndString, startType: token.TOKEN_COMMENT_BEGIN, endType: token.TOKEN_COMMENT_END},
	}
	if tl.config.LineStatementPrefix != "" {
		tl.codeTags = append(tl.codeTags, CodeTag{start: tl.config.LineStatementPrefix, end: "\n", startType: token.TOKEN_LINESTATEMENT_BEGIN, endType: token.TOKEN_LINESTATEMENT_END})
	}
	if tl.config.LineCommentPrefix != "" {
		tl.codeTags = append(tl.codeTags, CodeTag{start: tl.config.LineCommentPrefix, end: "\n", startType: token.TOKEN_LINECOMMENT_BEGIN, endType: token.TOKEN_LINECOMMENT_END})
	}
}

func (tl *TemplateLexer) processInput() {
	tl.input = normalizeAllNewlines(tl.input)
	if !tl.config.KeepTrailingNewline && tl.input != "" && strings.HasSuffix(tl.input, "\n") {
		tl.input = tl.input[:len(tl.input)-1]
	}
	tl.positioner = token.NewLineColumnIndex(tl.input)
}

// splitSegments 将输入字符串分割成数据段和代码段，并存储在 tl.segments 中
func (tl *TemplateLexer) splitSegments() {
	index := 0
	inRaw := false
	rawStartIndex := -1
	for {
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
		endFindOffset := startIndex + len(matchCodeTag.start)
		endIndex := -1
		for endFindOffset <= len(tl.input) {
			idx := FindByteIndex(tl.input, matchCodeTag.end, endFindOffset)
			if idx == -1 {
				tl.segments = append(tl.segments, token.Token{
					Type:    token.TOKEN_ILLEGAL,
					Literal: fmt.Sprintf("unclosed code tag %s", matchCodeTag.start),
					Start:   tl.getPos(startIndex),
					End:     tl.getPos(startIndex),
				})
				return
			}
			if tl.isBalance(startIndex+len(matchCodeTag.start), idx) {
				endIndex = idx
				break
			}
			endFindOffset = idx + 1 // +1 to avoid finding the same end tag again
		}
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
		stripSign := byte(0)
		startLen := len(matchCodeTag.start)
		if tl.input[startIndex+startLen] == '-' || tl.input[startIndex+startLen] == '+' {
			stripSign = tl.input[startIndex+startLen]
			startLen++
		}
		startToken := token.Token{
			Type:    matchCodeTag.startType,
			Literal: tl.input[startIndex : startIndex+startLen],
			Start:   tl.getPos(startIndex),
			End:     tl.getPos(startIndex + startLen),
		}
		tl.lstripBlocks(&dataToken, &matchCodeTag, stripSign)
		stripSign = 0
		endLen := len(matchCodeTag.end)
		if tl.input[endIndex-1] == '-' || tl.input[endIndex-1] == '+' {
			stripSign = tl.input[endIndex-1]
			endLen++
			endIndex--
		}
		// 行为与 jinja2 保持一致，将空白字符放在 endToken
		endLen += tl.trimBlocks(stripSign, &matchCodeTag, endIndex+endLen)
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
		switch matchCodeTag.startType {
		case token.TOKEN_COMMENT_BEGIN:
			codeToken.Type = token.TOKEN_COMMENT
		case token.TOKEN_LINECOMMENT_BEGIN:
			codeToken.Type = token.TOKEN_LINECOMMENT
		}

		if strings.TrimSpace(codeToken.Literal) == "raw" {
			inRaw = true
			rawStartIndex = endIndex + endLen
			tl.segments = append(tl.segments,
				dataToken,
				token.Token{
					Type:    token.TOKEN_RAW_BEGIN,
					Literal: tl.input[startIndex : endIndex+endLen],
					Start:   startToken.Start,
					End:     endToken.End,
				})
		} else if strings.TrimSpace(codeToken.Literal) == "endraw" {
			if !inRaw {
				tl.segments = append(tl.segments, token.Token{
					Type:    token.TOKEN_ILLEGAL,
					Literal: "unexpected endraw tag",
					Start:   startToken.Start,
					End:     endToken.End,
				})
				return
			}
			dataToken := token.Token{
				Type:    token.TOKEN_DATA,
				Literal: tl.input[rawStartIndex:startIndex],
				Start:   tl.getPos(rawStartIndex),
				End:     tl.getPos(startIndex),
			}
			tl.lstripBlocks(&dataToken, &matchCodeTag, startToken.Literal[len(startToken.Literal)-1])
			tl.segments = append(tl.segments,
				dataToken,
				token.Token{
					Type:    token.TOKEN_RAW_END,
					Literal: tl.input[startIndex : endIndex+endLen],
					Start:   startToken.Start,
					End:     endToken.End,
				},
			)
		} else if !inRaw {
			tl.segments = append(tl.segments, dataToken, startToken, codeToken, endToken)
		}

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

func (tl *TemplateLexer) normalizeNewlines(value string) string {
	replacer := strings.NewReplacer("\r\n", tl.config.NewlineSequence, "\r", tl.config.NewlineSequence, "\n", tl.config.NewlineSequence)
	return replacer.Replace(value)
}

// isBalance 检查 [start, end) 范围内的引号和括号是否平衡
func (tl *TemplateLexer) isBalance(start int, end int) bool {
	stack := NewByteStack()
	for i := start; i < end; i++ {
		ch := tl.input[i]
		if ch == '"' || ch == '\'' {
			quoteIndex := FindQuote(tl.input[i+1:end], ch)
			if quoteIndex == -1 {
				return false
			}
			i += quoteIndex + 1 // +1 to skip the closing quote
			continue
		}
		switch ch {
		case '(':
			stack.Push(')')
		case '[':
			stack.Push(']')
		case '{':
			stack.Push('}')
		case ')', ']', '}':
			top, ok := stack.Pop()
			if !ok || top != ch {
				return false
			}
		}
	}
	return stack.Empty()
}

// lstripBlocks 根据 stripSign 和配置决定是否需要 lstrip blocks 相关的空白字符
// 并修改 tk.Literal (block_begin 前的 segment token)
func (tl *TemplateLexer) lstripBlocks(tk *token.Token, codeTag *CodeTag, stripSign byte) {
	if tk.Literal == "" {
		return
	}
	if stripSign == '-' {
		// strip all whitespace before
		tk.Literal = strings.TrimRightFunc(tk.Literal, unicode.IsSpace)
	} else if stripSign != '+' && tl.config.LeftStripBlocks && codeTag.startType != token.TOKEN_VARIABLE_BEGIN {
		// strip whitespace from the beginning of a line to the start of a block
		// Nothing will be stripped if there are other characters before the start of the block.
		lastNewline := strings.LastIndex(tk.Literal, "\n")
		if lastNewline != -1 {
			if IsAllWhitespace(tk.Literal[lastNewline+1:]) {
				tk.Literal = tk.Literal[:lastNewline+1]
			}
		} else if IsAllWhitespace(tk.Literal) {
			tk.Literal = ""
		}
	}
}

// trimBlocks 根据 stripSign 和配置决定是否需要 trim blocks 相关的空白字符，并返回需要 trim 的字节数
func (tl *TemplateLexer) trimBlocks(stripSign byte, codeTag *CodeTag, start int) int {
	if stripSign == '-' {
		// strip all whitespace after
		return leftWhitespaceByteCount(tl.input[start:])
	} else if stripSign != '+' && tl.config.TrimBlocks && codeTag.endType != token.TOKEN_VARIABLE_END {
		// strip the first newline after a template tag
		if strings.HasPrefix(tl.input[start:], "\n") {
			return 1
		}
	}
	return 0
}
