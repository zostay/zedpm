package gitImpl

import (
	"context"

	"github.com/zostay/zedpm/pkg/goals"
	"github.com/zostay/zedpm/plugin"
)

// Verify that Plugin implements plugin.Interface.
var _ plugin.Interface = &Plugin{}

// Plugin implements the plugin.Interface for performing tasks related to git.
type Plugin struct{}

// Implements provides task descriptions for /release/mint/git and
// /release/publish/git tasks.
func (p *Plugin) Implements(context.Context) ([]plugin.TaskDescription, error) {
	release := goals.DescribeRelease()
	return []plugin.TaskDescription{
		release.Task("mint/git", "Verify work directory is clean and push a release branch."),
		release.Task("publish/git", "Push a release tag.",
			release.TaskName("mint")),
	}, nil
}

// Goal returns plugin.ErrUnsupportedGoal.
func (p *Plugin) Goal(context.Context, string) (plugin.GoalDescription, error) {
	return nil, plugin.ErrUnsupportedGoal
}

// Prepare returns plugin.Task implementations for the implemented tasks.
func (p *Plugin) Prepare(
	ctx context.Context,
	task string,
) (plugin.Task, error) {
	switch task {
	case "/release/mint/git":
		return &ReleaseMintTask{}, nil
	case "/release/publish/git":
		return &ReleasePublishTask{}, nil
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
