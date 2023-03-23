package cmd

import (
	"fmt"
	"os"

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
		graph := group.NewDepsGraph(goal)
		descs, err := graph.PhaseOrder()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to get the dependency order: %v", err)
			os.Exit(1)
		}

		for i, grp := range descs {
			fmt.Printf("Phase #%d:\n", i)
			for _, desc := range grp {
				fmt.Printf(" - %s\n", desc)
			}
		}

		return nil
	}
}
