package golangciImpl

import (
	"context"

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

	ps, err := RunLinter(ctx)

	plugin.Logger(ctx,
		"exitcode", ps.ExitCode(),
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
