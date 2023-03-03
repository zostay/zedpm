package changelogImpl

import (
	"context"

	"github.com/zostay/zedpm/plugin"
	"github.com/zostay/zedpm/plugin-goals/pkg/goals"
)

var changelogPlugin = plugin.ConfigName(Plugin{})

type Plugin struct{}

var _ plugin.Interface = &Plugin{}

func (p *Plugin) Goal(context.Context, string) (plugin.GoalDescription, error) {
	return nil, plugin.ErrUnsupportedGoal
}

func (p *Plugin) Implements(context.Context) ([]plugin.TaskDescription, error) {
	info := goals.DescribeInfo()
	lint := goals.DescribeLint()
	release := goals.DescribeRelease()
	return []plugin.TaskDescription{
		info.Task("release/description", "Explain the changes made for a release."),
		lint.Task("changelog", "Check changelog for correctness."),
		release.Task("mint/changelog", "Check and prepare changelog for release."),
		release.Task("publish/changelog", "Capture changelog data to prepare for release.",
			release.TaskName("mint")),
	}, nil
}

func (p *Plugin) Prepare(
	ctx context.Context,
	task string,
) (plugin.Task, error) {
	switch task {
	case "/lint/changelog":
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

func (p *Plugin) Cancel(ctx context.Context, task plugin.Task) error {
	return nil
}

func (p *Plugin) Complete(ctx context.Context, task plugin.Task) error {
	return nil
}
