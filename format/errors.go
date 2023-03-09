package format

import (
	"fmt"

	"google.golang.org/grpc/status"
)

type WrapError struct {
	msg string
	err error
}

func (err *WrapError) Error() string {
	return err.msg
}

func (err *WrapError) Unwrap() error {
	return err.err
}

func WrapErr(err error, msg string, args ...any) error {
	fmtMsg := fmt.Sprintf(msg, args...)

	// google.golang.org/grpc formats its errors in basically the ugliest way
	// possible, so let's not let them format their own errors.
	if status, ok := status.FromError(err); ok {
		return &WrapError{
			err: err,
			msg: fmt.Sprintf("%s: %s", fmtMsg, status.Message()),
		}
	}

	return &WrapError{
		err: err,
		msg: fmt.Sprintf("%s: %v", fmtMsg, err),
	}
}
