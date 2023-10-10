package forkid

import (
	// "encoding/json"
	// "math/big"
	// "reflect"
	// "time"
	// "sync"

	// "github.com/ethereum/go-ethereum/common"
	// "github.com/ethereum/go-ethereum/core/state"
	// "github.com/ethereum/go-ethereum/core/types"
	// "github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/plugins"
	// "github.com/ethereum/go-ethereum/plugins/wrappers"
	// "github.com/ethereum/go-ethereum/rlp"
	// "github.com/openrelayxyz/plugeth-utils/core"
)

func PluginForkIDs(pl *plugins.PluginLoader) ([]uint64, []uint64, bool) {
	f, ok := plugins.LookupOne[func() ([]uint64, []uint64)](pl, "ForkIDs")
	if !ok {
		return nil, nil, false
	}
	byBlock, byTime := f()

	return byBlock, byTime, ok
}

func pluginForkIDs() ([]uint64, []uint64, bool) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting PluginForkIDs, but default PluginLoader has not been initialized")
		return nil, nil, false
	}
	return PluginForkIDs(plugins.DefaultPluginLoader)
}
