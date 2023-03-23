package group

import (
	"path"

	"github.com/zostay/zedpm/plugin"
)

// Task is the smallest unit of execution. Normally, a Task is executed
// concurrently with other Tasks and shares state with them while running as
// part of a Phase.
type Task struct {
	// Phase is the phase to which this task belongs.
	Phase *Phase

	// Name is the name of the Task to run.
	Name string

	// Task is the description of this task.
	Task plugin.TaskDescription
}

// Path returns the /goalName/phaseName/taskName path.
func (t *Task) Path() string {
	return path.Join(t.Phase.Path(), t.Name)
}

// Short returns a short description of what the task does.
func (t *Task) Short() string {
	return t.Task.Short()
}

// Tasks returns just this task.
func (t *Task) Tasks() []plugin.TaskDescription {
	return []plugin.TaskDescription{t.Task}
}
