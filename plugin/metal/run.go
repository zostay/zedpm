package metal

import (
	goPlugin "github.com/hashicorp/go-plugin"

	"github.com/zostay/zedpm/plugin"
)

func RunPlugin(impl plugin.Interface) {
	goPlugin.Serve(&goPlugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins: goPlugin.PluginSet{
			"task-interface": NewPlugin(impl),
		},
		GRPCServer: goPlugin.DefaultGRPCServer,
	})
}
