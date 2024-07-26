package errors

import (
	"fmt"
	"net/http"

	"golang.org/x/xerrors"
)

func (e *appError) new(msg string) *appError {
	e.message = msg
	e.frame = xerrors.Caller(2)
	return e
}

func (e *appError) New(msg ...string) AppError {
	var m string
	if len(msg) == 0 {
		m = fmt.Sprint(e.Code())
	} else {
		m = msg[0]
	}
	return e.new(m)
}

func (e *appError) Errorf(format string, args ...interface{}) AppError {
	return e.new(fmt.Sprintf(format, args...))
}

func (e *appError) Wrap(err error, msg ...string) AppError {
	var m string
	if len(msg) == 0 {
		m = fmt.Sprint(e.Code())
	} else {
		m = msg[0]
	}
	ne := e.new(m)
	ne.next = err
	return ne
}

func (e *appError) Code() int {
	if e.code != -1 {
		return e.code
	}
	next := AsAppError(e.next)
	if next != nil {
		return next.Code()
	}
	return -1
}

func (e *appError) Status() int {
	if e.status != 0 {
		return e.status
	}
	next := AsAppError(e.next)
	if next != nil {
		return next.Status()
	}
	return http.StatusInternalServerError
}
