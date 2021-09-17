package state

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/plugins"
	"github.com/openrelayxyz/plugeth-utils/core"
)

func PluginStateUpdate(pl *plugins.PluginLoader, blockRoot, parentRoot common.Hash, destructs map[common.Hash]struct{}, accounts map[common.Hash][]byte, storage map[common.Hash]map[common.Hash][]byte) {
	fnList := pl.Lookup("StateUpdate", func(item interface{}) bool {
		_, ok := item.(func(core.Hash, core.Hash, map[core.Hash]struct{}, map[core.Hash][]byte, map[core.Hash]map[core.Hash][]byte))
		return ok
	})
	coreDestructs := make(map[core.Hash]struct{})
	for k, v := range destructs {
		coreDestructs[core.Hash(k)] = v
	}
	coreAccounts := make(map[core.Hash][]byte)
	for k, v := range accounts {
		coreAccounts[core.Hash(k)] = v
	}
	coreStorage := make(map[core.Hash]map[core.Hash][]byte)
	for k, v := range storage {
		coreStorage[core.Hash(k)] = make(map[core.Hash][]byte)
		for h, d := range v {
			coreStorage[core.Hash(k)][core.Hash(h)] = d
		}
	}

	for _, fni := range fnList {
		if fn, ok := fni.(func(core.Hash, core.Hash, map[core.Hash]struct{}, map[core.Hash][]byte, map[core.Hash]map[core.Hash][]byte)); ok {
			fn(core.Hash(blockRoot), core.Hash(parentRoot), coreDestructs, coreAccounts, coreStorage)
		}
	}
}

func pluginStateUpdate(blockRoot, parentRoot common.Hash, destructs map[common.Hash]struct{}, accounts map[common.Hash][]byte, storage map[common.Hash]map[common.Hash][]byte) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting StateUpdate, but default PluginLoader has not been initialized")
		return
	}
	PluginStateUpdate(plugins.DefaultPluginLoader, blockRoot, parentRoot, destructs, accounts, storage)
}
