package cmd

import (
	"context"
	"path"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"

	"github.com/zostay/zedpm/config"
	"github.com/zostay/zedpm/plugin/master"
	"github.com/zostay/zedpm/plugin/metal"
)

// configureTasks ests up the master plugin interface and plugin executor. Then,
// it contacts each plugin to see what it implements and uses that information
// to configure the available run subcommands on the command-line.
func configureTasks(
	logger hclog.Logger,
	cfg *config.Config,
	plugins metal.Clients,
	attachCmd *cobra.Command,
) error {
	ifaces, err := metal.DispenseAll(plugins)
	if err != nil {
		return err
	}

	m := master.NewInterface(logger, cfg, ifaces)
	e := master.NewExecutor(logger, m)

	ctx := context.Background()
	groups, err := e.TaskGroups(ctx)
	if err != nil {
		return err
	}

	cmds := make(map[string]*cobra.Command, len(groups))
	for _, group := range groups {
		cmd := configureGoalCommand(group, e)
		attachCmd.AddCommand(cmd)
		cmds[group.Tree] = cmd

		for _, sub := range group.SubTasks() {
			cmd := configureTaskCommand(sub, e)
			parent := path.Dir(sub.Tree)
			cmds[parent].AddCommand(cmd)
			cmds[sub.Tree] = cmd
		}
	}

	return nil
}

// configureTaskCommand builds and returns the configuration for a single
// subcommand for a given subtask.
func configureTaskCommand(
	group *master.TaskGroup,
	e *master.InterfaceExecutor,
) *cobra.Command {
	shorts := make([]string, len(group.Tasks))
	for i, task := range group.Tasks {
		shorts[i] = task.Short()
	}

	return &cobra.Command{
		Use:   path.Base(group.Tree),
		Short: strings.Join(shorts, " "),
		RunE:  RunGoal(e, group),
	}
}

// configureGoalCommand builds and returns the configuration for a single
// subcommand for a given goal.
func configureGoalCommand(
	group *master.TaskGroup,
	e *master.InterfaceExecutor,
) *cobra.Command {
	return &cobra.Command{
		Use:     group.Goal.Name(),
		Short:   group.Goal.Short(),
		Aliases: group.Goal.Aliases(),
		RunE:    RunGoal(e, group),
	}
}
