package state

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/plugins"
	"github.com/ethereum/go-ethereum/core/state/snapshot"
	"github.com/openrelayxyz/plugeth-utils/core"
)

type pluginSnapshot struct {
	root common.Hash
}

func (s *pluginSnapshot) Root() common.Hash {
	return s.root
}

func (s *pluginSnapshot) Account(hash common.Hash) (*snapshot.Account, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *pluginSnapshot) AccountRLP(hash common.Hash) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *pluginSnapshot) Storage(accountHash, storageHash common.Hash) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func PluginStateUpdate(pl *plugins.PluginLoader, blockRoot, parentRoot common.Hash, destructs map[common.Hash]struct{}, accounts map[common.Hash][]byte, storage map[common.Hash]map[common.Hash][]byte, codeUpdates map[common.Hash][]byte) {
	fnList := pl.Lookup("StateUpdate", func(item interface{}) bool {
		_, ok := item.(func(core.Hash, core.Hash, map[core.Hash]struct{}, map[core.Hash][]byte, map[core.Hash]map[core.Hash][]byte, map[core.Hash][]byte))
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
	coreCode := make(map[core.Hash][]byte)
	for k, v := range codeUpdates {
		coreCode[core.Hash(k)] = v
	}

	for _, fni := range fnList {
		if fn, ok := fni.(func(core.Hash, core.Hash, map[core.Hash]struct{}, map[core.Hash][]byte, map[core.Hash]map[core.Hash][]byte, map[core.Hash][]byte)); ok {
			fn(core.Hash(blockRoot), core.Hash(parentRoot), coreDestructs, coreAccounts, coreStorage, coreCode)
		}
	}
}

func pluginStateUpdate(blockRoot, parentRoot common.Hash, destructs map[common.Hash]struct{}, accounts map[common.Hash][]byte, storage map[common.Hash]map[common.Hash][]byte, codeUpdates map[common.Hash][]byte) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting StateUpdate, but default PluginLoader has not been initialized")
		return
	}
	PluginStateUpdate(plugins.DefaultPluginLoader, blockRoot, parentRoot, destructs, accounts, storage, codeUpdates)
}
