package goImpl

import (
	"context"
	"os/exec"

	"github.com/zostay/zedpm/pkg/log"
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

	cmd := exec.CommandContext(ctx, "go", "test", "-v", "./...")
	cmd.Stdout = logger.Output(log.LevelInfo)
	cmd.Stderr = logger.Output(log.LevelError)
	err := cmd.Run()

	plugin.Logger(ctx,
		"exitcode", cmd.ProcessState.ExitCode(),
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
