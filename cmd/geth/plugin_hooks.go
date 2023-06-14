package main

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/plugins"
	"github.com/ethereum/go-ethereum/plugins/wrappers"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted"
)

func apiTranslate(apis []core.API) []rpc.API {
	result := make([]rpc.API, len(apis))
	for i, api := range apis {
		result[i] = rpc.API{
			Namespace: api.Namespace,
			Version:   api.Version,
			Service:   api.Service,
			Public:    api.Public,
		}
	}
	return result
}

func GetAPIsFromLoader(pl *plugins.PluginLoader, stack *node.Node, backend restricted.Backend) []rpc.API {
	result := []core.API{}
	fnList := pl.Lookup("GetAPIs", func(item interface{}) bool {
		switch item.(type) {
		case func(core.Node, restricted.Backend) []core.API:
			return true
		case func(core.Node, core.Backend) []core.API:
			return true
		default:
			return false
		}
	})
	for _, fni := range fnList {
		switch fn := fni.(type) {
		case func(core.Node, restricted.Backend) []core.API:
			result = append(result, fn(wrappers.NewNode(stack), backend)...)
		case func(core.Node, core.Backend) []core.API:
			result = append(result, fn(wrappers.NewNode(stack), backend)...)
		default:
		}
	}
	return apiTranslate(result)
}

func pluginGetAPIs(stack *node.Node, backend restricted.Backend) []rpc.API {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting GetAPIs, but default PluginLoader has not been initialized")
		return []rpc.API{}
	}
	return GetAPIsFromLoader(plugins.DefaultPluginLoader, stack, backend)
}

func InitializeNode(pl *plugins.PluginLoader, stack *node.Node, backend restricted.Backend) {
	fnList := pl.Lookup("InitializeNode", func(item interface{}) bool {
		switch item.(type) {
		case func(core.Node, restricted.Backend):
			return true
		case func(core.Node, core.Backend):
			return true
		default:
			return false
		}
	})
	for _, fni := range fnList {
		switch fn := fni.(type) {
		case func(core.Node, restricted.Backend):
			fn(wrappers.NewNode(stack), backend)
		case func(core.Node, core.Backend):
			fn(wrappers.NewNode(stack), backend)
		default:
		}
	}
}

func pluginsInitializeNode(stack *node.Node, backend restricted.Backend) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting InitializeNode, but default PluginLoader has not been initialized")
		return
	}
	InitializeNode(plugins.DefaultPluginLoader, stack, backend)
}

func OnShutdown(pl *plugins.PluginLoader) {
	fnList := pl.Lookup("OnShutdown", func(item interface{}) bool {
		_, ok := item.(func())
		return ok
	})
	for _, fni := range fnList {
		fni.(func())()
	}
}

func pluginsOnShutdown() {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting OnShutdown, but default PluginLoader has not been initialized")
		return
	}
	OnShutdown(plugins.DefaultPluginLoader)
}

func BlockChain(pl *plugins.PluginLoader) {
	fnList := pl.Lookup("BlockChain", func(item interface{}) bool {
			_, ok := item.(func())
			return ok
	})
	for _, fni := range fnList {
			fni.(func())()
	}
}

func pluginBlockChain() {
	if plugins.DefaultPluginLoader == nil {
			log.Warn("Attempting BlockChain, but default PluginLoader has not been initialized")
			return
	}
	BlockChain(plugins.DefaultPluginLoader)
}
