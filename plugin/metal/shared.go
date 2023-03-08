package metal

import (
	"context"

	goPlugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	"github.com/zostay/zedpm/plugin"
	"github.com/zostay/zedpm/plugin/api"
	"github.com/zostay/zedpm/plugin/grpc/client"
	"github.com/zostay/zedpm/plugin/grpc/service"
)

// Handshake is the shared Hashicorp plugin handshake object.
var Handshake = goPlugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "ZEDPM_PLUGIN_MAGIC_COOKIE",
	MagicCookieValue: "Q0aHomIRxbv3sa9jlP28A3juUduYTyUnAh4MQnr3",
}

// InterfaceGRPCPlugin defines the plugin interface that zedpm uses to fill in
// the Hashicorp plugin interface.
type InterfaceGRPCPlugin struct {
	goPlugin.Plugin
	Impl plugin.Interface
}

// NewPlugin creates a new plugin object for the given plugin interface.
func NewPlugin(impl plugin.Interface) *InterfaceGRPCPlugin {
	return &InterfaceGRPCPlugin{Impl: impl}
}

// GRPCServer returns a new GRPC plugin server for implementing a plugin.
func (p *InterfaceGRPCPlugin) GRPCServer(
	_ *goPlugin.GRPCBroker,
	s *grpc.Server,
) error {
	api.RegisterTaskExecutionServer(s, service.NewGRPCTaskExecution(p.Impl))
	return nil
}

// GRPCClient returns a new GRPC plugin client for executing a plugin.
func (p *InterfaceGRPCPlugin) GRPCClient(
	_ context.Context,
	_ *goPlugin.GRPCBroker,
	c *grpc.ClientConn,
) (any, error) {
	return client.NewGRPCTaskInterface(api.NewTaskExecutionClient(c)), nil
}
