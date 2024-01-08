package turingpi

import "errors"

// ErrInvalidMode Define a static error.
var ErrInvalidMode = errors.New("invalid mode")

var ErrInvalidStatus = errors.New("invalid status")

type ResultError struct {
	Reason string
}

func (m *ResultError) Error() string {
	return m.Reason
}
