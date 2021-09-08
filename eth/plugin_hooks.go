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

// func PluginCreateConsensusEngine(pl *plugins.PluginLoader, stack *node.Node, chainConfig *params.ChainConfig, config *ethash.Config, notify []string, noverify bool, db ethdb.Database) consensus.Engine {
//   fnList := pl.Lookup("CreateConsensusEngine", func(item interface{}) bool {
//     _, ok := item.(func(*node.Node, *params.ChainConfig, *ethash.Config, []string, bool, ethdb.Database) consensus.Engine)
//     return ok
//   })
//   for _, fni := range fnList {
//     if fn, ok := fni.(func(*node.Node, *params.ChainConfig, *ethash.Config, []string, bool, ethdb.Database) consensus.Engine); ok {
//       return fn(stack, chainConfig, config, notify, noverify, db)
//     }
//   }
//   return ethconfig.CreateConsensusEngine(stack, chainConfig, config, notify, noverify, db)
// }
//
// func pluginCreateConsensusEngine(stack *node.Node, chainConfig *params.ChainConfig, config *ethash.Config, notify []string, noverify bool, db ethdb.Database) consensus.Engine {
//   if plugins.DefaultPluginLoader == nil {
// 		log.Warn("Attempting CreateConsensusEngine, but default PluginLoader has not been initialized")
// 		return ethconfig.CreateConsensusEngine(stack, chainConfig, config, notify, noverify, db)
// 	}
//   return PluginCreateConsensusEngine(plugins.DefaultPluginLoader, stack, chainConfig, config, notify, noverify, db)
// }

// TODO (philip): Translate to core.TracerResult instead of vm.Tracer, with appropriate type adjustments (let me know if this one is too hard)
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
		_, ok := item.(*vm.Tracer)
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
		cfg.Tracer = mt //I think this means we will need a vm.config wrapper although confugure doesnt sound very passive
	} else {
		log.Warn("Module is not tracer")
	}

	fnList := plugins.Lookup("UpdateBlockchainVMConfig", func(item interface{}) bool {
		_, ok := item.(func(*vm.Config))
		return ok
	})
	for _, fni := range fnList {
		if fn, ok := fni.(func(*vm.Config)); ok {
			fn(cfg)
			return
		}
	}
}

func pluginUpdateBlockchainVMConfig(cfg *vm.Config) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting CreateConsensusEngine, but default PluginLoader has not been initialized")
		return
	}
	PluginUpdateBlockchainVMConfig(plugins.DefaultPluginLoader, cfg)
}
