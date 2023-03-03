package metal

import (
	"fmt"
	"os/exec"
	"strings"

	goPlugin "github.com/hashicorp/go-plugin"

	"github.com/zostay/zedpm/config"
	"github.com/zostay/zedpm/plugin"
)

// TODO This is a cheap debugging aid, but should be made nicer and configgable somehow for debugging those sticky problems.
var runPluginServerLocally = map[string]plugin.Interface{}

type Clients map[string]*goPlugin.Client

const devModePluginPrefix = "go run "

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
				"task-interface": NewPlugin(iface),
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

func LoadPlugins(
	cfg *config.Config,
	stdOut *SyncBuffer,
	stdErr *SyncBuffer,
) (Clients, error) {
	clients := make(Clients, len(cfg.Plugins))
	for i := range cfg.Plugins {
		pcfg := &cfg.Plugins[i]

		logger := hclog.New(&hclog.LoggerOptions{
			Output: stdErr,
			Level:  hclog.Warn,
			Name:   "plugin-" + pcfg.Name,
		})

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

func KillPlugins(clients Clients) {
	for _, v := range clients {
		v.Kill()
	}
}
