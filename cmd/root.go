package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
	"github.com/zostay/go-std/generic"

	"github.com/zostay/zedpm/config"
	"github.com/zostay/zedpm/plugin/master"
	"github.com/zostay/zedpm/plugin/metal"
)

var (
	rootCmd = &cobra.Command{
		Use:   "zedpm",
		Short: "The Golang project management tool with the power to save the galaxy.",
	}

	logger        hclog.InterceptLogger
	screenAdapter hclog.SinkAdapter
)

func init() {
	rootCmd.AddCommand(templateFileCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(depsCmd)

	rootCmd.PersistentFlags().StringP("log-file", "o", "", "send the raw log to this file")
	rootCmd.PersistentFlags().StringP("log-level", "l", "info", "set the log level to use [trace, debug, info, warn, error]")
}

var levelMap = map[string]hclog.Level{
	"trace": hclog.Trace,
	"debug": hclog.Debug,
	"info":  hclog.Info,
	"warn":  hclog.Warn,
	"error": hclog.Error,
}

func logLevel() hclog.Level {
	l, _ := rootCmd.PersistentFlags().GetString("log-level")
	if level, hasLevel := levelMap[l]; hasLevel {
		return level
	}
	return hclog.Info
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

	lvl := logLevel()
	stdLvl := generic.Min(hclog.Info, lvl)
	logger := hclog.NewInterceptLogger(&hclog.LoggerOptions{
		Name:  "zedpm",
		Level: stdLvl,
	})

	logFile, _ := rootCmd.PersistentFlags().GetString("log-file")
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND, 0x644)
		if err != nil {
			panic(fmt.Sprintf("zedpm failed to open log file %q: %v", logFile, err))
		}

		fileLog := hclog.NewSinkAdapter(&hclog.LoggerOptions{
			Level:  lvl,
			Output: file,
		})

		logger.RegisterSink(fileLog)
	}

	screenAdapter = hclog.NewSinkAdapter(&hclog.LoggerOptions{
		Level:  lvl,
		Output: stdErr,
	})

	logger.RegisterSink(screenAdapter)

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
	ctx = hclog.WithContext(ctx, logger)
	goals, err := e.PotentialGoalsPhasesAndTasks(ctx)
	if err != nil {
		panic(fmt.Sprintf("zedpm failed to discover plugin goals: %v", err))
	}

	configureGoalsPhasesAndTasks(ctx, goals, e, runCmd, RunGoal)
	configureGoals(ctx, goals, e, depsCmd, RunDepsForGoal)

	err = rootCmd.Execute()
	cobra.CheckErr(err)
}
