package tracers

import (
	"reflect"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/plugins"
	"github.com/ethereum/go-ethereum/plugins/interfaces"
	"github.com/ethereum/go-ethereum/plugins/wrappers"
	"github.com/openrelayxyz/plugeth-utils/core"
)

func GetPluginTracer(pl *plugins.PluginLoader, name string) (func(*state.StateDB) interfaces.TracerResult, bool) {
	tracers := pl.Lookup("Tracers", func(item interface{}) bool {
		_, ok := item.(*map[string]func(*state.StateDB) interfaces.TracerResult)
		if !ok {
			log.Warn("Found tracer that did not match type", "tracer", reflect.TypeOf(item))
		}
		return ok
	})

	for _, tmap := range tracers {
		if tracerMap, ok := tmap.(*map[string]func(core.StateDB) core.TracerResult); ok {
			if tracer, ok := (*tracerMap)[name]; ok {
				return func(sdb *state.StateDB) interfaces.TracerResult {
					return wrappers.NewWrappedTracer(tracer(wrappers.NewWrappedStateDB(sdb)))
				}, true
			}
		}
	}
	log.Info("Tracer not found", "name", name, "tracers", len(tracers))
	return nil, false
}

func getPluginTracer(name string) (func(*state.StateDB) interfaces.TracerResult, bool) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting GetPluginTracer, but default PluginLoader has not been initialized")
		return nil, false
	}
	return GetPluginTracer(plugins.DefaultPluginLoader, name)
}
