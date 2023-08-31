// Package errcode 错误码
package errcode

var (
	Success            = NewError(0, "Success")
	InternalError      = NewError(1000, "Internal error")
	InvalidParamsError = NewError(1001, "Invalid parameters")
	NotFoundError      = NewError(1002, "Not found")
	TimeoutError       = NewError(1003, "Timeout")

	DataBaseInternalError = NewError(2000, "DataBase internal error")
	DataBaseCreateError   = NewError(2001, "DataBase creation error")
	DataBaseUpdateError   = NewError(2002, "DataBase update error")
	DataBaseReadError     = NewError(2003, "DataBase read error")
	DataBaseDeleteError   = NewError(2004, "DataBase delete error")

	NetworkError = NewError(3000, "Network error")

	Canceled = NewError(4000, "Operation canceled")
)
