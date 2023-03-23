package goalsImpl

import (
	"context"
	"os"

	"github.com/zostay/zedpm/pkg/goals"
	"github.com/zostay/zedpm/pkg/storage"
	"github.com/zostay/zedpm/plugin"
)

// Verify that Plugin is an implementation of plugin.Interface.
var _ plugin.Interface = &Plugin{}

// Plugin implements the built-in goals plugin.
type Plugin struct{}

// Implements returns that this plugin implements the /info/display task.
func (p *Plugin) Implements(context.Context) ([]plugin.TaskDescription, error) {
	info := goals.DescribeInfo()
	return []plugin.TaskDescription{
		info.Task("_finally", "display", "Display information."),
	}, nil
}

// Goal provides descriptions for the following goals: build, deploy, generate,
// info, init, install, lint, release, request, and test.
func (p *Plugin) Goal(
	_ context.Context,
	name string,
) (plugin.GoalDescription, error) {
	switch name {
	case goals.NameBuild:
		return goals.DescribeBuild(), nil
	case goals.NameDeploy:
		return goals.DescribeDeploy(), nil
	case goals.NameGenerate:
		return goals.DescribeGenerate(), nil
	case goals.NameInfo:
		return goals.DescribeInfo(), nil
	case goals.NameInit:
		return goals.DescribeInit(), nil
	case goals.NameInstall:
		return goals.DescribeInstall(), nil
	case goals.NameLint:
		return goals.DescribeLint(), nil
	case goals.NameRelease:
		return goals.DescribeRelease(), nil
	case goals.NameRequest:
		return goals.DescribeRequest(), nil
	case goals.NameTest:
		return goals.DescribeTest(), nil
	default:
		return nil, plugin.ErrUnsupportedGoal
	}
}

// Prepare returns the implementations for the implemented tasks.
func (p *Plugin) Prepare(
	_ context.Context,
	taskName string,
) (plugin.Task, error) {
	switch taskName {
	case "/info/_finally/display":
		return &InfoDisplayTask{}, nil
	}
	return nil, plugin.ErrUnsupportedTask
}

// Cancel is a no-op.
func (p *Plugin) Cancel(context.Context, plugin.Task) error {
	return nil
}

// Complete will output the accumulated properties if the /info/display task has
// been executed.
func (p *Plugin) Complete(ctx context.Context, task plugin.Task) error {
	var values storage.KV = plugin.KV(ctx)
	outputAll := goals.GetPropertyInfoOutputAll(ctx)
	if !outputAll {
		values = storage.ExportsOnly(values)
	}
	formatter := goals.InfoOutputFormatter(ctx)
	return formatter(os.Stdout, values)
}
