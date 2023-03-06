package plugin

import (
	"context"
)

// Task provides some operations to help perform a task. The task is executed
// in a series of stages and if multiple plugins implement a given task, they
// may be run in parallel. The task operations are executed in the following
// order:
//
// * Setup
//
// * Check
//
// * Begin (in ascending Order)
//
// * Run (in ascending Order)
//
// * End (in ascending Order)
//
// * Finish
//
// * Teardown
type Task interface {
	// Setup should be used exclusively for initial setup of the Task, such as
	// acquiring resources, setting clients and connections, and other
	// housekeeping tasks. No operations related to performing the task should
	// be performed here.
	Setup(context.Context) error

	// Check should only execute initial validation and guard functions to
	// prevent the operation from continuing. No writes or modifications ought
	// to be performed here. This can also be a reasonable stage at which to
	// detect and setup property values for later use.
	Check(context.Context) error

	// Begin should be used to prepare early stage operations for the task.
	Begin(context.Context) (Operations, error)

	// Run should be used to prepare middle stage operations for the task.
	Run(context.Context) (Operations, error)

	// End should be used to prepare late stage operations for the task.
	End(context.Context) (Operations, error)

	// Finish should be used for any final checks to ensure the task has been
	// performed correctly and sanely. Modifications should not be performed
	// at this late stage.
	Finish(context.Context) error

	// Teardown should only be used to releasing resources, closing connections,
	// and internal cleanup. No operations realted to performing the task should
	// be performed here.
	Teardown(context.Context) error
}

// Tasks is a list of Task objects.
type Tasks []Task

// TaskBoilerplate is intended at simplifying the implementation of a Task by
// providing empty, noop implementations for all Task methods. That way a Task
// only needs to implement the stages and operations it needs to implement.
type TaskBoilerplate struct{}

func (TaskBoilerplate) Setup(context.Context) error               { return nil }
func (TaskBoilerplate) Check(context.Context) error               { return nil }
func (TaskBoilerplate) Begin(context.Context) (Operations, error) { return nil, nil }
func (TaskBoilerplate) Run(context.Context) (Operations, error)   { return nil, nil }
func (TaskBoilerplate) End(context.Context) (Operations, error)   { return nil, nil }
func (TaskBoilerplate) Finish(context.Context) error              { return nil }
func (TaskBoilerplate) Teardown(context.Context) error            { return nil }
