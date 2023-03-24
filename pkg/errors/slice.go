package errors

import "strings"

// SliceErrors is returned by many of the processes here. It represents a list of
// errors. Since concurrency is involved with running multiple tasks at once, it
// is quite possible that multiple failures may occur simultaneously. This error
// implementation collects these errors into a super-error.
type SliceErrors []error

// Error returns all the errors inside it as a string.
func (e SliceErrors) Error() string {
	msgs := make([]string, len(e))
	for i, err := range e {
		msgs[i] = err.Error()
	}
	return strings.Join(msgs, "; ")
}

// Unwrap returns all the inner errors.
func (e SliceErrors) Unwrap() []error {
	return e
}
