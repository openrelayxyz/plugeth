package ethconfig

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/plugins"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/consensus"

	"github.com/ethereum/go-ethereum/plugins/wrappers/backendwrapper"
	wengine "github.com/ethereum/go-ethereum/plugins/wrappers/engine"
	
	"github.com/openrelayxyz/plugeth-utils/restricted"
	pparams "github.com/openrelayxyz/plugeth-utils/restricted/params"
	pconsensus "github.com/openrelayxyz/plugeth-utils/restricted/consensus"
)



func PluginGetEngine(pl *plugins.PluginLoader, chainConfig *params.ChainConfig, db ethdb.Database) consensus.Engine {
	fnList := pl.Lookup("CreateEngine", func(item interface{}) bool {
		_, ok := item.(func(*pparams.ChainConfig, restricted.Database) pconsensus.Engine)
		return ok
	})
	for _, fni := range fnList {
		if fn, ok := fni.(func(*pparams.ChainConfig, restricted.Database) pconsensus.Engine); ok { 
			clonedConfig := backendwrapper.CloneChainConfig(chainConfig)
			wrappedDb := backendwrapper.NewDb(db)
			if engine := fn(clonedConfig, wrappedDb); engine != nil {
				wrappedEngine := wengine.NewWrappedEngine(engine)
				return wrappedEngine
			}
			
		}
	}
	return nil
}

func pluginGetEngine(chainConfig *params.ChainConfig, db ethdb.Database) consensus.Engine {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting GetEngine, but default PluginLoader has not been initialized")
		return nil
	}
	return PluginGetEngine(plugins.DefaultPluginLoader, chainConfig, db)
}