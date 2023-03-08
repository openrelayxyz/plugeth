package trie

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/plugins"
)

func PluginPreTrieCommit(pl *plugins.PluginLoader, node common.Hash) {
	fnList := pl.Lookup("PreTrieCommit", func(item interface{}) bool {
		_, ok := item.(func(common.Hash))
		return ok
	})
	for _, fni := range fnList {
		if fn, ok := fni.(func(common.Hash)); ok {
			fn(node)
		}
	}
}

func pluginPreTrieCommit(node common.Hash,) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting PreTrieCommit, but default PluginLoader has not been initialized")
		return
	}
	PluginPreTrieCommit(plugins.DefaultPluginLoader, node)
}

func PluginPostTrieCommit(pl *plugins.PluginLoader, node common.Hash) {
	fnList := pl.Lookup("PostTrieCommit", func(item interface{}) bool {
		_, ok := item.(func(common.Hash))
		return ok
	})
	for _, fni := range fnList {
		if fn, ok := fni.(func(common.Hash)); ok {
			fn(node)
		}
	}
}

func pluginPostTrieCommit(node common.Hash,) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting PostTrieCommit, but default PluginLoader has not been initialized")
		return
	}
	PluginPostTrieCommit(plugins.DefaultPluginLoader, node)
}