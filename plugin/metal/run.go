package metal

import (
	goPlugin "github.com/hashicorp/go-plugin"

	"github.com/zostay/zedpm/plugin"
)

// RunPlugin is the function called by a plugin's main function to handle making
// the plugin.Interface it defines available for execution by the zedpm master
// process.
func RunPlugin(impl plugin.Interface) {
	goPlugin.Serve(&goPlugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins: goPlugin.PluginSet{
			"task-interface": NewPlugin(impl),
		},
		GRPCServer: goPlugin.DefaultGRPCServer,
	})
}
