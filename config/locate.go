package config

import (
	"os"
	"path/filepath"

	"github.com/zostay/zedpm/pkg/storage"
)

// LocateAndLoadHome will load the user-global configuration file from
//
//	~/.zedpm.conf
//
// This file is only used for goals outside of an existing project (such as
// init).
func LocateAndLoadHome() (*Config, error) {
	userDir, err := os.UserHomeDir()
	if err != nil {
		return nil, nil
	}

	homeConf := filepath.Join(userDir, ".zedpm.conf")
	r, err := os.Open(homeConf)
	if err != nil {
		return nil, nil
	}

	return Load(homeConf, r)
}

// LocateAndLoadProject will load the local project configuration file. This
// file is loaded from the current working directory if possible. If not, this
// function will try to find the file in one of the next three folders outside
// the current directory and will stop if it appears to encounter a project
// root, which is detected by looking for a .git directory or go.mod file.
func LocateAndLoadProject() (*Config, error) {
	// TODO LocateAndLoadProject might be too smart for its own good or not smart enough.
	curDir, err := os.Getwd()
	if err != nil {
		return nil, nil
	}

	for i := 0; i < 3; i++ {
		curConf := filepath.Join(curDir, "zedpm.conf")
		r, err := os.Open(curConf)
		if err == nil {
			return Load(curConf, r)
		}

		// if we encounter a go.mod, assume this is the project dir
		goMod := filepath.Join(curDir, "go.mod")
		goModStat, err := os.Stat(goMod)
		if err == nil && !goModStat.IsDir() {
			return nil, nil
		}

		// if we encounter a .git, assume this is the project dir
		gitDir := filepath.Join(curDir, ".git")
		gitDirStat, err := os.Stat(gitDir)
		if err == nil && !gitDirStat.IsDir() {
			return nil, nil
		}

		curDir = filepath.Dir(curDir)
	}

	return nil, nil
}

// DefaultConfig is the ultimate fallback configuration, used when no other
// configuration can be found.
func DefaultConfig() *Config {
	return &Config{
		Properties: storage.New().RO(),
		Plugins: []PluginConfig{
			{
				Name:       "goals",
				Command:    "zedpm-plugin-goals",
				Properties: storage.New().RO(),
			},
			{
				Name:       "changelog",
				Command:    "zedpm-plugin-changelog",
				Properties: storage.New().RO(),
			},
			{
				Name:       "git",
				Command:    "zedpm-plugin-git",
				Properties: storage.New().RO(),
			},
			{
				Name:       "github",
				Command:    "zedpm-plugin-github",
				Properties: storage.New().RO(),
			},
			{
				Name:       "go",
				Command:    "zedpm-plugin-go",
				Properties: storage.New().RO(),
			},
		},
	}
}

// LocateAndLoad will attempt to load the configuration from the project
// directory. If that fails, it will look for a global home directory
// configuration. If that fails, it will fall back onto the ultimate default
// configuration.
func LocateAndLoad() (*Config, error) {
	// TODO When fallbacks occur here, it might be worth logging a warning or something.
	cfg, err := LocateAndLoadProject()
	if err != nil || cfg != nil {
		return cfg, err
	}

	cfg, err = LocateAndLoadHome()
	if err != nil || cfg != nil {
		return cfg, err
	}

	return DefaultConfig(), nil
}
