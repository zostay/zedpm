package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/zostay/zedpm/plugin/master"
)

var depsCmd = &cobra.Command{
	Use:   "deps",
	Short: "Report the order that tasks will be run.",
}

// RunDepsForGoal returns a command runner for cobra that will report the
// order of execution for the goal or subtask.
func RunDepsForGoal(
	e *master.InterfaceExecutor,
	group *master.TaskGroup,
) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		graph := master.NewDepsGraph(group.Goal.Name(), group.Tasks)
		descs, err := graph.GroupOrder()
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to get the dependency order: %v", err)
			os.Exit(1)
		}

		for i, grp := range descs {
			fmt.Printf("Group #%d:\n", i)
			for _, desc := range grp {
				fmt.Printf(" - %s\n", desc.Name())
			}
		}

		return nil
	}
}
