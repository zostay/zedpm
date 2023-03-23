package group

import (
	"fmt"
	"path"
	"strings"

	"github.com/zostay/zedpm/plugin"
)

// Phase makes up part of the overall execution plan of a Goal. It represents a
// single unit of execution that shares state between the Tasks that are run
// within concurrently and interleaved according to the Staging plan.
type Phase struct {
	// Goal is the goal to which this phase belongs.
	Goal *Goal

	// Name is the name fo the phase.
	Name string

	// InterleavedTasks is the list of tasks to be interleaved when running this
	// phase.
	InterleavedTasks []*Task
}

// Path returns the /goalName/phaseName path.
func (p *Phase) Path() string {
	return path.Join(p.Goal.Path(), p.Name)
}

// Short returns the short description of the phase, as built from the tasks.
func (p *Phase) Short() string {
	short := &strings.Builder{}
	for i, task := range p.InterleavedTasks {
		if i > 0 {
			_, _ = fmt.Fprint(short, " ")
		}
		_, _ = fmt.Fprint(short, task.Short())
	}

	return short.String()
}

// Tasks returns the plugin.TaskDescription for all the tasks that will be
// executed during this Phase.
func (p *Phase) Tasks() []plugin.TaskDescription {
	tasks := make([]plugin.TaskDescription, 0, len(p.InterleavedTasks))
	for _, t := range p.InterleavedTasks {
		tasks = append(tasks, t.Task)
	}
	return tasks
}
