package cmd

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"

	"github.com/zostay/zedpm/config"
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

	err = configureTasks(logger, cfg, plugins, runCmd, RunGoal)
	if err != nil {
		panic(fmt.Sprintf("zedpm failed to configure run goals: %v\n", err))
	}

	err = configureTasks(logger, cfg, plugins, depsCmd, RunDepsForGoal)
	if err != nil {
		panic(fmt.Sprintf("zedpm failed to configure deps goals: %v\n", err))
	}

	err = rootCmd.Execute()
	cobra.CheckErr(err)
}
