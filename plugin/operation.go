package plugin

import (
	"context"
	"fmt"
)

// Ordering is a scheme for prioritizing operations. An Ordering must be in the
// range of 0-100, inclusive. Lower Ordering values precede larger ones when it
// comes to determining the order of operations.
type Ordering int

func (o Ordering) String() string {
	return fmt.Sprintf("%02d", o)
}

// OperationFunc is an implementation of OperationHandler defined as a function.
type OperationFunc func(ctx context.Context) error

// Call calls the OperationFunc.
func (op OperationFunc) Call(ctx context.Context) error {
	return op(ctx)
}

// OperationHandler defines the type interface for operations. An operation
// executes some part of a Task.
type OperationHandler interface {
	// Call performs the operation. An error will typically result in
	// termination of the entire task with the given error.
	Call(ctx context.Context) error
}

// Operations is a list of zero or more Operation objects.
type Operations []Operation

// Operation is an object that is used to give Order to the Actions that get
// performed while executing a Stage of a Task.
type Operation struct {
	Order  Ordering
	Action OperationHandler
}

// OperationLess provides a sorting function for Operations.
func OperationLess(ops Operations) func(int, int) bool {
	return func(i, j int) bool {
		return ops[i].Order < ops[j].Order
	}
}
