package backendwrapper

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/trie"

	"github.com/openrelayxyz/plugeth-utils/core"
)

type WrappedTrie struct {
	t state.Trie
}

func NewWrappedTrie(t state.Trie) core.Trie {
	return &WrappedTrie{t}
}

func (t *WrappedTrie) GetKey(b []byte) []byte {
	return t.t.GetKey(b)
}

func (t *WrappedTrie) GetAccount(address core.Address) (*core.StateAccount, error) {
	act, err := t.t.GetAccount(common.Address(address))
	if err != nil {
		return nil, err
	}
	return &core.StateAccount{
		Nonce: act.Nonce,
		Balance: act.Balance,
		Root: core.Hash(act.Root),
		CodeHash: act.CodeHash,
	}, nil
}

func (t *WrappedTrie) Hash() core.Hash {
	return core.Hash(t.t.Hash())
}

func (t *WrappedTrie) NodeIterator(startKey []byte) core.NodeIterator {
	itr := t.t.NodeIterator(startKey)
	return &WrappedNodeIterator{itr}
}

func (t *WrappedTrie) Prove(key []byte, fromLevel uint, proofDb core.KeyValueWriter) error {
	return nil
}

type WrappedNodeIterator struct {
	n trie.NodeIterator
}

func (n WrappedNodeIterator) Next(b bool) bool {
	return n.n.Next(b)
}

func (n WrappedNodeIterator) Error() error {
	return n.n.Error()
}

func (n WrappedNodeIterator) Hash() core.Hash {
	return core.Hash(n.n.Hash())
}

func (n WrappedNodeIterator) Parent() core.Hash {
	return core.Hash(n.n.Parent())
}

func (n WrappedNodeIterator) Path() []byte {
	return n.n.Path()
}

func (n WrappedNodeIterator) NodeBlob() []byte {
	return n.n.NodeBlob()
}

func (n WrappedNodeIterator) Leaf() bool {
	return n.n.Leaf()
}

func (n WrappedNodeIterator) LeafKey() []byte {
	return n.n.LeafKey()
}

func (n WrappedNodeIterator) LeafBlob() []byte {
	return n.n.LeafBlob()
}

func (n WrappedNodeIterator) LeafProof() [][]byte {
	return n.n.LeafProof()
}

func (n WrappedNodeIterator) AddResolver(c core.NodeResolver) {
	n.n.AddResolver(WrappedNodeResolver(c))
}

func WrappedNodeResolver(fn core.NodeResolver) trie.NodeResolver {
	return func(owner common.Hash, path []byte, hash common.Hash) []byte {
		return fn(core.Hash(owner), path, core.Hash(hash) )
	}
}