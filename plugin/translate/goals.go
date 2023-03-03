package translate

import (
	"github.com/zostay/zedpm/plugin"
	"github.com/zostay/zedpm/plugin-goals/pkg/goals"
	"github.com/zostay/zedpm/plugin/api"
)

func APITaskDescriptorToPluginTaskDescription(in *api.Descriptor_Task) *goals.TaskDescription {
	return goals.NewTaskDescription(
		in.GetName(),
		in.GetShort(),
		in.GetRequires(),
	)
}

func APITaskDescriptorsToPluginTaskDescriptions(ins []*api.Descriptor_Task) []plugin.TaskDescription {
	outs := make([]plugin.TaskDescription, len(ins))
	for i, in := range ins {
		outs[i] = APITaskDescriptorToPluginTaskDescription(in)
	}
	return outs
}

func PluginTaskDescriptionToAPITaskDescriptor(in plugin.TaskDescription) *api.Descriptor_Task {
	return &api.Descriptor_Task{
		Name:     in.Name(),
		Short:    in.Short(),
		Requires: in.Requires(),
	}
}

func PluginTaskDescriptionsToAPITaskDescriptors(ins []plugin.TaskDescription) []*api.Descriptor_Task {
	outs := make([]*api.Descriptor_Task, len(ins))
	for i, in := range ins {
		outs[i] = PluginTaskDescriptionToAPITaskDescriptor(in)
	}
	return outs
}

func APIGoalDescriptorToPluginGoalDescription(in *api.Descriptor_Goal) *goals.GoalDescription {
	return goals.NewGoalDescription(
		in.GetName(),
		in.GetShort(),
		in.GetAliases()...,
	)
}

func PluginGoalDescriptionToAPIGoalDescriptor(in plugin.GoalDescription) *api.Descriptor_Goal {
	return &api.Descriptor_Goal{
		Name:    in.Name(),
		Short:   in.Short(),
		Aliases: in.Aliases(),
	}
}
