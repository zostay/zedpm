package group

import (
	"path"

	"github.com/zostay/zedpm/plugin"
)

// Goal represents a top-level command. A goal is broken down into phases that
// run as a unit. Phases are run in a particular order based upon the run
// requirements of Tasks that are run in each phase. The Tasks run within a
// Phase are run concurrently and interleaved according to the Staging plan.
type Goal struct {
	// Name of the goal this group represents.
	Name string

	// Goal describes the goal for this group.
	Goal plugin.GoalDescription

	// Phases names all the phases for this goal and the tasks associated with
	// each phase.
	Phases map[string]*Phase

	// PhaseOrder is the order the phases should be run.
	PhaseOrder []string
}

// Path returns the /goalName.
func (g *Goal) Path() string {
	return path.Join("", g.Name)
}

// Short returns the short description of the goal.
func (g *Goal) Short() string {
	return g.Goal.Short()
}

// Aliases returns other names for this goal.
func (g *Goal) Aliases() []string {
	return g.Goal.Aliases()
}

// Tasks returns the plugin.TaskDescription of all tasks that will run at some
// point as part of this Goal.
func (g *Goal) Tasks() []plugin.TaskDescription {
	tasks := make([]plugin.TaskDescription, len(g.Phases))
	for _, name := range g.PhaseOrder {
		tasks = append(tasks, g.Phases[name].Tasks()...)
	}
	return tasks
}

// ExecutionPhases returns teh ordered list of Phases that are part of this
// Goal. The Phases will be logically executed in this order.
func (g *Goal) ExecutionPhases() []*Phase {
	phases := make([]*Phase, len(g.Phases))
	for i, name := range g.PhaseOrder {
		phases[i] = g.Phases[name]
	}
	return phases
}
