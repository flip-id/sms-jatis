package sms_jatis

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

// ErrNilArgs is returned when the argument is nil.
var ErrNilArgs = errors.New("nil arguments")

func getUnknownErr(resp *http.Response, byteBody []byte) (err error) {
	if resp.StatusCode < http.StatusBadRequest {
		return
	}

	err = &UnknownError{
		StatusCode: resp.StatusCode,
		Message:    byteBody,
	}
	return
}

// UnknownError is an error that is not defined by the documentation from the ValueFirst.
type UnknownError struct {
	StatusCode int    `json:"status_code"`
	Message    []byte `json:"message"`
}

// Error implements the error interface.
func (e UnknownError) Error() string {
	return fmt.Sprintf(
		"unknown error SMS Jatis: status code:%d message:%s",
		e.StatusCode,
		e.Message,
	)
}
