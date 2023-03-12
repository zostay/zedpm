package config

import (
	"errors"
	"io"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"github.com/zostay/zedpm/pkg/storage"
)

var (
	// ErrMissingPrefix is returned by GoalPhaseAndTaskName when the initial
	// slash is missing.
	ErrMissingPrefix = errors.New("task path is missing prefix")

	// ErrIncorrectName is returned by GoalPhaseAndTaskName when any of the
	// names are not correct. The names are made up of one or more words joined
	// by a hyphen. Each word must start with an underscore or letter and then
	// may be followed by zero or more underscores, letters, or numbers.
	ErrIncorrectName = errors.New("task path contains incorrect names")

	// ErrTooManyNames is returned by GoalPhaseAndTaskName when more than three
	// names are present in the task path.
	ErrTooManyNames = errors.New("task path contains too many names")
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

// ActionConfig defines common configuration for goals, phases, and tasks.
type ActionConfig struct {
	// Name is the name of the action being configured.
	Name string

	// EnabledPlugins creates an allow list of plugins to use when executing
	// this action. If an enable list is provided, then only plugins on this
	// list will be executed.
	//
	// TODO Implement EnabledPlugins functionality.
	EnabledPlugins []string

	// DisabledPlugins creates a block list of plugins to disable when executing
	// this action. If a disabled list is provided, then the listed plugins will
	// not be executed when running this goal, even if they are listed in the
	// EnabledPlugins list.
	//
	// TODO Implement DisabledPlugins functionality.
	DisabledPlugins []string

	// Properties provides settings that override globals when executing this
	// action.
	Properties storage.KV

	// Targets provides configuration of targets as applies ot this action.
	Targets []TargetConfig
}

// GoalConfig contains the configuration assigned to goals, which is also
// inherited by phases and tasks.
type GoalConfig struct {
	ActionConfig

	// Tasks provides configuration of sub-tasks of this goal.
	Phases []PhaseConfig
}

// PhaseConfig contains the configuration assigned to phases, which is also
// inherited by tasks.
type PhaseConfig struct {
	ActionConfig

	// Tasks provides configuration of sub-tasks of this goal.
	Tasks []TaskConfig
}

// TaskConfig is the configuration for a sub-task, which is always nested within
// a goal configuration. This configuration will be employed while running this
// sub-task regardless of how executed from the command-line.
//
// TODO The sub-sub-task configuration here seems inconsistent and needs a look.
type TaskConfig struct {
	ActionConfig
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

var legalName = regexp.MustCompile(`[_\PL][_\PL\PN]*(?:-[_\PL][_\PL\PN]*)*`)

// GoalPhaseAndTaskName splits a string of the form /goal/phase/task and returns
// it returns strings "goal", "phase", and "task". Returns an error if the
// task path is badly formed.
//
// The following error condition are possible:
//
// Returns ErrMissingParts if the task path is missing the initial slash or
// there are not enough slashes present.
//
// Returns ErrIncorrectName if a goal, phase, or task name is not legal. The
// legal names must contain one or more words. Each word is defined as value
// starting with an underscore or letter (any Unicode letter) followed by 0 or
// more underscores, letters, or numbers (again, any Unicode number). Multiple
// words must be joined with a hyphen. That is, this
// "/_foo44-blahblah/florby_flah/bl00p" defines goal, phase, and task with legal
// name values, but "/44_floo/feebly fly/%&$*" is illegal with all three names
// being disallowed for different reasons.
//
// Returns ErrTooManyNames if the task path contains more than three slashes.
func GoalPhaseAndTaskName(taskPath string) (string, string, string, error) {
	taskParts := strings.Split(taskPath, "/")
	if taskParts[0] != "" || len(taskParts) < 2 {
		return "", "", "", ErrMissingPrefix
	}

	taskParts = taskParts[1:]
	for _, taskPart := range taskParts[1:] {
		if !legalName.MatchString(taskPart) {
			return "", "", "", ErrIncorrectName
		}
	}

	if len(taskParts) > 3 {
		return "", "", "", ErrTooManyNames
	}

	var goalName, phaseName, taskName string

	goalName = taskParts[0]
	if len(taskParts) > 1 {
		phaseName = taskParts[1]
	}
	if len(taskParts) > 2 {
		taskName = taskParts[2]
	}

	return goalName, phaseName, taskName, nil
}

// GetGoalFromPath returns the goal configuration for the given task path.
func (c *Config) GetGoalFromPath(taskPath string) (*GoalConfig, error) {
	goalName, _, _, err := GoalPhaseAndTaskName(taskPath)
	if err != nil {
		return nil, err
	}

	return c.GetGoal(goalName), nil
}

// syntheticGoal is used to generate an empty GoalConfig in cases where we need
// such a thing.
func syntheticGoal(name string) *GoalConfig {
	return &GoalConfig{
		ActionConfig: ActionConfig{
			Name: name,
		},
	}
}

// GetGoalPhaseAndTaskConfig returns the goal, phase, and task configuration of
// the given task path.
func (c *Config) GetGoalPhaseAndTaskConfig(
	taskPath string,
) (*GoalConfig, *PhaseConfig, *TaskConfig, error) {
	goalName, phaseName, taskName, err := GoalPhaseAndTaskName(taskPath)
	if err != nil {
		return nil, nil, nil, err
	}

	goal := c.GetGoal(goalName)
	if goal == nil {
		goal = syntheticGoal(goalName)
	}

	var phase *PhaseConfig
	if phaseName != "" {
		for _, phaseConfig := range goal.Phases {
			if phaseConfig.Name == phaseName {
				phase = &phaseConfig
				break
			}
		}
	}

	var task *TaskConfig
	if phase != nil && taskName != "" {
		for _, taskConfig := range phase.Tasks {
			if taskConfig.Name == taskName {
				task = &taskConfig
				break
			}
		}
	}

	return goal, phase, task, nil
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

type targetable interface {
	GetTarget(string) *TargetConfig
	GetProperties() storage.KV
}

// targetableToKV adds one or more storage.KV objects to the given layer slice
// and returns the updated slice.
func targetableToKV[T targetable](
	in T,
	targetName string,
	layers []storage.KV,
) []storage.KV {
	var target *TargetConfig
	if targetName != "" {
		target = in.GetTarget(targetName)
	}

	if target != nil {
		layers = append(layers, target.Properties)
	}

	layers = append(layers, in.GetProperties())

	return layers
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
) (*storage.KVLayer, error) {
	var (
		goal   *GoalConfig
		phase  *PhaseConfig
		task   *TaskConfig
		plugin *PluginConfig
		err    error
	)

	if taskPath != "" {
		goal, phase, task, err = c.GetGoalPhaseAndTaskConfig(taskPath)
		if err != nil {
			return nil, err
		}
	}

	if pluginName != "" {
		plugin = c.GetPlugin(pluginName)
	}

	// topmost layer is for runtime properties
	layers := make([]storage.KV, 0, 8)
	layers = append(layers, properties)

	if task != nil {
		layers = targetableToKV[*TaskConfig](task, targetName, layers)
	}

	if phase != nil {
		layers = targetableToKV[*PhaseConfig](phase, targetName, layers)
	}

	if goal != nil {
		layers = targetableToKV[*GoalConfig](goal, targetName, layers)
	}

	if plugin != nil {
		layers = append(layers, plugin.Properties)
	}

	layers = append(layers, c.Properties)

	return storage.Layers(layers...), nil
}

// GetTarget returns the TargetConfig for the given target name.
func (a *ActionConfig) GetTarget(targetName string) *TargetConfig {
	for i := range a.Targets {
		if a.Targets[i].Name == targetName {
			return &a.Targets[i]
		}
	}
	return nil
}

// GetProperties returns the properties.
func (a *ActionConfig) GetProperties() storage.KV {
	return a.Properties
}
