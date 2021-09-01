package state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/plugins"
	"github.com/opoenrelayxyz/plugeth-utils/core"
)

// TODO (philip): change common.Hash to core.Hash,

func PluginStateUpdate(pl *plugins.PluginLoader, blockRoot, parentRoot core.Hash, destructs map[core.Hash]struct{}, accounts map[core.Hash][]byte, storage map[core.Hash]map[core.Hash][]byte) {
	fnList := pl.Lookup("StateUpdate", func(item interface{}) bool {
		_, ok := item.(func(core.Hash, core.Hash, map[core.Hash]struct{}, map[core.Hash][]byte, map[core.Hash]map[core.Hash][]byte))
		return ok
	})
	for _, fni := range fnList {
		if fn, ok := fni.(func(core.Hash, core.Hash, map[core.Hash]struct{}, map[core.Hash][]byte, map[core.Hash]map[core.Hash][]byte)); ok {
			fn(blockRoot, parentRoot, destructs, accounts, storage)
		}
	}
}

func pluginStateUpdate(blockRoot, parentRoot core.Hash, destructs map[core.Hash]struct{}, accounts map[core.Hash][]byte, storage map[core.Hash]map[core.Hash][]byte) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting StateUpdate, but default PluginLoader has not been initialized")
		return
	}
	PluginStateUpdate(plugins.DefaultPluginLoader, blockRoot, parentRoot, destructs, accounts, storage)
}
