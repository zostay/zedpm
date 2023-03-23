package githubImpl

import (
	"context"

	"github.com/zostay/zedpm/pkg/goals"
	"github.com/zostay/zedpm/plugin"
)

// Verifies that Plugin implements plugin.Interface.
var _ plugin.Interface = &Plugin{}

// Plugin implements plugin.Interface for handling github-related tasks.
type Plugin struct{}

// Implements returns the task descriptions for the /release/mint/github and
// /release/publish/github tasks.
func (p *Plugin) Implements(context.Context) ([]plugin.TaskDescription, error) {
	rel := goals.DescribeRelease()
	return []plugin.TaskDescription{
		rel.Task("mint", "github", "Create a Github pull request."),
		rel.Task("publish", "github", "Publish a release.", "mint"),
	}, nil
}

// Goal returns plugin.ErrUnsupportedGoal.
func (p *Plugin) Goal(context.Context, string) (plugin.GoalDescription, error) {
	return nil, plugin.ErrUnsupportedGoal
}

// Prepare returns the implemented tasks.
func (p *Plugin) Prepare(
	_ context.Context,
	task string,
) (plugin.Task, error) {
	switch task {
	case "/release/mint/github":
		return &ReleaseMintTask{}, nil
	case "/release/publish/github":
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
