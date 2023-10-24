package forkid

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/plugins"
)

func PluginForkIDs(pl *plugins.PluginLoader, byBlock, byTime []uint64) ([]uint64, []uint64, bool) {
	f, ok := plugins.LookupOne[func([]uint64, []uint64) ([]uint64, []uint64)](pl, "ForkIDs")
	if !ok {
		return nil, nil, false
	}
	pluginByBlock, pluginByTime := f(byBlock, byTime)

	return pluginByBlock, pluginByTime, ok

}

func pluginForkIDs(byBlock, byTime []uint64) ([]uint64, []uint64, bool) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting PluginForkIDs, but default PluginLoader has not been initialized")
		return nil, nil, false
	}
	return PluginForkIDs(plugins.DefaultPluginLoader, byBlock, byTime)
}
