package metal

import (
	"os"

	"github.com/hashicorp/go-hclog"
	goPlugin "github.com/hashicorp/go-plugin"

	"github.com/zostay/zedpm/plugin"
)

// RunPlugin is the function called by a plugin's main function to handle making
// the plugin.Interface it defines available for execution by the zedpm master
// process.
func RunPlugin(impl plugin.Interface) {
	name, err := os.Executable()
	if err != nil {
		name = "plugin"
	}

	logger := hclog.New(&hclog.LoggerOptions{
		Name:       name,
		JSONFormat: true,
	})

	goPlugin.Serve(&goPlugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins: goPlugin.PluginSet{
			"task-interface": NewPlugin(logger, impl),
		},
		GRPCServer: goPlugin.DefaultGRPCServer,
	})
}
