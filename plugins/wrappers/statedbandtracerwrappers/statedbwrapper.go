package statedbandtracerwrappers

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/openrelayxyz/plugeth-utils/core"
)

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

// HasSuicided(Address) bool
func (w *WrappedStateDB) HasSuicided(addr core.Address) bool { // I figured we'd skip some of the future labor and update the name now
	return w.s.HasSuicided(common.Address(addr))
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
