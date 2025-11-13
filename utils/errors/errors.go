package errors

import "fmt"

type Error struct {
	Code    string
	Status  int
	Message string
	Err     error
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) SetErr(err error) *Error {
	e.Err = err
	return e
}
