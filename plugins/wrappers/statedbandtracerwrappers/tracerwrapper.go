package statedbandtracerwrappers

import (
	"math/big"
	"time"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/openrelayxyz/plugeth-utils/core"
)

type WrappedScopeContext struct {
	s *vm.ScopeContext
}

func NewWrappedScopeContext(s *vm.ScopeContext) *WrappedScopeContext {
	return &WrappedScopeContext{s}
}

func (w *WrappedScopeContext) Memory() core.Memory {
	return w.s.Memory
}

func (w *WrappedScopeContext) Stack() core.Stack {
	return w.s.Stack
}

func (w *WrappedScopeContext) Contract() core.Contract {
	return &WrappedContract{w.s.Contract}
}

type WrappedTracer struct {
	r core.TracerResult
}

type WrappedContract struct {
	c *vm.Contract
}

func (w *WrappedContract) AsDelegate() core.Contract {
	return &WrappedContract{w.c.AsDelegate()}
}

func (w *WrappedContract) GetOp(n uint64) core.OpCode {
	return core.OpCode(w.c.GetOp(n))
}

func (w *WrappedContract) GetByte(n uint64) byte {
	return byte(w.c.GetOp(n))
}

func (w *WrappedContract) Caller() core.Address {
	return core.Address(w.c.Caller())
}

func (w *WrappedContract) Address() core.Address {
	return core.Address(w.c.Address())
}

func (w *WrappedContract) Value() *big.Int {
	return w.c.Value()
}

func (w *WrappedContract) Input() []byte {
	return w.c.Input
}

func (w *WrappedContract) Code() []byte {
	return w.c.Code
}

// added UseGas bc compiler compained without it. Should investigate if the false return with effect performance.
// take this out of core.interface
func (w *WrappedContract) UseGas(gas uint64) (ok bool) {

	return false
}

func NewWrappedTracer(r core.TracerResult) *WrappedTracer {
	return &WrappedTracer{r}
}
func (w WrappedTracer) CapturePreStart(from common.Address, to *common.Address, input []byte, gas uint64, value *big.Int) {
	if v, ok := w.r.(core.PreTracer); ok {
	v.CapturePreStart(core.Address(from), (*core.Address)(to), input, gas, value)}
}
func (w WrappedTracer) CaptureStart(env *vm.EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
	w.r.CaptureStart(core.Address(from), core.Address(to), create, input, gas, value)
}
func (w WrappedTracer) CaptureState(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, rData []byte, depth int, err error) {
	w.r.CaptureState(pc, core.OpCode(op), gas, cost, &WrappedScopeContext{scope}, rData, depth, err)
}
func (w WrappedTracer) CaptureEnter(typ vm.OpCode, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int) {
	w.r.CaptureEnter(core.OpCode(typ), core.Address(from), core.Address(to), input, gas, value)
}
func (w WrappedTracer) CaptureExit(output []byte, gasUsed uint64, err error) {
	w.r.CaptureExit(output, gasUsed, err)
}
func (w WrappedTracer) CaptureFault(pc uint64, op vm.OpCode, gas, cost uint64, scope *vm.ScopeContext, depth int, err error) {
	w.r.CaptureFault(pc, core.OpCode(op), gas, cost, &WrappedScopeContext{scope}, depth, err)
}
func (w WrappedTracer) CaptureEnd(output []byte, gasUsed uint64, t time.Duration, err error) {
	w.r.CaptureEnd(output, gasUsed, t, err)
}
func (w WrappedTracer) GetResult() (json.RawMessage, error) {
	data, err := w.r.Result()
	if err != nil { return nil, err}
	result, err := json.Marshal(data)
	return json.RawMessage(result), err
}

func (w WrappedTracer) CaptureTxStart (gasLimit uint64) {}

func (w WrappedTracer) CaptureTxEnd (restGas uint64) {}

func (w WrappedTracer) Stop(err error) {}
