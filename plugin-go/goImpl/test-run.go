package goImpl

import (
	"context"

	"github.com/zostay/zedpm/plugin"
)

type TestRunTask struct {
	plugin.TaskBoilerplate
}

func (s *TestRunTask) RunTests(ctx context.Context) error {
	logger := plugin.Logger(ctx,
		"operation", "RunTests",
		"task", "/test/run/go",
	)
	logger.Info("Running go test -v ./...")

	ps, err := RunTests(ctx)

	plugin.Logger(ctx,
		"exitcode", ps.ExitCode(),
	).Info("Exited go test -v ./...")

	return err
}

func (s *TestRunTask) Run(context.Context) (plugin.Operations, error) {
	return plugin.Operations{
		{
			Order:  50,
			Action: plugin.OperationFunc(s.RunTests),
		},
	}, nil
}
