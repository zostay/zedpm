package translate

import (
	"github.com/zostay/zedpm/pkg/goals"
	"github.com/zostay/zedpm/plugin"
	"github.com/zostay/zedpm/plugin/api"
)

// APITaskDescriptorToPluginTaskDescription translates an api.Descriptor_Task
// into a goals.TaskDescription object.
func APITaskDescriptorToPluginTaskDescription(in *api.Descriptor_Task) *goals.TaskDescription {
	return goals.NewTaskDescription(
		in.GetName(),
		in.GetShort(),
		in.GetRequires(),
	)
}

// APITaskDescriptorsToPluginTaskDescriptions translates zero or more
// api.Descriptor_Task objects into the same number of plugin.TaskDescription
// objects.
func APITaskDescriptorsToPluginTaskDescriptions(ins []*api.Descriptor_Task) []plugin.TaskDescription {
	outs := make([]plugin.TaskDescription, len(ins))
	for i, in := range ins {
		outs[i] = APITaskDescriptorToPluginTaskDescription(in)
	}
	return outs
}

// PluginTaskDescriptionToAPITaskDescriptor translates a plugin.TaskDescription
// into an api.Descriptor_Task.
func PluginTaskDescriptionToAPITaskDescriptor(in plugin.TaskDescription) *api.Descriptor_Task {
	return &api.Descriptor_Task{
		Name:     in.Name(),
		Short:    in.Short(),
		Requires: in.Requires(),
	}
}

// PluginTaskDescriptionsToAPITaskDescriptors translates zero or more
// plugin.TaskDescription objects into the same number of api.Descriptor_Task
// objects.
func PluginTaskDescriptionsToAPITaskDescriptors(ins []plugin.TaskDescription) []*api.Descriptor_Task {
	outs := make([]*api.Descriptor_Task, len(ins))
	for i, in := range ins {
		outs[i] = PluginTaskDescriptionToAPITaskDescriptor(in)
	}
	return outs
}

// APIGoalDescriptorToPluginGoalDescription translates an api.Descriptor_Goal
// into a goals.GoalDescription.
func APIGoalDescriptorToPluginGoalDescription(in *api.Descriptor_Goal) *goals.GoalDescription {
	return goals.NewGoalDescription(
		in.GetName(),
		in.GetShort(),
		in.GetAliases()...,
	)
}

// PluginGoalDescriptionToAPIGoalDescriptor translates a plugin.GoalDescription
// into an api.Descriptor_Goal.
func PluginGoalDescriptionToAPIGoalDescriptor(in plugin.GoalDescription) *api.Descriptor_Goal {
	return &api.Descriptor_Goal{
		Name:    in.Name(),
		Short:   in.Short(),
		Aliases: in.Aliases(),
	}
}
