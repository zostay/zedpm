package config

import (
	"io"
	"path"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/zostay/zedpm/storage"
)

// Config is the master configuration as we use it in the application. The
// configuration format is in HCL. The actual HCL definition are with the Raw*
// structures, which are converted into these structures to handle the
// conversion of generic JSON-style objects used for Properties into storage.KV
// objects as working with JSON-style objects in HCL is ugly without some
// conversion like this.
type Config struct {
	// Properties are the global properties that are used as the value if not
	// overridden by any other configuration section.
	Properties storage.KV

	// Goals is the configuration to apply to each goal.
	Goals []GoalConfig

	// Plugins is the configuration to apply to each plugin.
	Plugins []PluginConfig
}

type GoalConfig struct {
	// Name is the name of the goal being configured.
	Name string

	// EnabledPlugins creates an allow list of plugins to use when executing
	// this goal. If an enable list is provided, then only plugins on this list
	// will be executed.
	//
	// TODO Implement EnabledPlugins functionality for goals.
	EnabledPlugins []string

	// DisabledPlugins creates a block list of plugins to disable when executing
	// this goal. If a disabled list is provided, then the listed plugins will
	// not be executed when running this goal, even if they are listed in the
	// EnabledPlugins list.
	//
	// TODO Implement DisabledPlugins functionality for goals.
	DisabledPlugins []string

	// Properties provides settings that override globals when executing this
	// goal or one of its sub-tasks.
	Properties storage.KV

	// Tasks provides configuration of sub-tasks of this goal.
	Tasks []TaskConfig

	// Targets provides configuration of targets as applies ot this goal.
	Targets []TargetConfig
}

// PluginConfig holds the configuration to use for a particular plugin.
type PluginConfig struct {
	// Name is the name to give the plugin.
	Name string

	// Command is the command to execute to run the plugin.
	Command string

	// Properties provides settings that overwrite globals when executing this
	// plugin.
	Properties storage.KV
}

// TaskConfig is the configuration for a sub-task, which is always nested within
// a goal configuration. This configuration will be employed while running this
// sub-task regardless of how executed from the command-line.
//
// TODO The sub-sub-task configuration here seems inconsistent and needs a look.
type TaskConfig struct {
	// Name is the name of the sub-task.
	Name string

	// SubTask is the name of the sub-sub-task (and may be empty).
	SubTask string

	// EnabledPlugins creates an allow list of plugins to use when executing
	// this sub-task. If an enable list is provided, then only plugins on this
	// list will be executed.
	//
	// TODO Implement EnabledPlugins functionality for tasks.
	EnabledPlugins []string

	// DisabledPlugins creates a block list of plugins to disable when executing
	// this sub-task. If a disabled list is provided, then the listed plugins
	// will not be executed when running this sub-task, even if they are listed
	// in the EnabledPlugins list.
	//
	// TODO Implement DisabledPlugins functionality for tasks.
	DisabledPlugins []string

	// Properties are the settings used to override globals and goal settings
	// when executing this sub-task.
	Properties storage.KV

	// Targets is configuration that should apply just to this sub-task.
	Targets []TargetConfig

	// Tasks is nested sub-tasks.
	Tasks []TaskConfig
}

// TargetConfig is the configuration of a target, which allows for multiple
// configurations for each goal in case you need a different configuration per
// environment or per output binary or whatever.
type TargetConfig struct {
	// Name is the name to give the target.
	Name string

	// EnabledPlugins creates an allow list of plugins to use when executing
	// this target. If an enable list is provided, then only plugins on this
	// list will be executed.
	//
	// TODO Implement EnabledPlugins functionality for targets.
	EnabledPlugins []string

	// DisabledPlugins creates a block list of plugins to disable when executing
	// this target. If a disabled list is provided, then the listed plugins
	// will not be executed when running this sub-task, even if they are listed
	// in the EnabledPlugins list.
	//
	// TODO Implement DisabledPlugins functionality for targets.
	DisabledPlugins []string

	// Properties are the settings that will override those of the global
	// settings or the parent goal or task.
	Properties storage.KV
}

// Load will load the HCL configuration from the given io.Reader, using the
// given filename as the one passed to the HCL library, which it uses to help
// generate useful error messages.
func Load(filename string, in io.Reader) (*Config, error) {
	var raw RawConfig
	fileBytes, err := io.ReadAll(in)
	if err != nil {
		return nil, err
	}

	file, diags := hclsyntax.ParseConfig(fileBytes, filename, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, diags
	}

	diags = gohcl.DecodeBody(file.Body, nil, &raw)
	if diags.HasErrors() {
		return nil, diags
	}

	return decodeRawConfig(&raw)
}

// GoalAndTaskNames splits a string of the form /goal/task/subtask and returns
// it as "goal" and []string{"task", "subtask"}.
func GoalAndTaskNames(taskPath string) (string, []string) {
	taskPath = path.Clean(taskPath)
	if taskPath == "" || taskPath == "/" {
		return "", nil
	}

	taskParts := strings.Split(taskPath, "/")
	if taskParts[0] == "" {
		taskParts = taskParts[1:]
	}

	if len(taskParts) == 1 {
		return taskParts[0], nil
	}

	return taskParts[0], taskParts[1:]
}

// GetGoalFromPath returns the goal configuration for the given task path.
func (c *Config) GetGoalFromPath(taskPath string) *GoalConfig {
	goalName, _ := GoalAndTaskNames(taskPath)
	return c.GetGoal(goalName)
}

// syntheticGoal is used to generate an empty GoalConfig in cases where we need
// such a thing.
func syntheticGoal(name string) *GoalConfig {
	return &GoalConfig{
		Name: name,
	}
}

// GetGoalAndTasks returns the goal and configuration of applicable sub-tasks
// for the given task path.
func (c *Config) GetGoalAndTasks(taskPath string) (*GoalConfig, []*TaskConfig) {
	goalName, taskNames := GoalAndTaskNames(taskPath)

	goal := c.GetGoal(goalName)
	if goal == nil {
		goal = syntheticGoal(goalName)
	}

	tasks := make([]*TaskConfig, 0, len(taskNames))
	taskList := goal.Tasks

TaskLoop:
	for _, taskName := range taskNames {
		for j := range taskList {
			if taskList[j].Name == taskName {
				tasks = append(tasks, &taskList[j])
				taskList = taskList[j].Tasks
				continue TaskLoop
			}
		}
		break
	}
	return goal, tasks
}

// GetGoal returns the GoalConfig for the given goal name.
func (c *Config) GetGoal(goalName string) *GoalConfig {
	for i := range c.Goals {
		if c.Goals[i].Name == goalName {
			return &c.Goals[i]
		}
	}
	return nil
}

// GetPlugin returns the PluginConfig for the given plugin name.
func (c *Config) GetPlugin(pluginName string) *PluginConfig {
	for i := range c.Plugins {
		if c.Plugins[i].Name == pluginName {
			return &c.Plugins[i]
		}
	}
	return nil
}

// ToKV builds and returns a storage.KVLayer containing the configuration layers
// matching the given taskPath, targetName, and pluginName in proper order
// (i.e., so that scope overrides happen correctly). The given properties store
// will be the top-most layer.
//
// The scopes override each other in the following order, with the first
// mentioned item overriding everything below it:
//
// 1. Property Settings (from the given properties argument)
//
// 2. Target Settings on Task
//
// 3. Task Settings
//
// 4. Target Settings on Goal
//
// 5. Goal Settings
//
// 6. Plugin Settings
//
// 7. Global Settings
func (c *Config) ToKV(
	properties storage.KV,
	taskPath,
	targetName,
	pluginName string,
) *storage.KVLayer {
	var (
		goal   *GoalConfig
		tasks  []*TaskConfig
		plugin *PluginConfig
	)

	if taskPath != "" {
		goal, tasks = c.GetGoalAndTasks(taskPath)
	}
	if pluginName != "" {
		plugin = c.GetPlugin(pluginName)
	}

	// topmost layer is for runtime properties
	layers := make([]storage.KV, 0, (len(tasks)+2)*2+1)
	layers = append(layers, properties)

	for _, task := range tasks {
		var target *TargetConfig
		if targetName != "" {
			target = task.GetTarget(targetName)
		}

		if target != nil {
			layers = append(layers, target.Properties)
		}

		layers = append(layers, task.Properties)

	}

	if goal != nil {
		var target *TargetConfig
		if targetName != "" {
			target = goal.GetTarget(targetName)
		}

		if target != nil {
			layers = append(layers, target.Properties)
		}

		layers = append(layers, goal.Properties)
	}

	if plugin != nil {
		layers = append(layers, plugin.Properties)
	}

	layers = append(layers, c.Properties)

	return storage.Layers(layers...)
}

// GetTarget returns the TargetConfig for the given target name.
func (g *GoalConfig) GetTarget(targetName string) *TargetConfig {
	for i := range g.Targets {
		if g.Targets[i].Name == targetName {
			return &g.Targets[i]
		}
	}
	return nil
}

// GetTarget returns the TargetConfig for the given target name.
func (t *TaskConfig) GetTarget(targetName string) *TargetConfig {
	for i := range t.Targets {
		if t.Targets[i].Name == targetName {
			return &t.Targets[i]
		}
	}
	return nil
}
