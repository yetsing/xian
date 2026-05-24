package ast

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

func reprNodeList(nodes []Node) string {
	args := make([]string, len(nodes))
	for i, node := range nodes {
		args[i] = node.String()
	}

	return "[" + strings.Join(args, ", ") + "]"
}

func reprExpressionList(exprs []Expression) string {
	args := make([]string, len(exprs))
	for i, expr := range exprs {
		args[i] = expr.String()
	}

	return "[" + strings.Join(args, ", ") + "]"
}

func reprStringList(strs []string) string {
	args := make([]string, len(strs))
	for i, s := range strs {
		args[i] = reprString(s)
	}

	return "[" + strings.Join(args, ", ") + "]"
}

// reprString returns python style string representation of s
func reprString(s string) string {
	qs := fmt.Sprintf("%q", s)
	if strings.ContainsRune(s, '\'') && !strings.ContainsRune(s, '"') {
		// 包含单引号 '，但无双引号 "，使用双引号包裹
		return qs
	}
	// 使用单引号包裹，并转义其中的单引号
	cs := qs[1 : len(qs)-1]
	cs = strings.ReplaceAll(cs, `'`, `\'`) // 转义单引号
	cs = strings.ReplaceAll(cs, `\"`, `"`) // 恢复双引号
	return "'" + cs + "'"
}

func reprFloat64(f float64) string {
	// 检查是否为整数
	if f == float64(int64(f)) {
		return fmt.Sprintf("%.1f", f) // Python 会显示 .0
	}

	// 检查是否需要用科学计数法
	abs := math.Abs(f)
	if (abs > 0 && abs < 0.0001) || abs >= 1e6 {
		// 使用科学计数法，但格式要像 Python
		return strconv.FormatFloat(f, 'e', -1, 64)
	}

	// 使用 %g，但去除尾随的 .0
	s := strconv.FormatFloat(f, 'f', -1, 64)
	return s
}

func reprAny(v any) string {
	switch vv := v.(type) {
	case string:
		return reprString(vv)
	case float64:
		return reprFloat64(vv)
	default:
		return fmt.Sprintf("%v", v)
	}
}

type iStringBuilder struct {
	strings.Builder

	indent     int
	indentBase string
	indentStr  string
}

func newStringBuilder(indent int) *iStringBuilder {
	sb := &iStringBuilder{
		Builder:    strings.Builder{},
		indent:     indent,
		indentBase: "  ",
	}
	if sb.indent < 0 {
		sb.indent = 0
	}
	sb.indentStr = strings.Repeat(sb.indentBase, sb.indent)
	return sb
}

// WriteLineIndent writes order: s -> newline -> indent
func (sb *iStringBuilder) WriteLineIndent(s string) {
	sb.WriteString(s)
	sb.WriteString("\n")
	sb.WriteString(sb.indentStr)
}

func (sb *iStringBuilder) WriteIndent() {
	sb.WriteString(sb.indentStr)
}

func (sb *iStringBuilder) WriteIndentN(n int) {
	sb.WriteString(strings.Repeat(sb.indentBase, n))
}

func (sb *iStringBuilder) WriteNodeList(name string, nodes []Node) {
	if len(nodes) == 0 {
		fmt.Fprintf(sb, "  %s=[],\n", name)
	} else {
		fmt.Fprintf(sb, "  %s=[\n", name)
		for _, node := range nodes {
			sb.WriteIndentN(sb.indent + 2)
			sb.WriteString(node.Dumps(sb.indent + 2))
			sb.WriteString(",\n")
		}
		sb.WriteString(sb.indentStr)
		sb.WriteString("  ],\n")
	}
}

func (sb *iStringBuilder) WriteExpressionList(name string, nodes []Expression) {
	if len(nodes) == 0 {
		fmt.Fprintf(sb, "  %s=[],\n", name)
	} else {
		fmt.Fprintf(sb, "  %s=[\n", name)
		for _, node := range nodes {
			sb.WriteIndentN(sb.indent + 2)
			sb.WriteString(node.Dumps(sb.indent + 2))
			sb.WriteString(",\n")
		}
		sb.WriteString(sb.indentStr)
		sb.WriteString("  ],\n")
	}
}

func (sb *iStringBuilder) WriteStringList(name string, strs []string) {
	if len(strs) == 0 {
		fmt.Fprintf(sb, "  %s=[],\n", name)
	} else {
		fmt.Fprintf(sb, "  %s=[\n", name)
		for _, s := range strs {
			sb.WriteIndentN(sb.indent + 2)
			fmt.Fprintf(sb, "%s,\n", reprString(s))
		}
		sb.WriteString(sb.indentStr)
		sb.WriteString("  ],\n")
	}
}

func SetCtx(node Node, ctx string) {
	todo := []Node{node}
	for len(todo) > 0 {
		node := todo[0]
		switch n := node.(type) {
		case *Name:
			n.Ctx = ctx
		case *Tuple:
			n.Ctx = ctx
		case *Getitem:
			n.Ctx = ctx
		case *Getattr:
			n.Ctx = ctx
		}

		todo = todo[1:]
		for _, child := range node.ChildNodes() {
			if child.Value != nil {
				todo = append(todo, child.Value)
			} else {
				todo = append(todo, child.Values...)
			}
		}

	}
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}

	vi := reflect.ValueOf(i)
	// 只有指针、通道、函数、接口、切片、映射等类型才能调用 IsNil
	switch vi.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Chan, reflect.Func, reflect.Slice, reflect.Interface:
		return vi.IsNil()
	}
	return false
}

func safeDumps(node Node, indent int) string {
	if isNil(node) {
		return "none" // Python 中的 None
	}
	return node.Dumps(indent)
}
