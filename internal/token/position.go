package token

import (
	"errors"
	"fmt"
	"sort"
)

type Position struct {
	Line   int
	Column int
}

func (p Position) Equal(other *Position) bool {
	return p.Line == other.Line && p.Column == other.Column
}

func (p Position) IsZero() bool {
	return p.Line == 0 && p.Column == 0
}

func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Column)
}

func (p Position) Add(lineOffset, columnOffset int) Position {
	return Position{
		Line:   p.Line + lineOffset,
		Column: p.Column + columnOffset,
	}
}

// LineColumnIndex 针对多次查询优化的索引结构（支持多种换行符）
type LineColumnIndex struct {
	lineStarts []int // 每行的起始字节索引（包括第0行作为哨兵）
	maxIndex   int   // 最大有效索引
}

// NewLineColumnIndex 预处理字符串，构建行起始位置索引
// 支持 \n, \r\n, \r 三种换行符
func NewLineColumnIndex(s string) *LineColumnIndex {
	// 预分配容量
	lineStarts := make([]int, 0, len(s)/50+1)
	lineStarts = append(lineStarts, 0) // 第1行从索引0开始

	i := 0
	for i < len(s) {
		if s[i] == '\n' {
			// Unix/Linux/macOS 换行符: \n
			lineStarts = append(lineStarts, i+1)
			i++
		} else if s[i] == '\r' {
			// 检查是否是 \r\n (Windows)
			if i+1 < len(s) && s[i+1] == '\n' {
				// \r\n 作为一个换行符
				lineStarts = append(lineStarts, i+2)
				i += 2
			} else {
				// 单独的 \r (旧版 Mac)
				lineStarts = append(lineStarts, i+1)
				i++
			}
		} else {
			i++
		}
	}

	return &LineColumnIndex{
		lineStarts: lineStarts,
		maxIndex:   len(s),
	}
}

// GetLineColumn O(log n) 时间复杂度，n 为行数
func (lci *LineColumnIndex) GetLineColumn(index int) (Position, error) {
	if index < 0 || index > lci.maxIndex {
		return Position{}, errors.New("index out of range")
	}

	// 使用 sort.Search 进行二分查找
	line := sort.Search(len(lci.lineStarts), func(i int) bool {
		return lci.lineStarts[i] > index
	})

	// 计算列号（字节偏移 + 1）
	col := index - lci.lineStarts[line-1] + 1

	return Position{Line: line, Column: col}, nil
}

// MustGetLineColumn 如果确定索引有效，可以使用这个版本（跳过错误检查）
func (lci *LineColumnIndex) MustGetLineColumn(index int) Position {
	line := sort.Search(len(lci.lineStarts), func(i int) bool {
		return lci.lineStarts[i] > index
	})
	col := index - lci.lineStarts[line-1] + 1
	return Position{Line: line, Column: col}
}

// GetLineCount 返回总行数
func (lci *LineColumnIndex) GetLineCount() int {
	return len(lci.lineStarts)
}

// GetLineStart 获取指定行的起始索引（行号从1开始）
func (lci *LineColumnIndex) GetLineStart(line int) (int, error) {
	if line < 1 || line > len(lci.lineStarts) {
		return 0, errors.New("line out of range")
	}
	return lci.lineStarts[line-1], nil
}
