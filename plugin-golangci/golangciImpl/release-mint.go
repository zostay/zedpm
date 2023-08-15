package golangciImpl

import (
	"context"

	"github.com/zostay/zedpm/plugin"
)

// ReleaseMintTask is the implementation of the /release/mint/golangci task.
type ReleaseMintTask struct {
	plugin.TaskBoilerplate
}

// Check runs golangci linter to ensure the project is ready for release.
func (s *ReleaseMintTask) Check(ctx context.Context) error {
	logger := plugin.Logger(ctx,
		"operation", "Check",
		"task", "/release/mint/golangci",
	)
	logger.Info("Running golangci-lint run ./...")

	ps, err := RunLinter(ctx)

	plugin.Logger(ctx,
		"exitcode", ps.ExitCode(),
	).Info("Exited golangci-lint run ./...")

	return err
}
