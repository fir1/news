package service

import (
	"fmt"
	"time"
)

// RetriableError is a custom error that contains a positive duration for the next retry
type RetriableError struct {
	Err        error
	RetryAfter time.Duration
}

// Error returns error message and a Retry-After duration
func (e *RetriableError) Error() string {
	return fmt.Sprintf("%s (retry after %v)", e.Err.Error(), e.RetryAfter)
}

// ErrArgument is a custom error that contains an error message for validation
type ErrArgument struct {
	Err error
}

// Error returns error message and a Retry-After duration
func (e ErrArgument) Error() string {
	return fmt.Sprintf("invalid argument: %s", e.Err.Error())
}
