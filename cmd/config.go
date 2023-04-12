package cmd

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zostay/zedpm/pkg/group"
	"github.com/zostay/zedpm/plugin/master"
)

type CmdBuilder func(context.Context, *master.InterfaceExecutor, []*group.Phase) func(*cobra.Command, []string) error

// configureGoals goes through all goals for all available plugins and
// configures the given runner for the given command as sub-commands named for
// each "goal".
func configureGoals(
	ctx context.Context,
	goals []*group.Goal,
	e *master.InterfaceExecutor,
	attachCmd *cobra.Command,
	runner CmdBuilder,
) {
	for _, goal := range goals {
		goalCmd := configureGoalCommand(ctx, goal, e, runner)
		attachCmd.AddCommand(goalCmd)
	}
}

// configureGoalsPhasesAndTasks goes through all goals, phases, and tasks for
// all available plugins and configures the given runner for the given command
// as sub-commands named for each "goal", each "goal phase", and each "goal
// phase task".
func configureGoalsPhasesAndTasks(
	ctx context.Context,
	goals []*group.Goal,
	e *master.InterfaceExecutor,
	attachCmd *cobra.Command,
	runner CmdBuilder,
) {
	for _, goal := range goals {
		goalCmd := configureGoalCommand(ctx, goal, e, runner)
		attachCmd.AddCommand(goalCmd)

		for _, phase := range goal.Phases {
			if strings.HasPrefix(phase.Name, "_") {
				continue
			}

			phaseCmd := configurePhaseCommand(ctx, phase, e, runner)
			goalCmd.AddCommand(phaseCmd)

			for _, task := range phase.InterleavedTasks {
				if strings.HasPrefix(task.Name, "_") {
					continue
				}

				taskCmd := configureTaskCommand(ctx, task, e, runner)
				phaseCmd.AddCommand(taskCmd)
			}
		}
	}
}

// configureTaskCommand builds and returns the configuration for a single
// subcommand for a given subtask.
func configureTaskCommand(
	ctx context.Context,
	task *group.Task,
	e *master.InterfaceExecutor,
	runner CmdBuilder,
) *cobra.Command {
	return &cobra.Command{
		Use:   task.Name,
		Short: task.Short(),
		RunE: runner(ctx, e, []*group.Phase{
			{InterleavedTasks: []*group.Task{task}},
		}),
	}
}

// configurePhaseCommand builds and returns the configuration for a single
// subcommand for a given phase.
func configurePhaseCommand(
	ctx context.Context,
	phase *group.Phase,
	e *master.InterfaceExecutor,
	runner CmdBuilder,
) *cobra.Command {
	return &cobra.Command{
		Use:   phase.Name,
		Short: phase.Short(),
		RunE:  runner(ctx, e, []*group.Phase{phase}),
	}
}

// configureGoalCommand builds and returns the configuration for a single
// subcommand for a given goal.
func configureGoalCommand(
	ctx context.Context,
	goal *group.Goal,
	e *master.InterfaceExecutor,
	runner CmdBuilder,
) *cobra.Command {
	return &cobra.Command{
		Use:     goal.Name,
		Short:   goal.Short(),
		Aliases: goal.Aliases(),
		RunE:    runner(ctx, e, goal.ExecutionPhases()),
	}
}
