// Package errcode 错误码
package errcode

var (
	Success            = NewError(0, "Success")               // Success code 0 means success
	InternalError      = NewError(1000, "Internal error")     // InternalError code 1000 means something has gone wrong
	InvalidParamsError = NewError(1001, "Invalid parameters") // InvalidParamsError code 1001 means parameters are invalid
	NotFoundError      = NewError(1002, "Not found")          // NotFoundError code 1002 means we cannot find something at anywhere
	TimeoutError       = NewError(1003, "Timeout")            // TimeoutError code 1003 means operation timed out

	DataBaseInternalError = NewError(2000, "DataBase internal error") // DataBaseInternalError code 2000 means something has gone wrong with the database
	DataBaseCreateError   = NewError(2001, "DataBase creation error") // DataBaseCreateError code 2001 means creating something error
	DataBaseUpdateError   = NewError(2002, "DataBase update error")   // DataBaseUpdateError code 2002 means updating something error
	DataBaseReadError     = NewError(2003, "DataBase read error")     // DataBaseReadError code 2003 means reading something error
	DataBaseDeleteError   = NewError(2004, "DataBase delete error")   // DataBaseDeleteError code 2004 means deleting something error

	NetworkError = NewError(3000, "Network error") // NetworkError code 3000 means something has gone wrong when doing http request

	Canceled = NewError(4000, "Operation canceled") // Canceled code 4000 means operation has been canceled
)
