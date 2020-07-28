package serverchan

import (
	"fmt"
)

const (
	ErrnoSuccess = 0
	ErrnoDefault = 1024

	ErrmsgBadPushToken = "bad pushtoken"
	errmsgDuplicate    = "不要重复发送同样的内容"
	ErrmsgDuplicate    = "duplicate message"
)

// Serverchan's response model
type ResponseObject struct {
	// some errno: 0, 1024
	Errno int32 `json:"errno"`

	// some errmsg: success, bad pushtoken
	Errmsg string `json:"errmsg"`

	// demo dataset: done
	Dataset string `json:"dataset"`
}

// Error that describe ResponseError
type ResponseError struct {
	Errno  int32  `json:"errno"`
	Errmsg string `json:"errmsg"`
}

// new response error (inner)
func newResponseError(obj *ResponseObject) *ResponseError {
	return &ResponseError{Errno: obj.Errno, Errmsg: obj.Errmsg}
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("%d: %v", e.Errno, e.Errmsg)
}

// Check if err is responseError
func IsResponseError(err error) bool {
	_, ok := err.(*ResponseError)
	return ok
}
