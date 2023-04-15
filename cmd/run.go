package cmd

import (
	"context"
	"path"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

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

		caser := cases.Title(language.AmericanEnglish)
		phaseNames := make([]string, len(phases))
		for i, phase := range phases {
			phaseNames[i] = caser.String(phase.Name)
		}

		if progress != nil {
			progress.SetPhases(phaseNames)
		}

		for i, phase := range phases {
			if progress != nil {
				progress.StartPhase(i, len(phase.Tasks()))
				for _, task := range phase.Tasks() {
					progress.RegisterTask(task.Name(), caser.String(path.Base(task.Name())))
				}
			}
		}

		phasePlan := e.PreparePhasePlan(phases)
		for phasePlan.NextPhase() {
			err := phasePlan.ExecutePhase(ctx)
			if err != nil {
				logger.Error("failed to execute phase", "phase", phasePlan.CurrentPhase().Name, "error", err)
				exitStatus = 1
			}
		}

		return nil
	}
}
