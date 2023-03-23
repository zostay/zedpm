package changelogImpl

import (
	"context"

	"github.com/zostay/zedpm/pkg/goals"
	"github.com/zostay/zedpm/plugin"
)

// Plugin is the plugin.Interface implementation for handling various
// changelog-related tasks.
type Plugin struct{}

// Ensures that Plugin is an implementation of plugin.Interface.
var _ plugin.Interface = &Plugin{}

// Goal always returns plugin.ErrUnsupportedGoal.
func (p *Plugin) Goal(context.Context, string) (plugin.GoalDescription, error) {
	return nil, plugin.ErrUnsupportedGoal
}

// Implements returns the following tasks:
//
//	/info/release/description
//	/lint/changelog
//	/release/mint/changelog
//	/release/publish/changelog
func (p *Plugin) Implements(context.Context) ([]plugin.TaskDescription, error) {
	info := goals.DescribeInfo()
	lint := goals.DescribeLint()
	release := goals.DescribeRelease()
	return []plugin.TaskDescription{
		info.Task("release", "description", "Explain the changes made for a release."),
		lint.Task("project-files", "changelog", "Check changelog for correctness."),
		release.Task("mint", "changelog", "Check and prepare changelog for release."),
		release.Task("publish", "changelog", "Capture changelog data to prepare for release.", "mint"),
	}, nil
}

// Prepare returns task implementations for each of the implemented tasks.
func (p *Plugin) Prepare(
	ctx context.Context,
	task string,
) (plugin.Task, error) {
	switch task {
	case "/lint/project-files/changelog":
		return &LintChangelogTask{}, nil
	case "/info/release/description":
		return &InfoChangelogTask{}, nil
	case "/release/mint/changelog":
		return &ReleaseMintTask{}, nil
	case "/release/publish/changelog":
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
