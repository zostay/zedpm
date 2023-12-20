package golangciImpl

import (
	"context"

	"github.com/zostay/zedpm/pkg/goals"
	"github.com/zostay/zedpm/plugin"
)

// Verify that Plugin implements the plugin.Interface.
var _ plugin.Interface = &Plugin{}

// Plugin implements the plugin.Interface for performing tasks related to
// golangci.
type Plugin struct{}

// Goal always returns plugin.ErrUnsupportedGoal.
func (p *Plugin) Goal(ctx context.Context, name string) (plugin.GoalDescription, error) {
	return nil, plugin.ErrUnsupportedGoal
}

// Implements provides task descriptions for /lint/project-files/golangci tasks.
func (p *Plugin) Implements(context.Context) ([]plugin.TaskDescription, error) {
	release := goals.DescribeRelease()
	lint := goals.DescribeLint()
	return []plugin.TaskDescription{
		release.Task("mint", "golangci", "Run golangci-lint to ensure the project is ready for release."),
		lint.Task("project-files", "golangci", "Check project files for correctness."),
	}, nil
}

// Prepare returns plugin.Task implementations for the implemented tasks.
func (p *Plugin) Prepare(
	_ context.Context,
	task string,
) (plugin.Task, error) {
	switch task {
	case "/lint/project-files/golangci":
		return &LintGolangciTask{}, nil
	case "/release/mint/golangci":
		return &ReleaseMintTask{}, nil
	}
	return nil, plugin.ErrUnsupportedTask
}

// Cancel is a no-op.
func (p *Plugin) Cancel(ctx context.Context, task plugin.Task) error {
	return nil
}

// Complete is a no-op.
func (p *Plugin) Complete(ctx context.Context, task plugin.Task) error {
	return nil
}
