package tracers

import (
	"reflect"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/plugins"
	"github.com/ethereum/go-ethereum/plugins/interfaces"
	"github.com/ethereum/go-ethereum/plugins/wrappers"
	"github.com/openrelayxyz/plugeth-utils/core"
)

func GetPluginTracer(pl *plugins.PluginLoader, name string) (func(*state.StateDB, vm.BlockContext) interfaces.TracerResult, bool) {
	tracers := pl.Lookup("Tracers", func(item interface{}) bool {
		_, ok := item.(*map[string]func(core.StateDB) core.TracerResult)
		_, ok2 := item.(*map[string]func(core.StateDB, core.BlockContext) core.TracerResult)
		if !(ok || ok2) {
			log.Warn("Found tracer that did not match type", "tracer", reflect.TypeOf(item))
		}
		return ok || ok2
	})

	for _, tmap := range tracers {
		if tracerMap, ok := tmap.(*map[string]func(core.StateDB) core.TracerResult); ok {
			if tracer, ok := (*tracerMap)[name]; ok {
				return func(sdb *state.StateDB, vmctx vm.BlockContext) interfaces.TracerResult {
					return wrappers.NewWrappedTracer(tracer(wrappers.NewWrappedStateDB(sdb)))
				}, true
			}
		}
		if tracerMap, ok := tmap.(*map[string]func(core.StateDB, core.BlockContext) core.TracerResult); ok {
			if tracer, ok := (*tracerMap)[name]; ok {
				return func(sdb *state.StateDB, vmctx vm.BlockContext) interfaces.TracerResult {
					return wrappers.NewWrappedTracer(tracer(wrappers.NewWrappedStateDB(sdb), core.BlockContext{
						Coinbase:    core.Address(vmctx.Coinbase),
						GasLimit:    vmctx.GasLimit,
						BlockNumber: vmctx.BlockNumber,
						Time:        vmctx.Time,
						Difficulty:  vmctx.Difficulty,
						BaseFee:     vmctx.BaseFee,
					}))
				}, true
			}
		}
	}
	log.Info("Tracer not found", "name", name, "tracers", len(tracers))
	return nil, false
}

func getPluginTracer(name string) (func(*state.StateDB, vm.BlockContext) interfaces.TracerResult, bool) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting GetPluginTracer, but default PluginLoader has not been initialized")
		return nil, false
	}
	return GetPluginTracer(plugins.DefaultPluginLoader, name)
}
