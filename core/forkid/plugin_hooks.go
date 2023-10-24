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

	// fnList := pl.Lookup("ForkIDs", func(item interface{}) bool {
	// 	_, ok := item.(func([]uint64, []uint64) ([]uint64, []uint64))
	// 	return ok
	// })
	// for _, fni := range fnList {
	// 	if fn, ok := fni.(func([]uint64, []uint64) ([]uint64, []uint64)); ok {
	// 		pluginByBlock, pluginByTime := fn(byBlock, byTime)
	// 		return pluginByBlock, pluginByTime, true
	// 	}
	// }

	// return nil, nil, false
}

func pluginForkIDs(byBlock, byTime []uint64) ([]uint64, []uint64, bool) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting PluginForkIDs, but default PluginLoader has not been initialized")
		return nil, nil, false
	}
	return PluginForkIDs(plugins.DefaultPluginLoader, byBlock, byTime)
}
