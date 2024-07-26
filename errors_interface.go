package serrors

import (
	"fmt"

	"golang.org/x/xerrors"
)

type AppError interface {
	Error() string
	Unwrap() error
	Format(s fmt.State, v rune)
	FormatError(p xerrors.Printer) error
	AddData(field string, data interface{}) AppError

	BadRequest() AppError
	Unauthorized() AppError
	NotFound() AppError
	InternalServerError() AppError

	New(msg ...string) AppError
	Errorf(format string, args ...interface{}) AppError
	Wrap(err error, msg ...string) AppError
	Code() int
	Status() int
	Is(err error) bool
	IsStatus(status int) bool
}
