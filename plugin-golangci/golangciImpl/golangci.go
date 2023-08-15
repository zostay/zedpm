package golangciImpl

import (
	"context"
	"os"
	"os/exec"

	"github.com/zostay/zedpm/pkg/log"
	"github.com/zostay/zedpm/plugin"
)

func RunLinter(ctx context.Context) (*os.ProcessState, error) {
	logger := plugin.Logger(ctx)

	cmd := exec.CommandContext(ctx, "golangci-lint", "run", "./...")
	cmd.Stdout = logger.Output(log.LevelInfo)
	cmd.Stderr = logger.Output(log.LevelError)
	err := cmd.Run()

	return cmd.ProcessState, err
}
