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

	// TODO If the plugin panics, that panic might not be received by zedpm for some reason, which makes debugging hard. Investigate and resolve.
	//
	// In the meantime, the solution I've found is to use the
	// runPluginServerLocally setting in load.go to run the plugin in the same
	// process as zedpm and debug cases when I am getting "connection refused"
	// errors from plugins.

	goPlugin.Serve(&goPlugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins: goPlugin.PluginSet{
			"task-interface": NewPlugin(logger, impl),
		},
		GRPCServer: goPlugin.DefaultGRPCServer,
	})
}
