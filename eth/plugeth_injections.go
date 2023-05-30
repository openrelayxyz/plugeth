package eth

import (
	"github.com/ethereum/go-ethereum/core"
)

func (b *EthAPIBackend) BlockChain() *core.BlockChain {
	return b.eth.BlockChain()
}