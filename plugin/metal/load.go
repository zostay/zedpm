package metal

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-hclog"
	goPlugin "github.com/hashicorp/go-plugin"

	"github.com/zostay/zedpm/config"
	"github.com/zostay/zedpm/plugin"
	"github.com/zostay/zedpm/plugin-changelog/changelogImpl"
)

// runPluginServerLocally is a variable that can be configured to replace or add
// a plugin configuration that will be run from the same process as the master
// process, which can allow for easier debugging.
//
// TODO This is a cheap debugging aid, but should be made nicer and configgable somehow for debugging those sticky problems.
var runPluginServerLocally = map[string]plugin.Interface{
	"changelog": &changelogImpl.Plugin{},
}

// Clients represents a list of Hashicorp plugins we are running to implement
// the plugin interface of zedpm.
type Clients map[string]*goPlugin.Client

// devModePluginPrefix is the special prefix to note that a plugin is running in
// developer mode. When running a plugin this way, the DEV_MODE property must be
// set to "true" globally. This should make it relatively easy to run local
// plugins and plugins that are under active development.
const devModePluginPrefix = "go run "

// LoadLocalPlugin will run the plugin server as a goroutine on the local
// process and connect the master process to it, essentially talking to itself.
func LoadLocalPlugin(
	iface plugin.Interface,
	logger hclog.Logger,
	stdOut *SyncBuffer,
	stdErr *SyncBuffer,
) (*goPlugin.Client, error) {
	reattach := make(chan *goPlugin.ReattachConfig)
	go func() {
		goPlugin.Serve(&goPlugin.ServeConfig{
			Test: &goPlugin.ServeTestConfig{
				ReattachConfigCh: reattach,
			},
			HandshakeConfig: Handshake,
			Plugins: goPlugin.PluginSet{
				"task-interface": NewPlugin(logger, iface),
			},
			GRPCServer: goPlugin.DefaultGRPCServer,
		})
	}()

	rc := <-reattach

	client := goPlugin.NewClient(&goPlugin.ClientConfig{
		HandshakeConfig: Handshake,
		Plugins: map[string]goPlugin.Plugin{
			"task-interface": &InterfaceGRPCPlugin{},
		},
		Reattach:         rc,
		AllowedProtocols: []goPlugin.Protocol{goPlugin.ProtocolGRPC},
		Logger:           logger,
		SyncStdout:       stdOut,
		SyncStderr:       stdErr,
		// TODO Implement SecureConfig
	})

	return client, nil
}

// LoadDevModePlugin initializes a plugin that is compiled and run in a single
// step via the "go run" command.
func LoadDevModePlugin(
	cfg *config.Config,
	pcfg *config.PluginConfig,
	logger hclog.Logger,
	stdOut *SyncBuffer,
	stdErr *SyncBuffer,
) (*goPlugin.Client, error) {
	if !cfg.Properties.GetBool("DEV_MODE") {
		return nil, fmt.Errorf("plugin configuration has plugins in development, but DEV_MODE is not set to true")
	}

	cmd := []string{"go", "run", pcfg.Command[len(devModePluginPrefix):]}

	return NewGoPluginClient(cmd, logger, stdOut, stdErr), nil
}

// NewGoPluginClient creates a new Hashicorp plugin client to connect to a
// single configured plugin.
func NewGoPluginClient(
	cmd []string,
	logger hclog.Logger,
	stdOut *SyncBuffer,
	stdErr *SyncBuffer,
) *goPlugin.Client {
	client := goPlugin.NewClient(&goPlugin.ClientConfig{
		HandshakeConfig: Handshake,
		Plugins: map[string]goPlugin.Plugin{
			"task-interface": &InterfaceGRPCPlugin{},
		},
		Cmd:              exec.Command(cmd[0], cmd[1:]...), //nolint:gosec // foot guns have been handed to user, so tainted value here is expected
		AllowedProtocols: []goPlugin.Protocol{goPlugin.ProtocolGRPC},
		Logger:           logger,
		SyncStderr:       stdErr,
		SyncStdout:       stdOut,
		// TODO Implement SecureConfig
	})
	return client
}

// LoadPlugins will load all the configured plugins by executing their plugin
// program via the Hashicorp plugin interface for each.
func LoadPlugins(
	logger hclog.Logger,
	cfg *config.Config,
	stdOut *SyncBuffer,
	stdErr *SyncBuffer,
) (Clients, error) {
	clients := make(Clients, len(cfg.Plugins))
	for i := range cfg.Plugins {
		pcfg := &cfg.Plugins[i]

		var client *goPlugin.Client
		if plugin, wantsLocal := runPluginServerLocally[pcfg.Name]; wantsLocal {
			var err error
			client, err = LoadLocalPlugin(plugin, logger, stdOut, stdErr)
			if err != nil {
				return nil, err
			}
		} else if strings.HasPrefix(pcfg.Command, devModePluginPrefix) {
			var err error
			client, err = LoadDevModePlugin(cfg, pcfg, logger, stdOut, stdErr)
			if err != nil {
				return nil, err
			}
		} else {
			cmd := []string{"sh", "-c", pcfg.Command}
			client = NewGoPluginClient(cmd, logger, stdOut, stdErr)
		}

		clients[pcfg.Name] = client
	}
	return clients, nil
}

// Dispense returns a plugin.Interface for executing parts of a single plugin
// interface.
func Dispense(clients Clients, name string) (plugin.Interface, error) {
	client, err := clients[name].Client()
	if err != nil {
		return nil, fmt.Errorf("error connecting to plugin %q: %w", name, err)
	}

	raw, err := client.Dispense("task-interface")
	if err != nil {
		return nil, fmt.Errorf("error dispensing plugin %q: %w", name, err)
	}

	iface := raw.(plugin.Interface)
	return iface, nil
}

// DispenseAll returns a mapping from plugin name (found in the configuration)
// to the matching plugin.Interface which is able to execute tasks and operations
// and so on.
func DispenseAll(clients Clients) (map[string]plugin.Interface, error) {
	ifaces := make(map[string]plugin.Interface, len(clients))
	for name := range clients {
		iface, err := Dispense(clients, name)
		if err != nil {
			return nil, err
		}

		ifaces[name] = iface
	}
	return ifaces, nil
}

// KillPlugins will kill all plugins that have been started by LoadPlugins.
func KillPlugins(clients Clients) {
	for _, v := range clients {
		v.Kill()
	}
}
