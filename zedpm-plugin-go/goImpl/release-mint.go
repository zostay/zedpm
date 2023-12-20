package goImpl

import (
	"context"

	"github.com/zostay/zedpm/plugin"
)

// ReleaseMintTask is the implementation of the /release/mint/go task.
type ReleaseMintTask struct {
	plugin.TaskBoilerplate
}

// Check runs tests to ensure the project is ready for release.
func (s *ReleaseMintTask) Check(ctx context.Context) error {
	logger := plugin.Logger(ctx,
		"operation", "Check",
		"task", "/release/mint/go",
	)
	logger.Info("Running go test -v ./...")

	ps, err := RunTests(ctx)

	plugin.Logger(ctx,
		"exitcode", ps.ExitCode(),
	).Info("Exited go test -v ./...")

	return err
}
