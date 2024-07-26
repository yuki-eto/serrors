package errors

import (
	"errors"
	"fmt"

	"golang.org/x/xerrors"
)

func create(msg string) *appError {
	var e appError
	e.message = msg
	e.frame = xerrors.Caller(2)
	e.code = -1
	return &e
}

func New(msg string) AppError {
	return create(msg)
}

func Errorf(format string, args ...interface{}) AppError {
	return create(fmt.Sprintf(format, args...))
}

func BadRequestf(format string, args ...any) AppError {
	err := create("BadRequest:" + fmt.Sprintf(format, args...))
	return err.BadRequest()
}

func Forbiddenf(format string, args ...interface{}) AppError {
	err := create("forbidden: " + fmt.Sprintf(format, args...))
	return err.Forbidden()
}

func NotFoundf(format string, args ...interface{}) AppError {
	err := create("not found: " + fmt.Sprintf(format, args...))
	return err.NotFound()
}

func Wrap(err error, msg ...string) AppError {
	if err == nil {
		return nil
	}

	var m string
	if len(msg) != 0 {
		m = msg[0]
	}
	e := create(m)
	e.next = err
	return e
}

func AsAppError(err error) AppError {
	if err == nil {
		return nil
	}

	var e *appError
	if errors.As(err, &e) {
		return e
	}
	return nil
}

func IsStatus(err error, status int) bool {
	ae := AsAppError(err)
	if ae == nil {
		return false
	}
	if ae.IsStatus(status) {
		return true
	}
	ue := ae.Unwrap()
	return IsStatus(ue, status)
}

type appError struct {
	next    error
	message string
	frame   xerrors.Frame

	data   []map[string]interface{}
	code   int
	status int
}

func (e *appError) Error() string {
	// 一番下位層のメッセージを取り出す
	next := AsAppError(e.next)
	if next != nil {
		return next.Error()
	}
	if e.next == nil {
		if e.message != `` {
			return e.message
		}
		return `no message`
	}
	return e.next.Error()
}

func (e *appError) Is(err error) bool {
	if er := AsAppError(err); er != nil {
		return e.Code() == er.Code()
	}
	return false
}

func (e *appError) Unwrap() error              { return e.next }
func (e *appError) Format(s fmt.State, v rune) { xerrors.FormatError(e, s, v) }

func (e *appError) FormatError(p xerrors.Printer) error {
	var message string
	if e.code > 0 {
		message += fmt.Sprintf("[%d] ", e.code)
	}
	if e.message != "" {
		message += fmt.Sprintf("%s", e.message)
	}
	if len(e.data) != 0 {
		if message != "" {
			message += "\n"
		}
		message += fmt.Sprintf("data: %+v", e.data)
	}

	p.Print(message)
	e.frame.Format(p)
	return e.next
}

func (e *appError) AddData(field string, data interface{}) AppError {
	if e.data == nil {
		e.data = make([]map[string]interface{}, 0)
	}
	e.data = append(e.data, map[string]interface{}{field: data})
	return e
}

func RequestError(err error, code int) error {
	e := create("request error")
	e.code = code
	e.next = err
	e.status = 400
	return e
}

func ResponseError(err error, code int) error {
	e := create("response error")
	e.code = code
	e.next = err
	e.status = 500
	return e
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}
