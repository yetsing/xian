package ast

import (
	"fmt"
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

// func dumpsNodeList(sb *strings.Builder, name string, nodes []Node, indent int) {
// 	sb.WriteString(strings.Repeat("  ", indent))
// 	if len(nodes) == 0 {
// 		fmt.Fprintf(sb, "  %s=[]\n", name)
// 	} else {
// 		fmt.Fprintf(sb, "  %s=[\n", name)
// 		for i, node := range nodes {
// 			sb.WriteString(node.Dumps(indent + 2))
// 			if i < len(nodes)-1 {
// 				sb.WriteString(",\n")
// 			}
// 		}
// 		sb.WriteString(strings.Repeat("  ", indent))
// 		sb.WriteString("  ]\n")
// 	}
// }

// func dumpsExpressionList(sb *strings.Builder, name string, nodes []Expression, indent int) {
// 	sb.WriteString(strings.Repeat("  ", indent))
// 	if len(nodes) == 0 {
// 		fmt.Fprintf(sb, "  %s=[]\n", name)
// 	} else {
// 		fmt.Fprintf(sb, "  %s=[\n", name)
// 		for i, node := range nodes {
// 			sb.WriteString(node.Dumps(indent + 2))
// 			if i < len(nodes)-1 {
// 				sb.WriteString(",\n")
// 			}
// 		}
// 		sb.WriteString(strings.Repeat("  ", indent))
// 		sb.WriteString("  ]\n")
// 	}
// }

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
