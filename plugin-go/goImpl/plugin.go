package goImpl

import (
	"context"

	"github.com/zostay/zedpm/pkg/goals"
	"github.com/zostay/zedpm/plugin"
)

// Verify that Plugin implements plugin.Interface
var _ plugin.Interface = &Plugin{}

// Plugin implements the plugin.Interface for performing tasks related to go.
type Plugin struct{}

// Implements provides task descriptions for /test/run/go tasks.
func (p *Plugin) Implements(ctx context.Context) ([]plugin.TaskDescription, error) {
	release := goals.DescribeRelease()
	test := goals.DescribeTest()
	return []plugin.TaskDescription{
		release.Task("mint", "go", "Run tests to ensure the project is ready for release."),
		test.Task("run", "go", "Run the go test command."),
	}, nil
}

// Goal returns plugin.ErrUnsupportedGoal.
func (p *Plugin) Goal(ctx context.Context, name string) (plugin.GoalDescription, error) {
	return nil, plugin.ErrUnsupportedGoal
}

// Prepare returns plugin.Task implementations for the implemented tasks.
func (p *Plugin) Prepare(
	_ context.Context,
	task string,
) (plugin.Task, error) {
	switch task {
	case "/release/mint/go":
		return &ReleaseMintTask{}, nil
	case "/test/run/go":
		return &TestRunTask{}, nil
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
