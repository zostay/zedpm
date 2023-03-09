package errors

import "strings"

// SliceError is returned by many of the processes here. It represents a list of
// errors. Since concurrency is involved with running multiple tasks at once, it
// is quite possible that multiple failures may occur simultaneously. This error
// implementation collects these errors into a super-error.
type SliceError []error

// Error returns all the errors inside it as a string.
func (e SliceError) Error() string {
	msgs := make([]string, len(e))
	for i, err := range e {
		msgs[i] = err.Error()
	}
	return strings.Join(msgs, "; ")
}

// Unwrap returns all the inner errors.
func (e SliceError) Unwrap() []error {
	return e
}
