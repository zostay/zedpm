package goals

import (
	"github.com/zostay/zedpm/plugin"
)

// Verify that TaskDescription implements plugin.TaskDescription.
var _ plugin.TaskDescription = &TaskDescription{}

// TaskDescription is a generic implementation of plugin.TaskDescription.
type TaskDescription struct {
	name     string
	short    string
	requires []string
}

// NewTaskDescription constructs a new TaskDescription.
func NewTaskDescription(name, short string, requires []string) *TaskDescription {
	return &TaskDescription{name, short, requires}
}

// Name is the task path of the task described.
func (t *TaskDescription) Name() string {
	return t.name
}

// Short is a short description of the task.
func (t *TaskDescription) Short() string {
	return t.short
}

// Requires is a list of zero or more task paths that this task requires to run
// prior to the current task.
func (t *TaskDescription) Requires() []string {
	return t.requires
}
