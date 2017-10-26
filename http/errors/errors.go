// Package errors provides http error functionality
package errors

// Wrap wraps provided err and http response code
func Wrap(err error, code int) *Error {
	return &Error{
		code: code,
		err:  err,
	}
}

// Error represents http error
type Error struct {
	code int
	err  error
}

// Code returns http response code associated with Error
func (e *Error) Code() int { return e.code }

// Err returns wrapped error
func (e *Error) Err() error { return e.err }

// Error returns error description
func (e *Error) Error() string { return e.err.Error() }
