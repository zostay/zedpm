package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/zostay/zedpm/pkg/group"
	"github.com/zostay/zedpm/plugin/master"
)

type CmdBuilder func(*master.InterfaceExecutor, []*group.Phase) func(*cobra.Command, []string) error

// configureGoals goes through all goals for all available plugins and
// configures the given runner for the given command as sub-commands named for
// each "goal".
func configureGoals(
	goals []*group.Goal,
	e *master.InterfaceExecutor,
	attachCmd *cobra.Command,
	runner CmdBuilder,
) error {
	for _, goal := range goals {
		goalCmd := configureGoalCommand(goal, e, runner)
		attachCmd.AddCommand(goalCmd)
	}

	return nil
}

// configureGoalsPhasesAndTasks goes through all goals, phases, and tasks for
// all available plugins and configures the given runner for the given command
// as sub-commands named for each "goal", each "goal phase", and each "goal
// phase task".
func configureGoalsPhasesAndTasks(
	goals []*group.Goal,
	e *master.InterfaceExecutor,
	attachCmd *cobra.Command,
	runner CmdBuilder,
) error {
	for _, goal := range goals {
		goalCmd := configureGoalCommand(goal, e, runner)
		attachCmd.AddCommand(goalCmd)

		for _, phase := range goal.Phases {
			if strings.HasPrefix(phase.Name, "_") {
				continue
			}

			phaseCmd := configurePhaseCommand(phase, e, runner)
			goalCmd.AddCommand(phaseCmd)

			for _, task := range phase.InterleavedTasks {
				if strings.HasPrefix(task.Name, "_") {
					continue
				}

				taskCmd := configureTaskCommand(task, e, runner)
				phaseCmd.AddCommand(taskCmd)
			}
		}
	}

	return nil
}

// configureTaskCommand builds and returns the configuration for a single
// subcommand for a given subtask.
func configureTaskCommand(
	task *group.Task,
	e *master.InterfaceExecutor,
	runner CmdBuilder,
) *cobra.Command {
	return &cobra.Command{
		Use:   task.Name,
		Short: task.Short(),
		RunE: runner(e, []*group.Phase{
			{InterleavedTasks: []*group.Task{task}},
		}),
	}
}

// configurePhaseCommand builds andr eturns teh configruation for a single
// subcommand for a given phase.
func configurePhaseCommand(
	phase *group.Phase,
	e *master.InterfaceExecutor,
	runner CmdBuilder,
) *cobra.Command {
	return &cobra.Command{
		Use:   phase.Name,
		Short: phase.Short(),
		RunE:  runner(e, []*group.Phase{phase}),
	}
}

// configureGoalCommand builds and returns the configuration for a single
// subcommand for a given goal.
func configureGoalCommand(
	goal *group.Goal,
	e *master.InterfaceExecutor,
	runner CmdBuilder,
) *cobra.Command {
	return &cobra.Command{
		Use:     goal.Name,
		Short:   goal.Short(),
		Aliases: goal.Aliases(),
		RunE:    runner(e, goal.ExecutionPhases()),
	}
}
