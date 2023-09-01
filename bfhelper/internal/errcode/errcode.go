// Package errcode 错误码标准化
package errcode

import (
	"fmt"

	zero "github.com/wdvxdr1123/ZeroBot"
)

// Error 错误
type Error struct {
	code    int
	msg     string
	details map[string]any
}

var codes = map[int]string{}

// NewError 创建一个错误到codes并返回Error指针
func NewError(code int, msg string) *Error {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("code %d already exists", code))
	}
	codes[code] = msg
	return &Error{code: code, msg: msg, details: make(map[string]any)}
}

// Code 返回错误码
func (e *Error) Code() int {
	return e.code
}

// Message 返回错误信息
func (e *Error) Message() string {
	return e.msg
}

// Error 打印错误信息
func (e *Error) Error() string {
	return fmt.Sprintf("code: %d, message: %s, details: %v", e.Code(), e.Message(), e.Details())
}

// Messagef 格式化输出错误信息
func (e *Error) Messagef(args ...any) string {
	return fmt.Sprintf(e.msg, args)
}

// Details 输出错误有关细节
func (e *Error) Details() map[string]any {
	return e.details
}

// WithDetails 添加错误细节
func (e *Error) WithDetails(k string, v any) *Error {
	newError := *e
	newDetails := make(map[string]any)
	for key, val := range e.details {
		newDetails[key] = val
	}
	newDetails[k] = v
	newError.details = newDetails
	return &newError
}

// WithZeroContext 添加zero.Ctx 相关细节
func (e *Error) WithZeroContext(ctx *zero.Ctx) *Error {
	return e.
		WithDetails("User", ctx.Event.UserID).
		WithDetails("Command", ctx.Event.RawMessage).
		WithDetails("Time", ctx.Event.Time).
		WithDetails("Bot", ctx.Event.SelfID).
		WithDetails("Group", ctx.Event.GroupID)
}
