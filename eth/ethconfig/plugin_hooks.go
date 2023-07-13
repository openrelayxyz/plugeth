package ethconfig

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/plugins"
	// "github.com/ethereum/go-ethereum/plugins/wrappers"
	wengine "github.com/ethereum/go-ethereum/plugins/wrappers/engine"
	// "github.com/ethereum/go-ethereum/plugins/wrappers/backendwrapper"
	// "github.com/openrelayxyz/plugeth-utils/core"
	// "github.com/openrelayxyz/plugeth-utils/restricted"

	"github.com/ethereum/go-ethereum/consensus"
	// "github.com/ethereum/go-ethereum/ethdb"
	// "github.com/ethereum/go-ethereum/node"

	pconsensus "github.com/openrelayxyz/plugeth-utils/restricted/consensus"
)



func PluginGetEngine(pl *plugins.PluginLoader) consensus.Engine {
	fnList := pl.Lookup("CreateEngine", func(item interface{}) bool {
		_, ok := item.(func() pconsensus.Engine)
		return ok
	})
	for _, fni := range fnList {
		if fn, ok := fni.(func() pconsensus.Engine); ok { 
			if engine := fn(); engine != nil {
				wrappedEngine := wengine.NewWrappedEngine(engine)
				return wrappedEngine
			}
			
		}
	}
	return nil
}

func pluginGetEngine() consensus.Engine {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting GetEngine, but default PluginLoader has not been initialized")
		return nil
	}
	return PluginGetEngine(plugins.DefaultPluginLoader)
}