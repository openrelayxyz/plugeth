package wrappers

import (
	"math/big"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/node"
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

type WrappedTracer struct {
	r core.TracerResult
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
// passing zero as a dummy value is foundation PluGeth only, it is being done to preserve compatability with other networks
func (w WrappedTracer) CaptureEnd(output []byte, gasUsed uint64, err error) {
	w.r.CaptureEnd(output, gasUsed, 0, err)
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


type WrappedStateDB struct {
	s *state.StateDB
}

func NewWrappedStateDB(d *state.StateDB) *WrappedStateDB {
	return &WrappedStateDB{d}
}

// GetBalance(Address) *big.Int
func (w *WrappedStateDB) GetBalance(addr core.Address) *big.Int {
	return w.s.GetBalance(common.Address(addr))
}

// GetNonce(Address) uint64
func (w *WrappedStateDB) GetNonce(addr core.Address) uint64 {
	return w.s.GetNonce(common.Address(addr))
}

// GetCodeHash(Address) Hash
func (w *WrappedStateDB) GetCodeHash(addr core.Address) core.Hash {
	return core.Hash(w.s.GetCodeHash(common.Address(addr)))
} // sort this out

// GetCode(Address) []byte
func (w *WrappedStateDB) GetCode(addr core.Address) []byte {
	return w.s.GetCode(common.Address(addr))
}

// GetCodeSize(Address) int
func (w *WrappedStateDB) GetCodeSize(addr core.Address) int {
	return w.s.GetCodeSize(common.Address(addr))
}

//GetRefund() uint64
func (w *WrappedStateDB) GetRefund() uint64 { //are we sure we want to include this? getting a refund seems like changing state
	return w.s.GetRefund()
}

// GetCommittedState(Address, Hash) Hash
func (w *WrappedStateDB) GetCommittedState(addr core.Address, hsh core.Hash) core.Hash {
	return core.Hash(w.s.GetCommittedState(common.Address(addr), common.Hash(hsh)))
}

// GetState(Address, Hash) Hash
func (w *WrappedStateDB) GetState(addr core.Address, hsh core.Hash) core.Hash {
	return core.Hash(w.s.GetState(common.Address(addr), common.Hash(hsh)))
}

func (w *WrappedStateDB) HasSuicided(addr core.Address) bool { 
	return w.s.HasSelfDestructed(common.Address(addr))
}

// // Exist reports whether the given account exists in state.
// // Notably this should also return true for suicided accounts.
// Exist(Address) bool
func (w *WrappedStateDB) Exist(addr core.Address) bool {
	return w.s.Exist(common.Address(addr))
}

// // Empty returns whether the given account is empty. Empty
// // is defined according to EIP161 (balance = nonce = code = 0).
// Empty(Address) bool
func (w *WrappedStateDB) Empty(addr core.Address) bool {
	return w.s.Empty(common.Address(addr))
}

// AddressInAccessList(addr Address) bool
func (w *WrappedStateDB) AddressInAccessList(addr core.Address) bool {
	return w.s.AddressInAccessList(common.Address(addr))
}

// SlotInAccessList(addr Address, slot Hash) (addressOk bool, slotOk bool)
func (w *WrappedStateDB) SlotInAccessList(addr core.Address, slot core.Hash) (addressOK, slotOk bool) {
	return w.s.SlotInAccessList(common.Address(addr), common.Hash(slot))
}

// IntermediateRoot(deleteEmptyObjects bool) common.Hash 
func (w *WrappedStateDB) IntermediateRoot(deleteEmptyObjects bool) core.Hash {
	return core.Hash(w.s.IntermediateRoot(deleteEmptyObjects))
}
	

type Node struct {
	n *node.Node
}

func NewNode(n *node.Node) *Node {
	return &Node{n}
}

func (n *Node) Server() core.Server {
	return n.n.Server()
}

func (n *Node) DataDir() string {
	return n.n.DataDir()
}
func (n *Node) InstanceDir() string {
	return n.n.InstanceDir()
}
func (n *Node) IPCEndpoint() string {
	return n.n.IPCEndpoint()
}
func (n *Node) HTTPEndpoint() string {
	return n.n.HTTPEndpoint()
}
func (n *Node) WSEndpoint() string {
	return n.n.WSEndpoint()
}
func (n *Node) ResolvePath(x string) string {
	return n.n.ResolvePath(x)
}
func (n *Node) Attach() (core.Client, error) {
	return n.n.Attach(), nil
}
func (n *Node) Close() error {
	return n.n.Close()
}

type WrappedBlockContext struct {
	b vm.BlockContext
}

// type WrappedBlockContext vm.BlockContext

func NewWrappedBlockContext(c vm.BlockContext) *WrappedBlockContext {
	return &WrappedBlockContext{c}
}
func (w *WrappedBlockContext) Coinbase() core.Address {
	return core.Address(w.b.Coinbase)
}
func (w *WrappedBlockContext) GasLimit() uint64 {
	return w.b.GasLimit
}
func (w *WrappedBlockContext) BlockNumber() *big.Int {
	return w.b.BlockNumber
}
func (w *WrappedBlockContext) Time() *big.Int {
	return new(big.Int).SetInt64(int64(w.b.Time))
}
func (w *WrappedBlockContext) Difficulty() *big.Int {
	return w.b.Difficulty
}
func (w *WrappedBlockContext) BaseFee() *big.Int {
	return w.b.BaseFee
}
