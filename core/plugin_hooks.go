package core

import (
	"encoding/json"
	"math/big"
	"reflect"
	"time"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/plugins"
	"github.com/ethereum/go-ethereum/plugins/wrappers"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/openrelayxyz/plugeth-utils/core"
)

func PluginPreProcessBlock(pl *plugins.PluginLoader, block *types.Block) {
	fnList := pl.Lookup("PreProcessBlock", func(item interface{}) bool {
		_, ok := item.(func(core.Hash, uint64, []byte))
		return ok
	})
	encoded, _ := rlp.EncodeToBytes(block)
	for _, fni := range fnList {
		if fn, ok := fni.(func(core.Hash, uint64, []byte)); ok {
			fn(core.Hash(block.Hash()), block.NumberU64(), encoded)
		}
	}
}
func pluginPreProcessBlock(block *types.Block) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting PreProcessBlock, but default PluginLoader has not been initialized")
		return
	}
	PluginPreProcessBlock(plugins.DefaultPluginLoader, block) // TODO
}
func PluginPreProcessTransaction(pl *plugins.PluginLoader, tx *types.Transaction, block *types.Block, i int) {
	fnList := pl.Lookup("PreProcessTransaction", func(item interface{}) bool {
		_, ok := item.(func([]byte, core.Hash, core.Hash, int))
		return ok
	})
	txBytes, _ := tx.MarshalBinary()
	for _, fni := range fnList {
		if fn, ok := fni.(func([]byte, core.Hash, core.Hash, int)); ok {
			fn(txBytes, core.Hash(tx.Hash()), core.Hash(block.Hash()), i)
		}
	}
}
func pluginPreProcessTransaction(tx *types.Transaction, block *types.Block, i int) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting PreProcessTransaction, but default PluginLoader has not been initialized")
		return
	}
	PluginPreProcessTransaction(plugins.DefaultPluginLoader, tx, block, i)
}
func PluginBlockProcessingError(pl *plugins.PluginLoader, tx *types.Transaction, block *types.Block, err error) {
	fnList := pl.Lookup("BlockProcessingError", func(item interface{}) bool {
		_, ok := item.(func(core.Hash, core.Hash, error))
		return ok
	})
	for _, fni := range fnList {
		if fn, ok := fni.(func(core.Hash, core.Hash, error)); ok {
			fn(core.Hash(tx.Hash()), core.Hash(block.Hash()), err)
		}
	}
}
func pluginBlockProcessingError(tx *types.Transaction, block *types.Block, err error) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting BlockProcessingError, but default PluginLoader has not been initialized")
		return
	}
	PluginBlockProcessingError(plugins.DefaultPluginLoader, tx, block, err)
}
func PluginPostProcessTransaction(pl *plugins.PluginLoader, tx *types.Transaction, block *types.Block, i int, receipt *types.Receipt) {
	fnList := pl.Lookup("PostProcessTransaction", func(item interface{}) bool {
		_, ok := item.(func(core.Hash, core.Hash, int, []byte))
		return ok
	})
	receiptBytes, _ := json.Marshal(receipt)
	for _, fni := range fnList {
		if fn, ok := fni.(func(core.Hash, core.Hash, int, []byte)); ok {
			fn(core.Hash(tx.Hash()), core.Hash(block.Hash()), i, receiptBytes)
		}
	}
}
func pluginPostProcessTransaction(tx *types.Transaction, block *types.Block, i int, receipt *types.Receipt) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting PostProcessTransaction, but default PluginLoader has not been initialized")
		return
	}
	PluginPostProcessTransaction(plugins.DefaultPluginLoader, tx, block, i, receipt)
}
func PluginPostProcessBlock(pl *plugins.PluginLoader, block *types.Block) {
	fnList := pl.Lookup("PostProcessBlock", func(item interface{}) bool {
		_, ok := item.(func(core.Hash))
		return ok
	})
	for _, fni := range fnList {
		if fn, ok := fni.(func(core.Hash)); ok {
			fn(core.Hash(block.Hash()))
		}
	}
}
func pluginPostProcessBlock(block *types.Block) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting PostProcessBlock, but default PluginLoader has not been initialized")
		return
	}
	PluginPostProcessBlock(plugins.DefaultPluginLoader, block)
}
func PluginNewHead(pl *plugins.PluginLoader, block *types.Block, hash common.Hash, logs []*types.Log, td *big.Int) {
	fnList := pl.Lookup("NewHead", func(item interface{}) bool {
		_, ok := item.(func([]byte, core.Hash, [][]byte, *big.Int))
		return ok
	})
	blockBytes, _ := rlp.EncodeToBytes(block)
	logBytes := make([][]byte, len(logs))
	for i, l := range logs {
		logBytes[i], _ = rlp.EncodeToBytes(l)
	}
	for _, fni := range fnList {
		if fn, ok := fni.(func([]byte, core.Hash, [][]byte, *big.Int)); ok {
			fn(blockBytes, core.Hash(hash), logBytes, td)
		}
	}
}
func pluginNewHead(block *types.Block, hash common.Hash, logs []*types.Log, td *big.Int) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting NewHead, but default PluginLoader has not been initialized")
		return
	}
	PluginNewHead(plugins.DefaultPluginLoader, block, hash, logs, td)
}

func PluginNewSideBlock(pl *plugins.PluginLoader, block *types.Block, hash common.Hash, logs []*types.Log) {
	fnList := pl.Lookup("NewSideBlock", func(item interface{}) bool {
		_, ok := item.(func([]byte, core.Hash, [][]byte))
		return ok
	})
	blockBytes, _ := rlp.EncodeToBytes(block)
	logBytes := make([][]byte, len(logs))
	for i, l := range logs {
		logBytes[i], _ = rlp.EncodeToBytes(l)
	}
	for _, fni := range fnList {
		if fn, ok := fni.(func([]byte, core.Hash, [][]byte)); ok {
			fn(blockBytes, core.Hash(hash), logBytes)
		}
	}
}
func pluginNewSideBlock(block *types.Block, hash common.Hash, logs []*types.Log) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting NewSideBlock, but default PluginLoader has not been initialized")
		return
	}
	PluginNewSideBlock(plugins.DefaultPluginLoader, block, hash, logs)
}

func PluginReorg(pl *plugins.PluginLoader, commonBlock *types.Block, oldChain, newChain types.Blocks) {
	fnList := pl.Lookup("Reorg", func(item interface{}) bool {
		_, ok := item.(func(core.Hash, []core.Hash, []core.Hash))
		return ok
	})
	oldChainHashes := make([]core.Hash, len(oldChain))
	for i, block := range oldChain {
		oldChainHashes[i] = core.Hash(block.Hash())
	}
	newChainHashes := make([]core.Hash, len(newChain))
	for i, block := range newChain {
		newChainHashes[i] = core.Hash(block.Hash())
	}
	for _, fni := range fnList {
		if fn, ok := fni.(func(core.Hash, []core.Hash, []core.Hash)); ok {
			fn(core.Hash(commonBlock.Hash()), oldChainHashes, newChainHashes)
		}
	}
}
func pluginReorg(commonBlock *types.Block, oldChain, newChain types.Blocks) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting Reorg, but default PluginLoader has not been initialized")
		return
	}
	PluginReorg(plugins.DefaultPluginLoader, commonBlock, oldChain, newChain)
}

type PreTracer interface {
	CapturePreStart(from common.Address, to *common.Address, input []byte, gas uint64, value *big.Int)

}

type metaTracer struct {
	tracers []core.BlockTracer
}

func (mt *metaTracer) PreProcessBlock(block *types.Block) {
	if len(mt.tracers) == 0 { return }
	blockHash := core.Hash(block.Hash())
	blockNumber := block.NumberU64()
	encoded, _ := rlp.EncodeToBytes(block)
	for _, tracer := range mt.tracers {
		tracer.PreProcessBlock(blockHash, blockNumber, encoded)
	}
}
func (mt *metaTracer) PreProcessTransaction(tx *types.Transaction, block *types.Block, i int) {
	if len(mt.tracers) == 0 { return }
	blockHash := core.Hash(block.Hash())
	transactionHash := core.Hash(tx.Hash())
	for _, tracer := range mt.tracers {
		tracer.PreProcessTransaction(transactionHash, blockHash, i)
	}
}
func (mt *metaTracer) BlockProcessingError(tx *types.Transaction, block *types.Block, err error) {
	if len(mt.tracers) == 0 { return }
	blockHash := core.Hash(block.Hash())
	transactionHash := core.Hash(tx.Hash())
	for _, tracer := range mt.tracers {
		tracer.BlockProcessingError(transactionHash, blockHash, err)
	}
}
func (mt *metaTracer) PostProcessTransaction(tx *types.Transaction, block *types.Block, i int, receipt *types.Receipt) {
	if len(mt.tracers) == 0 { return }
	blockHash := core.Hash(block.Hash())
	transactionHash := core.Hash(tx.Hash())
	receiptBytes, _ := json.Marshal(receipt)
	for _, tracer := range mt.tracers {
		tracer.PostProcessTransaction(transactionHash, blockHash, i, receiptBytes)
	}
}
func (mt *metaTracer) PostProcessBlock(block *types.Block) {
	if len(mt.tracers) == 0 { return }
	blockHash := core.Hash(block.Hash())
	for _, tracer := range mt.tracers {
		tracer.PostProcessBlock(blockHash)
	}
}
func (mt *metaTracer) CaptureStart(env *vm.EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
	for _, tracer := range mt.tracers {
		tracer.CaptureStart(core.Address(from), core.Address(to), create, input, gas, value)
	}
}
func (mt *metaTracer) CaptureState(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, rData []byte, depth int, err error) {
	for _, tracer := range mt.tracers {
		tracer.CaptureState(pc, core.OpCode(op), gas, cost, wrappers.NewWrappedScopeContext(scope), rData, depth, err)
	}
}
func (mt *metaTracer) CaptureFault(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, depth int, err error) {
	for _, tracer := range mt.tracers {
		tracer.CaptureFault(pc, core.OpCode(op), gas, cost, wrappers.NewWrappedScopeContext(scope), depth, err)
	}
}
func (mt *metaTracer) CaptureEnd(output []byte, gasUsed uint64, err error) {
	for _, tracer := range mt.tracers {
		tracer.CaptureEnd(output, gasUsed, err)
	}
}

func (mt *metaTracer) CaptureEnter(typ vm.OpCode, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int) {
	for _, tracer := range mt.tracers {
		tracer.CaptureEnter(core.OpCode(typ), core.Address(from), core.Address(to), input, gas, value)
	}
}

func (mt *metaTracer) CaptureExit(output []byte, gasUsed uint64, err error) {
	for _, tracer := range mt.tracers {
		tracer.CaptureExit(output, gasUsed, err)
	}
}

func (mt metaTracer) CaptureTxStart (gasLimit uint64) {}

func (mt metaTracer) CaptureTxEnd (restGas uint64) {}

func PluginGetBlockTracer(pl *plugins.PluginLoader, hash common.Hash, statedb *state.StateDB) (*metaTracer, bool) {
	//look for a function that takes whatever the ctx provides and statedb and returns a core.blocktracer append into meta tracer
	tracerList := plugins.Lookup("GetLiveTracer", func(item interface{}) bool {
		_, ok := item.(func(core.Hash, core.StateDB) core.BlockTracer)
		log.Info("Item is LiveTracer", "ok", ok, "type", reflect.TypeOf(item))
		return ok
	})
	mt := &metaTracer{tracers: []core.BlockTracer{}}
	for _, tracer := range tracerList {
		if v, ok := tracer.(func(core.Hash, core.StateDB) core.BlockTracer); ok {
			bt := v(core.Hash(hash), wrappers.NewWrappedStateDB(statedb))
			if bt != nil {
				mt.tracers = append(mt.tracers, bt)
			}
		}
	}
	return mt, (len(mt.tracers) > 0)
}
func pluginGetBlockTracer(hash common.Hash, statedb *state.StateDB) (*metaTracer, bool) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting GetBlockTracer, but default PluginLoader has not been initialized")
		return &metaTracer{}, false
	}
	return PluginGetBlockTracer(plugins.DefaultPluginLoader, hash, statedb)
}

func PluginSetTrieFlushIntervalClone(pl *plugins.PluginLoader, flushInterval time.Duration) time.Duration {
	fnList := pl.Lookup("SetTrieFlushIntervalClone", func(item interface{}) bool{
		_, ok := item.(func(time.Duration) time.Duration)
		return ok
	})
	var snc sync.Once
	if len(fnList) > 0 {
		snc.Do(func() {log.Warn("The blockChain flushInterval value is being accessed by multiple plugins")})
	}
	for _, fni := range fnList {
		if fn, ok := fni.(func(time.Duration) time.Duration); ok {
			flushInterval = fn(flushInterval) 
		}
	}
	return flushInterval
}

func pluginSetTrieFlushIntervalClone(flushInterval time.Duration) time.Duration {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting setTreiFlushIntervalClone, but default PluginLoader has not been initialized")
		return flushInterval
	}
	return PluginSetTrieFlushIntervalClone(plugins.DefaultPluginLoader, flushInterval)
}