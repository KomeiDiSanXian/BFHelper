// Package logger 日志
package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"runtime"
	"time"
)

// Level 日志等级
type Level int8

// Fields 日志字段
type Fields map[string]any

const (
	LevelDebug Level = iota // LevelDebug Debug
	LevelInfo               // LevelInfo Info
	LevelWarn               // LevelWarn Warn
	LevelError              // LevelError Error
	LevelFatal              // LevelFatal Fatal
	LevelPanic              // LevelPanic Panic
)

// String returns a string representation of level
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "Debug"
	case LevelInfo:
		return "Info"
	case LevelWarn:
		return "Warning"
	case LevelError:
		return "Error"
	case LevelFatal:
		return "Fatal"
	case LevelPanic:
		return "Panic"
	}
	return ""
}

// Logger 日志结构体
type Logger struct {
	newLogger *log.Logger
	ctx       context.Context
	fields    Fields
	callers   []string
}

// NewLogger 新建日志
func NewLogger(w io.Writer, prefix string, flag int) *Logger {
	return &Logger{newLogger: log.New(w, prefix, flag)}
}

func (l *Logger) clone() *Logger {
	ll := *l
	return &ll
}

// WithContext 设置日志上下文字段
func (l *Logger) WithContext(ctx context.Context) *Logger {
	ll := l.clone()
	ll.ctx = ctx
	return ll
}

// WithFields 设置日志公共字段
func (l *Logger) WithFields(f Fields) *Logger {
	ll := l.clone()
	if ll.fields == nil {
		ll.fields = make(Fields)
	}
	for k, v := range f {
		ll.fields[k] = v
	}
	return ll
}

// WithCaller 设置当前某一层的调用栈信息
func (l *Logger) WithCaller(skip int) *Logger {
	ll := l.clone()
	pc, file, line, ok := runtime.Caller(skip)
	if ok {
		f := runtime.FuncForPC(pc)
		// 打出文件名, 行号和函数名
		ll.callers = []string{fmt.Sprintf("%s:%d %s", file, line, f.Name())}
	}
	return ll
}

// WithCallerFrames 设置当前的整个调用栈信息
func (l *Logger) WithCallerFrames() *Logger {
	// 限制调用栈深度
	maxCallerDepth := 25
	minCallerDepth := 1

	callers := make([]string, 0, maxCallerDepth)
	pcs := make([]uintptr, maxCallerDepth)
	depth := runtime.Callers(minCallerDepth, pcs)

	frames := runtime.CallersFrames(pcs[:depth])
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		currentCaller := fmt.Sprintf("%s:%d %s", frame.File, frame.Line, frame.Function)
		callers = append(callers, currentCaller)
		if !more {
			break
		}
	}

	ll := l.clone()
	ll.callers = callers
	return ll
}

// JSONFormat 将日志格式化为JSON
func (l *Logger) JSONFormat(level Level, message string) map[string]any {
	data := make(Fields, len(l.fields)+4)
	data["level"] = level.String()
	data["time"] = time.Now().Local().Unix()
	data["message"] = message
	data["callers"] = l.callers

	if len(l.fields) > 0 {
		for k, v := range l.fields {
			if _, ok := data[k]; !ok {
				data[k] = v
			}
		}
	}

	return data
}

// Output 输出日志
func (l *Logger) Output(level Level, message string) {
	body, _ := json.Marshal(l.WithCaller(3).JSONFormat(level, message))
	content := string(body)
	switch level {
	case LevelDebug, LevelInfo, LevelWarn, LevelError:
		l.newLogger.Print(content)
	case LevelFatal:
		l.newLogger.Fatal(content)
	case LevelPanic:
		l.newLogger.Panic(content)
	}
}

// Debug 输出Debug级别日志
func (l *Logger) Debug(v ...interface{}) {
	l.Output(LevelDebug, fmt.Sprint(v...))
}

// Debugf 格式化输出Debug级别日志
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Output(LevelDebug, fmt.Sprintf(format, v...))
}

// Info 输出Info级别日志
func (l *Logger) Info(v ...interface{}) {
	l.Output(LevelInfo, fmt.Sprint(v...))
}

// Infof 格式化输出Info级别日志
func (l *Logger) Infof(format string, v ...interface{}) {
	l.Output(LevelInfo, fmt.Sprintf(format, v...))
}

// Warn 输出Warn级别日志
func (l *Logger) Warn(v ...interface{}) {
	l.Output(LevelWarn, fmt.Sprint(v...))
}

// Warnf 格式化输出Warn级别日志
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Output(LevelWarn, fmt.Sprintf(format, v...))
}

// Error 输出Error级别日志
func (l *Logger) Error(v ...interface{}) {
	l.Output(LevelError, fmt.Sprint(v...))
}

// Errorf 格式化输出Error级别日志
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Output(LevelError, fmt.Sprintf(format, v...))
}

// Fatal 输出Fatal级别日志
func (l *Logger) Fatal(v ...interface{}) {
	l.Output(LevelFatal, fmt.Sprint(v...))
}

// Fatalf 格式化输出Fatal级别日志
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.Output(LevelFatal, fmt.Sprintf(format, v...))
}

// Panic 输出Panic级别日志，并触发 panic
func (l *Logger) Panic(v ...interface{}) {
	l.Output(LevelPanic, fmt.Sprint(v...))
}

// Panicf 格式化输出Panic级别日志，并触发 panic
func (l *Logger) Panicf(format string, v ...interface{}) {
	l.Output(LevelPanic, fmt.Sprintf(format, v...))
}
