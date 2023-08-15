package golangciImpl

import (
	"context"
	"os/exec"

	"github.com/zostay/zedpm/pkg/log"
	"github.com/zostay/zedpm/plugin"
)

type LintGolangciTask struct {
	plugin.TaskBoilerplate
}

func (l *LintGolangciTask) RunLinter(ctx context.Context) error {
	logger := plugin.Logger(ctx,
		"operation", "RunLinter",
		"task", "/lint/project-files/golangci",
	)
	logger.Info("Running golangci-lint run ./...")

	cmd := exec.CommandContext(ctx, "golangci-lint", "run", "./...")
	cmd.Stdout = logger.Output(log.LevelInfo)
	cmd.Stderr = logger.Output(log.LevelError)
	err := cmd.Run()

	plugin.Logger(ctx,
		"exitcode", cmd.ProcessState.ExitCode(),
	).Info("Exited golangci-lint run ./...")

	return err
}

func (l *LintGolangciTask) Run(ctx context.Context) (plugin.Operations, error) {
	return plugin.Operations{
		{
			Order:  50,
			Action: plugin.OperationFunc(l.RunLinter),
		},
	}, nil
}
