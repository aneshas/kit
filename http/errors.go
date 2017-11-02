package http

// NewError wraps provided errs and http response code
// thus creating new http error
func NewError(code int, errs ...error) *Error {
	return &Error{
		code: code,
		errs: errs,
	}
}

// Error represents http error
type Error struct {
	code int
	errs []error
}

// Code returns http response code associated with Error
func (e *Error) Code() int { return e.code }

// Errs returns wrapped errors
func (e *Error) Errs() []error { return e.errs }

// Error returns error description
func (e *Error) Error() string {
	var str string
	for _, err := range e.errs {
		str += err.Error() + "\n"
	}
	return str
}
