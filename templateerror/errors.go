package templateerror

import (
	"errors"
	"fmt"
	"strings"
)

// TemplateError 是所有模板错误的基接口
type TemplateError interface {
	error
	Message() string
}

// 基础错误结构体
type baseTemplateError struct {
	msg string
}

func (e *baseTemplateError) Error() string {
	if e.msg == "" {
		return "template error"
	}
	return e.msg
}

func (e *baseTemplateError) Message() string {
	return e.msg
}

// ============================================================================
// TemplateNotFound 错误
// ============================================================================

// TemplateNotFoundError 当模板不存在时抛出
type TemplateNotFoundError struct {
	baseTemplateError
	Name      string   // 模板名称
	Templates []string // 尝试查找的模板列表
}

// NewTemplateNotFoundError 创建 TemplateNotFoundError
func NewTemplateNotFoundError(name string, message string) *TemplateNotFoundError {
	if message == "" {
		message = name
	}
	return &TemplateNotFoundError{
		baseTemplateError: baseTemplateError{msg: message},
		Name:              name,
		Templates:         []string{name},
	}
}

// 实现 errors.Is 支持
func (e *TemplateNotFoundError) Is(target error) bool {
	_, ok := target.(*TemplateNotFoundError)
	return ok
}

// ============================================================================
// TemplatesNotFound 错误 (多个模板未找到)
// ============================================================================

// TemplatesNotFoundError 当多个模板都未找到时抛出
// 它是 TemplateNotFoundError 的子类型
type TemplatesNotFoundError struct {
	TemplateNotFoundError
}

// NewTemplatesNotFoundError 创建 TemplatesNotFoundError
func NewTemplatesNotFoundError(names []string, message string) *TemplatesNotFoundError {
	if message == "" && len(names) > 0 {
		var parts []string
		for _, name := range names {
			parts = append(parts, name)
		}
		message = fmt.Sprintf("none of the templates given were found: %s", strings.Join(parts, ", "))
	}

	var lastName string
	if len(names) > 0 {
		lastName = names[len(names)-1]
	}

	return &TemplatesNotFoundError{
		TemplateNotFoundError: TemplateNotFoundError{
			baseTemplateError: baseTemplateError{msg: message},
			Name:              lastName,
			Templates:         names,
		},
	}
}

// 实现 errors.Is 支持：TemplatesNotFoundError 也是 TemplateNotFoundError
func (e *TemplatesNotFoundError) Is(target error) bool {
	if _, ok := target.(*TemplateNotFoundError); ok {
		return true
	}
	_, ok := target.(*TemplatesNotFoundError)
	return ok
}

// ============================================================================
// TemplateSyntaxError 语法错误
// ============================================================================

// TemplateSyntaxError 模板语法错误
type TemplateSyntaxError struct {
	baseTemplateError
	Lineno     int    // 错误行号
	Name       string // 模板名称（如果知道）
	Filename   string // 文件名
	Source     string // 出错的源代码行
	Translated bool   // 是否已被转换
}

// NewTemplateSyntaxError 创建 TemplateSyntaxError
func NewTemplateSyntaxError(message string, lineno int, name string, filename string) *TemplateSyntaxError {
	return &TemplateSyntaxError{
		baseTemplateError: baseTemplateError{msg: message},
		Lineno:            lineno,
		Name:              name,
		Filename:          filename,
		Translated:        false,
	}
}

// Error 实现 error 接口，带有位置信息
func (e *TemplateSyntaxError) Error() string {
	if e.Translated {
		return e.msg
	}

	// 构建位置信息
	location := fmt.Sprintf("line %d", e.Lineno)
	if e.Filename != "" || e.Name != "" {
		name := e.Filename
		if name == "" {
			name = e.Name
		}
		location = fmt.Sprintf("File %q, %s", name, location)
	}

	lines := []string{e.msg, "  " + location}

	// 如果有源代码，添加源代码行
	if e.Source != "" && e.Lineno > 0 {
		sourceLines := strings.Split(e.Source, "\n")
		if e.Lineno-1 < len(sourceLines) {
			line := strings.TrimSpace(sourceLines[e.Lineno-1])
			if line != "" {
				lines = append(lines, "    "+line)
			}
		}
	}

	return strings.Join(lines, "\n")
}

// WithSource 设置源代码并返回自身（链式调用）
func (e *TemplateSyntaxError) WithSource(source string) *TemplateSyntaxError {
	e.Source = source
	return e
}

// 实现 errors.Is 支持
func (e *TemplateSyntaxError) Is(target error) bool {
	_, ok := target.(*TemplateSyntaxError)
	return ok
}

// ============================================================================
// TemplateAssertionError 断言错误（编译时错误）
// ============================================================================

// TemplateAssertionError 模板断言错误，是 TemplateSyntaxError 的子类型
type TemplateAssertionError struct {
	TemplateSyntaxError
}

// NewTemplateAssertionError 创建 TemplateAssertionError
func NewTemplateAssertionError(message string, lineno int, name string, filename string) *TemplateAssertionError {
	return &TemplateAssertionError{
		TemplateSyntaxError: TemplateSyntaxError{
			baseTemplateError: baseTemplateError{msg: message},
			Lineno:            lineno,
			Name:              name,
			Filename:          filename,
			Translated:        false,
		},
	}
}

// 实现 errors.Is 支持
func (e *TemplateAssertionError) Is(target error) bool {
	if _, ok := target.(*TemplateSyntaxError); ok {
		return true
	}
	_, ok := target.(*TemplateAssertionError)
	return ok
}

// ============================================================================
// TemplateRuntimeError 运行时错误
// ============================================================================

// TemplateRuntimeError 模板运行时错误
type TemplateRuntimeError struct {
	baseTemplateError
}

// NewTemplateRuntimeError 创建 TemplateRuntimeError
func NewTemplateRuntimeError(message string) *TemplateRuntimeError {
	return &TemplateRuntimeError{
		baseTemplateError: baseTemplateError{msg: message},
	}
}

func (e *TemplateRuntimeError) Is(target error) bool {
	_, ok := target.(*TemplateRuntimeError)
	return ok
}

// ============================================================================
// UndefinedError 未定义变量错误
// ============================================================================

// UndefinedError 当模板尝试操作未定义变量时抛出
type UndefinedError struct {
	TemplateRuntimeError
}

// NewUndefinedError 创建 UndefinedError
func NewUndefinedError(name string) *UndefinedError {
	var msg string
	if name == "" {
		msg = "undefined value"
	} else {
		msg = fmt.Sprintf("'%s' is undefined", name)
	}
	return &UndefinedError{
		TemplateRuntimeError: TemplateRuntimeError{
			baseTemplateError: baseTemplateError{msg: msg},
		},
	}
}

func (e *UndefinedError) Is(target error) bool {
	if _, ok := target.(*TemplateRuntimeError); ok {
		return true
	}
	_, ok := target.(*UndefinedError)
	return ok
}

// ============================================================================
// SecurityError 安全错误（沙箱模式）
// ============================================================================

// SecurityError 当沙箱模式下模板尝试不安全操作时抛出
type SecurityError struct {
	TemplateRuntimeError
}

// NewSecurityError 创建 SecurityError
func NewSecurityError(message string) *SecurityError {
	return &SecurityError{
		TemplateRuntimeError: TemplateRuntimeError{
			baseTemplateError: baseTemplateError{msg: message},
		},
	}
}

func (e *SecurityError) Is(target error) bool {
	if _, ok := target.(*TemplateRuntimeError); ok {
		return true
	}
	_, ok := target.(*SecurityError)
	return ok
}

// ============================================================================
// FilterArgumentError 过滤器参数错误
// ============================================================================

// FilterArgumentError 当过滤器被调用时使用了不正确的参数
type FilterArgumentError struct {
	TemplateRuntimeError
	FilterName string
	Argument   string
}

// NewFilterArgumentError 创建 FilterArgumentError
func NewFilterArgumentError(filterName, argument, reason string) *FilterArgumentError {
	message := fmt.Sprintf("filter %q called with invalid argument %q: %s", filterName, argument, reason)
	return &FilterArgumentError{
		TemplateRuntimeError: TemplateRuntimeError{
			baseTemplateError: baseTemplateError{msg: message},
		},
		FilterName: filterName,
		Argument:   argument,
	}
}

// NewFilterArgumentErrorf 创建带格式化的 FilterArgumentError
func NewFilterArgumentErrorf(filterName, argument, format string, args ...interface{}) *FilterArgumentError {
	reason := fmt.Sprintf(format, args...)
	return NewFilterArgumentError(filterName, argument, reason)
}

func (e *FilterArgumentError) Is(target error) bool {
	if _, ok := target.(*TemplateRuntimeError); ok {
		return true
	}
	_, ok := target.(*FilterArgumentError)
	return ok
}

// ============================================================================
// 哨兵错误变量（用于 errors.Is 判断）
// ============================================================================

var (
	// ErrTemplateNotFound 是哨兵错误，表示模板不存在
	ErrTemplateNotFound = errors.New("template not found")
)

// ============================================================================
// 辅助函数：错误类型判断
// ============================================================================

// IsTemplateNotFound 判断是否为模板未找到错误（包括单模板和多模板）
func IsTemplateNotFound(err error) bool {
	var notFoundErr *TemplateNotFoundError
	return errors.As(err, &notFoundErr)
}

// IsTemplatesNotFound 判断是否为多模板未找到错误
func IsTemplatesNotFound(err error) bool {
	_, ok := err.(*TemplatesNotFoundError)
	return ok
}

// IsSyntaxError 判断是否为语法错误（包括 AssertionError）
func IsSyntaxError(err error) bool {
	var syntaxErr *TemplateSyntaxError
	return errors.As(err, &syntaxErr)
}

// IsRuntimeError 判断是否为运行时错误
func IsRuntimeError(err error) bool {
	var runtimeErr *TemplateRuntimeError
	return errors.As(err, &runtimeErr)
}

// IsUndefinedError 判断是否为未定义错误
func IsUndefinedError(err error) bool {
	_, ok := err.(*UndefinedError)
	return ok
}

// IsSecurityError 判断是否为安全错误
func IsSecurityError(err error) bool {
	_, ok := err.(*SecurityError)
	return ok
}

// IsFilterArgumentError 判断是否为过滤器参数错误
func IsFilterArgumentError(err error) bool {
	_, ok := err.(*FilterArgumentError)
	return ok
}
