package plugins

import (
	"github.com/hashicorp/go-plugin"
	plugins "karod/plugins/storage/leveldb_plugin"
)

var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "QUANTOS_KAROD_GRPC_PLUGIN",
	MagicCookieValue: "QKGP",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	//"vault_grpc":   &VaulrGrpcPlugin{},
	//"vault":        &VaultPlugin{},
	"leveldb":      &plugins.LevelDB{},
	"leveldb_grpc": &plugins.LevelDBGrpc{},
}
