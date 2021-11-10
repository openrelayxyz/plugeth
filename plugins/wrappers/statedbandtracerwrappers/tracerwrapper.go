package statedbandtracerwrappers

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/openrelayxyz/plugeth-utils/core"
)

type WrappedTracer struct {
	r core.TracerResult
}

func NewWrappedTracer(r core.TracerResult) *WrappedTracer {
	return &WrappedTracer{r}
}
func (w WrappedTracer) CaptureStart(env *vm.EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
	w.r.CaptureStart(core.Address(from), core.Address(to), create, input, gas, value)
}
func (w WrappedTracer) CaptureState(env *vm.EVM, pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, rData []byte, depth int, err error) {
	w.r.CaptureState(pc, core.OpCode(op), gas, cost, &WrappedScopeContext{scope}, rData, depth, err)
}
func (w WrappedTracer) CaptureFault(env *vm.EVM, pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, depth int, err error) {
	w.r.CaptureFault(pc, core.OpCode(op), gas, cost, &WrappedScopeContext{scope}, depth, err)
}
func (w WrappedTracer) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) {
	w.r.CaptureEnd(output, gasUsed, t, err)
}
func (w WrappedTracer) CaptureEnter(typ vm.OpCode, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int) {
	w.r.CaptureEnter(core.OpCode(typ), core.Address(from), core.Address(to), input, gas, value)
}
func (w WrappedTracer) CaptureExit(output []byte, gasUsed uint64, err error) {
	w.r.CaptureExit(output, gasUsed, err)
}
func (w WrappedTracer) GetResult() (interface{}, error) {
	return w.r.Result()
}
