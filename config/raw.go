package config

import (
	"fmt"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"

	"github.com/zostay/zedpm/storage"
)

// RawConfig is the configuration specification for HCL. See Config for details
// on what the fields represent.
type RawConfig struct {
	Properties cty.Value `hcl:"properties,optional"`

	Goals   []RawGoalConfig   `hcl:"goal,block"`
	Plugins []RawPluginConfig `hcl:"plugin,block"`
}

// RawGoalConfig is the configuration specification for HCL for goal
// configuration. See GoalConfig for details on what the fields represent.
type RawGoalConfig struct {
	Name string `hcl:"name,label"`

	EnabledPlugins  []string `hcl:"enabled,optional"`
	DisabledPlugins []string `hcl:"disabled,optional"`

	Properties cty.Value `hcl:"properties,optional"`

	Tasks   []RawTaskConfig   `hcl:"task,block"`
	Targets []RawTargetConfig `hcl:"target,block"`
}

// RawPluginConfig is the configuration specification for HCL for plugin
// configuration. See PluginConfig for details on what the fields represent.
type RawPluginConfig struct {
	Name    string `hcl:"name,label"`
	Command string `hcl:"command,label"`

	Properties cty.Value `hcl:"properties,optional"`
}

// RawTaskConfig is the configuration specification for HCL for task
// configuration. See TaskConfig for details on what the fields represent.
type RawTaskConfig struct {
	Name    string `hcl:"name,label"`
	SubTask string `hcl:"subtask,label"`

	EnabledPlugins  []string `hcl:"enabled,optional"`
	DisabledPlugins []string `hcl:"disabled,optional"`

	Properties cty.Value `hcl:"properties,optional"`

	Targets []RawTargetConfig `hcl:"target,block"`
	Tasks   []RawTaskConfig   `hcl:"task,block"`
}

// RawTargetConfig is the configuration specification for HCL for target
// configuration. See TargetConfig for details on what the fields represent.
type RawTargetConfig struct {
	Name string `hcl:"name,label"`

	EnabledPlugins  []string `hcl:"enabled,optional"`
	DisabledPlugins []string `hcl:"disabled,optional"`

	Properties cty.Value `hcl:"properties,optional"`
}

// p is a helper used by decodeRawProperties to create prefixes.
func p(prefix, key string) string {
	return prefix + key + "."
}

// decodeRawProperties takes the generic cty.Value that HCL can load generic
// values into and decodes it into a storage.KV. The resulting storage.KV is
// read-only, which is handy for detecting certain internal bugs.
func decodeRawProperties(prefix string, in cty.Value) (storage.KV, error) {
	if in.IsNull() {
		return storage.New().RO(), nil
	}

	if !in.Type().IsObjectType() {
		return nil, fmt.Errorf("%s properties must be set to an key/value map", p(prefix, "properties"))
	}

	inm := in.AsValueMap()
	out := storage.New()

	for k, v := range inm {
		if v.Type().IsMapType() {
			dv, err := decodeRawProperties(p(prefix, k), v)
			if err != nil {
				return nil, fmt.Errorf("%s %w", p(prefix, k), err)
			}

			out.Set(k, dv)
			continue
		}

		if v.Type().IsCapsuleType() {
			ev := v.EncapsulatedValue()
			out.Set(k, ev)
		}

		if v.Type() == cty.Bool {
			var val bool
			_ = gocty.FromCtyValue(v, &val)
			out.Set(k, val)
		} else if v.Type() == cty.Number {
			var val float64
			_ = gocty.FromCtyValue(v, &val)
			out.Set(k, val)
		} else if v.Type() == cty.String {
			var val string
			_ = gocty.FromCtyValue(v, &val)
			out.Set(k, val)
		} else {
			return nil, fmt.Errorf("unknown value type for key %q", k)
		}
	}

	return out.RO(), nil
}

// decodeRawConfig turns a RawConfig into a Config.
func decodeRawConfig(rc *RawConfig) (*Config, error) {
	props, err := decodeRawProperties("", rc.Properties)
	if err != nil {
		return nil, err
	}

	goals, err := decodeRawList[RawGoalConfig, GoalConfig]("", rc.Goals, decodeRawGoal)
	if err != nil {
		return nil, err
	}

	plugins, err := decodeRawList[RawPluginConfig, PluginConfig]("", rc.Plugins, decodeRawPlugin)
	if err != nil {
		return nil, err
	}

	return &Config{
		Properties: props,
		Goals:      goals,
		Plugins:    plugins,
	}, nil
}

// decodeRawList returns some []Raw* into a []* object via the given decoder.
func decodeRawList[In any, Out any](
	prefix string,
	rs []In,
	decoder func(string, *In) (*Out, error),
) ([]Out, error) {
	out := make([]Out, len(rs))
	for i := range rs {
		r := &rs[i]
		c, err := decoder(prefix, r)
		if err != nil {
			return nil, err
		}

		out[i] = *c
	}
	return out, nil
}

// decodeRawGoal converts a RawGoalConfig into a GoalConfig.
func decodeRawGoal(prefix string, in *RawGoalConfig) (*GoalConfig, error) {
	pn := p(prefix, in.Name)
	props, err := decodeRawProperties(pn, in.Properties)
	if err != nil {
		return nil, err
	}

	tasks, err := decodeRawList[RawTaskConfig, TaskConfig](pn, in.Tasks, decodeRawTask)
	if err != nil {
		return nil, err
	}

	targets, err := decodeRawList[RawTargetConfig, TargetConfig](pn, in.Targets, decodeRawTarget)
	if err != nil {
		return nil, err
	}

	return &GoalConfig{
		Name:            in.Name,
		EnabledPlugins:  in.EnabledPlugins,
		DisabledPlugins: in.DisabledPlugins,
		Properties:      props,
		Tasks:           tasks,
		Targets:         targets,
	}, nil
}

// decodeRawPlugin converts a RawPluginConfig into a PluginConfig.
func decodeRawPlugin(prefix string, in *RawPluginConfig) (*PluginConfig, error) {
	pn := p(prefix, in.Name)
	props, err := decodeRawProperties(pn, in.Properties)
	if err != nil {
		return nil, err
	}

	return &PluginConfig{
		Name:       in.Name,
		Command:    in.Command,
		Properties: props,
	}, nil
}

// decodeRawTask converts a RawTaskConfig into a TaskConfig.
func decodeRawTask(prefix string, in *RawTaskConfig) (*TaskConfig, error) {
	pn := p(prefix, in.Name)
	props, err := decodeRawProperties(pn, in.Properties)
	if err != nil {
		return nil, err
	}

	targets, err := decodeRawList[RawTargetConfig, TargetConfig](pn, in.Targets, decodeRawTarget)
	if err != nil {
		return nil, err
	}

	tasks, err := decodeRawList[RawTaskConfig, TaskConfig](pn, in.Tasks, decodeRawTask)
	if err != nil {
		return nil, err
	}

	return &TaskConfig{
		Name:            in.Name,
		SubTask:         in.SubTask,
		EnabledPlugins:  in.EnabledPlugins,
		DisabledPlugins: in.DisabledPlugins,
		Properties:      props,
		Targets:         targets,
		Tasks:           tasks,
	}, nil
}

// decodeRawTarget converts a RawTargetConfig into a TargetConfig.
func decodeRawTarget(prefix string, in *RawTargetConfig) (*TargetConfig, error) {
	pn := p(prefix, in.Name)
	props, err := decodeRawProperties(pn, in.Properties)
	if err != nil {
		return nil, err
	}

	return &TargetConfig{
		Name:            in.Name,
		EnabledPlugins:  in.EnabledPlugins,
		DisabledPlugins: in.DisabledPlugins,
		Properties:      props,
	}, nil
}
