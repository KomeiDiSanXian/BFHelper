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

	NetworkError            = NewError(3000, "Network error")                  // NetworkError code 3000 means something has gone wrong when doing http request
	ServerNotFoundError     = NewError(3001, "Server not found error")         // ServerNotFoundError code 3001 means we cannot find server at EA gateway
	InvalidAuthError        = NewError(3002, "Invalid map id or invalid auth") // InvalidAuthError code 3002 means map id is invalid or invalid permission
	InvalidPermissionsError = NewError(3003, "Invalid permissions")            // InvalidPermissionsError code 3003 means we don't have permission
	InvalidPlayerError      = NewError(3004, "Invalid player")                 // InvalidPlayerError code 3004 means EA server does not have info of player
	ServerNotStartError     = NewError(3005, "Server does not started")        // ServerNotStartError code 3005 means server is not started yet

	Canceled = NewError(4000, "Operation canceled") // Canceled code 4000 means operation has been canceled
)
