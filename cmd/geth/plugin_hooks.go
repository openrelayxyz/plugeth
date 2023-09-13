package main

import (
	"github.com/ethereum/go-ethereum/cmd/utils"
	gcore "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/ethereum/go-ethereum/trie/triedb/hashdb"
	"github.com/ethereum/go-ethereum/trie/triedb/pathdb"
	
	"github.com/ethereum/go-ethereum/plugins"
	"github.com/ethereum/go-ethereum/plugins/wrappers"
	
	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted"

	"github.com/urfave/cli/v2"
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

func plugethCaptureTrieConfig(ctx *cli.Context, stack *node.Node) *trie.Config {

	ec := new(ethconfig.Config)
	utils.SetEthConfig(ctx, stack, ec)

	cc := &gcore.CacheConfig{
		TrieCleanLimit:      ec.TrieCleanCache,
		TrieCleanNoPrefetch: ec.NoPrefetch,
		TrieDirtyLimit:      ec.TrieDirtyCache,
		TrieDirtyDisabled:   ec.NoPruning,
		TrieTimeLimit:       ec.TrieTimeout,
		SnapshotLimit:       ec.SnapshotCache,
		Preimages:           ec.Preimages,
		StateHistory:        ec.StateHistory,
		StateScheme:         ec.StateScheme,
	}

	config := &trie.Config{Preimages: cc.Preimages}
	if cc.StateScheme == rawdb.HashScheme {
		config.HashDB = &hashdb.Config{
			CleanCacheSize: cc.TrieCleanLimit * 1024 * 1024,
		}
	}
	if cc.StateScheme == rawdb.PathScheme {
		config.PathDB = &pathdb.Config{
			StateHistory:   cc.StateHistory,
			CleanCacheSize: cc.TrieCleanLimit * 1024 * 1024,
			DirtyCacheSize: cc.TrieDirtyLimit * 1024 * 1024,
		}
	}

	return config
}