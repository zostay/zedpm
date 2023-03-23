package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/zostay/zedpm/pkg/group"
	"github.com/zostay/zedpm/plugin/master"
)

var depsCmd = &cobra.Command{
	Use:   "deps",
	Short: "Report the order that phases will be run.",
}

// RunDepsForGoal returns a command runner for cobra that will report the
// order of execution for the goal or subtask.
func RunDepsForGoal(
	e *master.InterfaceExecutor,
	phases []*group.Phase,
) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		goal := phases[0].Goal
		for i, phase := range goal.ExecutionPhases() {
			fmt.Printf("Phase #%d - %s:\n", i, phase.Name)
			for _, task := range phase.InterleavedTasks {
				fmt.Printf(" - %s\n", task.Name)
			}
		}

		return nil
	}
}
