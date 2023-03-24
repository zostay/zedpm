package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"

	"github.com/zostay/zedpm/config"
	"github.com/zostay/zedpm/plugin/master"
	"github.com/zostay/zedpm/plugin/metal"
)

var (
	rootCmd = &cobra.Command{
		Use:   "zedpm",
		Short: "The Golang project management tool with the power to save the galaxy.",
	}
)

func init() {
	rootCmd.AddCommand(templateFileCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(depsCmd)
}

// Execute locates and loads configuration, loads the configured plugins, sets
// up the root command, and attaches the various run subcommands.
func Execute() {
	cfg, err := config.LocateAndLoad()
	if err != nil {
		panic(fmt.Sprintf("zedpm failed to load: %v\n", err))
	}

	stdOut := metal.NewSyncBuffer(os.Stdout)
	stdErr := metal.NewSyncBuffer(os.Stderr)

	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "zedpm",
		Level:  hclog.Warn,
		Output: stdErr,
	})

	plugins, err := metal.LoadPlugins(logger, cfg, stdOut, stdErr)
	if err != nil {
		panic(err) // TODO Fix this panic, it's temporary
	}
	defer metal.KillPlugins(plugins)

	ifaces, err := metal.DispenseAll(plugins)
	if err != nil {
		panic(fmt.Sprintf("zedpm failed to dispense plugins: %v", err))
	}

	m := master.NewInterface(logger, cfg, ifaces)
	e := master.NewExecutor(logger, m)

	ctx := context.Background()
	goals, err := e.PotentialGoalsPhasesAndTasks(ctx)
	if err != nil {
		panic(fmt.Sprintf("zedpm failed to discover plugin goals: %v", err))
	}

	configureGoalsPhasesAndTasks(goals, e, runCmd, RunGoal)
	configureGoals(goals, e, depsCmd, RunDepsForGoal)

	err = rootCmd.Execute()
	cobra.CheckErr(err)
}
