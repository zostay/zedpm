package goals

import (
	"path"

	"github.com/zostay/zedpm/plugin"
)

// Verifies that GoalDescription is an implementation of plugin.GoalDescription.
var _ plugin.GoalDescription = &GoalDescription{}

// GoalDescription provides a generic, concrete implementation of
// plugin.GoalDescription.
type GoalDescription struct {
	name    string
	short   string
	aliases []string
}

// NewGoalDescription constructs a new GoalDescription.
func NewGoalDescription(name, short string, aliases ...string) *GoalDescription {
	return &GoalDescription{name, short, aliases}
}

// Task creates a TaskDescription as a sub-task of this goal.
func (g *GoalDescription) Task(name, short string, requires ...string) *TaskDescription {
	return &TaskDescription{
		name:     g.TaskName(name),
		short:    short,
		requires: requires,
	}
}

// TaskName returns the goal as a task path.
func (g *GoalDescription) TaskName(name string) string {
	return path.Join("/"+g.name, name)
}

// Name returns the name of the goal.
func (g *GoalDescription) Name() string {
	return g.name
}

// Short returns a short description of the goal.
func (g *GoalDescription) Short() string {
	return g.short
}

// Aliases returns zero or more alternate names for this goal.
func (g *GoalDescription) Aliases() []string {
	return g.aliases
}
