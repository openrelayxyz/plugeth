package trie

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/plugins"
	"github.com/ethereum/go-ethereum/common"
	"github.com/openrelayxyz/plugeth-utils/core"
)

func PluginPreTrieCommit(pl *plugins.PluginLoader, node common.Hash) {
	fnList := pl.Lookup("PreTrieCommit", func(item interface{}) bool {
		_, ok := item.(func(core.Hash))
		return ok
	})
	for _, fni := range fnList {
		if fn, ok := fni.(func(core.Hash)); ok {
			fn(core.Hash(node))
		}
	}
}

func pluginPreTrieCommit(node common.Hash) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting PreTrieCommit, but default PluginLoader has not been initialized")
		return
	}
	PluginPreTrieCommit(plugins.DefaultPluginLoader, node)
}

func PluginPostTrieCommit(pl *plugins.PluginLoader, node common.Hash) {
	fnList := pl.Lookup("PostTrieCommit", func(item interface{}) bool {
		_, ok := item.(func(core.Hash))
		return ok
	})
	for _, fni := range fnList {
		if fn, ok := fni.(func(core.Hash)); ok {
			fn(core.Hash(node))
		}
	}
}

func pluginPostTrieCommit(node common.Hash) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting PostTrieCommit, but default PluginLoader has not been initialized")
		return
	}
	PluginPostTrieCommit(plugins.DefaultPluginLoader, node)
}

// func PluginGetRPCCalls(pl *plugins.PluginLoader, id, method, params string) {
// 	fnList := pl.Lookup("GetRPCCalls", func(item interface{}) bool {
// 		_, ok := item.(func(string, string, string))
// 		return ok
// 	})
// 	for _, fni := range fnList {
// 		if fn, ok := fni.(func(string, string, string)); ok {
// 			fn(id, method, params)
// 		}
// 	}
// }

// func pluginGetRPCCalls(id, method, params string) {
// 	if plugins.DefaultPluginLoader == nil {
// 		log.Warn("Attempting GerRPCCalls, but default PluginLoader has not been initialized")
// 		return
// 	}
// 	PluginGetRPCCalls(plugins.DefaultPluginLoader, id, method, params)
// }