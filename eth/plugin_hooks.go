package eth

import (
	"math/big"
	"reflect"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/plugins"
	"github.com/ethereum/go-ethereum/plugins/wrappers"
	"github.com/openrelayxyz/plugeth-utils/core"
)

type metaTracer struct {
	tracers []core.TracerResult
}

func (mt *metaTracer) CaptureStart(env *vm.EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
	for _, tracer := range mt.tracers {
		tracer.CaptureStart(core.Address(from), core.Address(to), create, input, gas, value)
	}
}
func (mt *metaTracer) CaptureState(env *vm.EVM, pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, rData []byte, depth int, err error) {
	for _, tracer := range mt.tracers {
		tracer.CaptureState(pc, core.OpCode(op), gas, cost, wrappers.NewWrappedScopeContext(scope), rData, depth, err)
	}
}
func (mt *metaTracer) CaptureFault(env *vm.EVM, pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, depth int, err error) {
	for _, tracer := range mt.tracers {
		tracer.CaptureFault(pc, core.OpCode(op), gas, cost, wrappers.NewWrappedScopeContext(scope), depth, err)
	}
}
func (mt *metaTracer) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) {
	for _, tracer := range mt.tracers {
		tracer.CaptureEnd(output, gasUsed, t, err)
	}
}

func PluginUpdateBlockchainVMConfig(pl *plugins.PluginLoader, cfg *vm.Config) {
	tracerList := plugins.Lookup("LiveTracer", func(item interface{}) bool {
		_, ok := item.(core.TracerResult)
		log.Info("Item is LiveTracer", "ok", ok, "type", reflect.TypeOf(item))
		return ok
	})
	if len(tracerList) > 0 {
		mt := &metaTracer{tracers: []core.TracerResult{}}
		for _, tracer := range tracerList {
			if v, ok := tracer.(core.TracerResult); ok {
				log.Info("LiveTracer registered")
				mt.tracers = append(mt.tracers, v)
			} else {
				log.Info("Item is not tracer")
			}
		}
		cfg.Debug = true
		cfg.Tracer = mt
	} else {
		log.Warn("Module is not tracer")
	}

}

func pluginUpdateBlockchainVMConfig(cfg *vm.Config) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting CreateConsensusEngine, but default PluginLoader has not been initialized")
		return
	}
	PluginUpdateBlockchainVMConfig(plugins.DefaultPluginLoader, cfg)
}
