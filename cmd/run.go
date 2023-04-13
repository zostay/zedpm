package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/zostay/zedpm/pkg/group"
	"github.com/zostay/zedpm/plugin/master"
)

var runCmd = &cobra.Command{
	Use:   "run [ -t <target> ] *[ -d <key>=<value> ]",
	Short: "Execute the tasks to achieve the named goal.",
}

func init() {
	runCmd.PersistentFlags().StringP("target", "t", "default", "the target configuration to use")
	runCmd.PersistentFlags().StringToStringP("define", "d", nil, "define a variable in a=b format")
	// TODO Figure out a contract that we can use when defining plugins to enforce dry-run that we can be relatively sure will always work...
	// runCmd.PersistentFlags().Bool("dry-run", false, "describe what would happen if the command run without doing it")
}

// RunGoal returns a command runner for cobra that will execute a particular
// goal or subtask.
func RunGoal(
	ctx context.Context,
	e *master.InterfaceExecutor,
	phases []*group.Phase,
) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		target, _ := cmd.Flags().GetString("target")
		e.SetTargetName(target)

		values, _ := cmd.Flags().GetStringToString("define")
		e.Define(values)

		phaseNames := make([]string, len(phases))
		for i, phase := range phases {
			phaseNames[i] = phase.Name
		}

		if progress != nil {
			progress.SetPhases(phaseNames)
		}

		for i, phase := range phases {
			if progress != nil {
				progress.StartPhase(i, len(phase.Tasks()))
				for _, task := range phase.Tasks() {
					progress.RegisterTask(task.Name(), task.Short())
				}
			}

			err := e.ExecutePhase(ctx, phase)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "failed to execute phase %q: %v\n", phase.Name, err)
				os.Exit(1)
			}
		}

		return nil
	}
}
