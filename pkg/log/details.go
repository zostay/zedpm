package log

type DetailedError interface {
	Details() []any
}

func GetDetails(err error) []any {
	if dErr, hasDetails := err.(DetailedError); hasDetails {
		return dErr.Details()
	}
	return []any{}
}

type detailedError struct {
	err     error
	details []any
}

func (err *detailedError) Unwrap() error {
	return err.err
}

func (err *detailedError) Error() string {
	return err.err.Error()
}

func (err *detailedError) Details() []any {
	return err.details
}

func WithDetails(err error, args ...any) error {
	return &detailedError{err, args}
}
