package plugins

import (
	"github.com/hashicorp/go-plugin"
)

type PluginManager interface {
	ServePlugins()
}

type Plugin interface {
	plugin.Plugin
}

type PluginsManager struct {
	PluginManager
}

func (pm *PluginsManager) ServePlugins() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins:         PluginMap,
		GRPCServer:      plugin.DefaultGRPCServer,
	})
}

func NewPluginManager() PluginManager {
	return &PluginsManager{}
}
