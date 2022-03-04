package plugins

import (
	"context"
	"github.com/libp2p/go-libp2p-core/transport"
	"karod/plugins/storage"
)

type PluginManager interface {
	ListAll() map[PluginType]string
	LoadStoragePlugin(storage.StoragePlugin) error
	LoadTransportPlugin(transport transport.Transport)
	GetStoragePluginInterface() *storage.StorePlugin
}

type Plugin interface {
	New() interface{}
	Load() Loader
	Name() string
	Type() PluginType
	SortOrder() int
	Initializer()
	LoadBefore(pluginName string) func()
	LoadAfter(pluginName string)
	Wrap(ctx context.Context, wrapped interface{}) WrapperFunction
	Verify() error
	Signature() []byte
	Config() []map[string]string
}

type WrapperFunction func(ctx context.Context, wrapped interface{}) error

type PluginType uint32

func (pt PluginType) String() string {
	return ""
}

const (
	PluginStorage PluginType = iota
	PluginTransport
	PluginBlockChain
	PluginConsensus
	PluginSecurity
)

type Loader func(pluginType, pluginName string, pluginInterface interface{}) error
