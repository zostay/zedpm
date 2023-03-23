package group

import (
	"fmt"

	"github.com/zostay/zedpm/config"
	"github.com/zostay/zedpm/plugin"
)

type Group interface {
	Path() string
}

type PhaseGroup interface {
	ExecutionPhases() []*Phase
}

type TaskGroup interface {
	Tasks() []plugin.TaskDescription
}

func setupGoal(
	goalMap map[string]*Goal,
	name string,
	desc plugin.GoalDescription,
	n int,
) *Goal {
	if goal, hasGoal := goalMap[name]; hasGoal {
		return goal
	}

	goal := &Goal{
		Name:       name,
		Goal:       desc,
		Phases:     make(map[string]*Phase, n),
		PhaseOrder: make([]string, 0, n),
	}
	goalMap[name] = goal
	return goal
}

func setupPhaseAndTask(
	goal *Goal,
	phaseName string,
	taskName string,
	taskDesc plugin.TaskDescription,
	n int,
) {
	phase, hasPhase := goal.Phases[phaseName]
	if !hasPhase {
		phase = &Phase{
			Goal:             goal,
			Name:             phaseName,
			InterleavedTasks: make([]*Task, 0, n),
		}
		goal.Phases[phaseName] = phase
	}

	task := &Task{
		Phase: phase,
		Name:  taskName,
		Task:  taskDesc,
	}

	phase.InterleavedTasks = append(phase.InterleavedTasks, task)
}

func SetupGroups(
	taskDescs []plugin.TaskDescription,
	goalDescs map[string]plugin.GoalDescription,
) ([]*Goal, error) {
	n := len(taskDescs)
	goalMap := make(map[string]*Goal, n)

	for _, taskDesc := range taskDescs {
		goalName, phaseName, taskName, err := config.GoalPhaseAndTaskName(taskDesc.Name())
		if err != nil {
			return nil, err
		}

		goalDesc := goalDescs[goalName]
		if goalDesc == nil {
			return nil, fmt.Errorf("unable to locate goal description %q for task %q", goalName, taskDesc.Name())
		}

		goal := setupGoal(goalMap, goalName, goalDesc, n)
		setupPhaseAndTask(goal, phaseName, taskName, taskDesc, n)
	}

	goals := make([]*Goal, 0, len(goalMap))
	for _, goal := range goalMap {
		graph := NewDepsGraph(goal)
		order, err := graph.PhaseOrder()
		if err != nil {
			return nil, fmt.Errorf("unable to order phases for goal %q: %w", goal.Name, err)
		}
		goal.PhaseOrder = order
		goals = append(goals, goal)
	}

	return goals, nil
}
