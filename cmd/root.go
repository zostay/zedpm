package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	config2 "github.com/zostay/zedpm/config"
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
}

func Execute() {
	cfg, err := config2.LocateAndLoad()
	if err != nil {
		panic(fmt.Sprintf("zedpm failed to load: %v\n", err))
	}

	stdOut := metal.NewSyncBuffer(os.Stdout)
	stdErr := metal.NewSyncBuffer(os.Stderr)

	plugins, err := metal.LoadPlugins(cfg, stdOut, stdErr)
	if err != nil {
		panic(err) // TODO Fix this panic, it's temporary
	}
	defer metal.KillPlugins(plugins)

	err = configureTasks(cfg, plugins, runCmd)
	if err != nil {
		panic(fmt.Sprintf("zedpm failed to configure goals: %v\n", err))
	}

	err = rootCmd.Execute()
	cobra.CheckErr(err)
}
