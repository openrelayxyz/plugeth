package main

import(
	"errors"
	"math/big"

	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted"
	"github.com/openrelayxyz/plugeth-utils/restricted/types"
	"github.com/openrelayxyz/plugeth-utils/restricted/hasher"
	"github.com/openrelayxyz/plugeth-utils/restricted/consensus"
)

var (
	pl      core.PluginLoader
	backend restricted.Backend
	log     core.Logger
	events  core.Feed
)

var httpApiFlagName = "http.api"

func Initialize(ctx core.Context, loader core.PluginLoader, logger core.Logger) { 
	pl = loader
	events = pl.GetFeed()
	log = logger
	v := ctx.String(httpApiFlagName)
	if v != "" {
		ctx.Set(httpApiFlagName, v+",plugeth")
	} else {
		ctx.Set(httpApiFlagName, "eth,net,web3,plugeth")
		log.Info("Loaded consensus engine plugin")
	}
}

type engine struct {
}

func (e *engine) Author(header *types.Header) (core.Address, error) {
	return header.Coinbase, nil
}
func (e *engine) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, seal bool) error {
	return nil
}
func (e *engine) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	quit := make(chan struct{})
	err := make(chan error)
	go func () {
		for i, h := range headers {
			select {
			case <-quit:
				return 
			case err<- e.VerifyHeader(chain, h, seals[i]):
			}
		} 
	} ()
	return quit, err
}
func (e *engine) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	return nil
}
func (e *engine) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	header.Difficulty = new(big.Int).SetUint64(123456789)
	header.UncleHash = core.HexToHash("1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347")
	return nil
}
func (e *engine) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state core.RWStateDB, txs []*types.Transaction,uncles []*types.Header, withdrawals []*types.Withdrawal) {
}
func (e *engine) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state core.RWStateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt, withdrawals []*types.Withdrawal) (*types.Block, error) {
	header.Root = state.IntermediateRoot(false)
	hasher := hasher.NewStackTrie(nil)
	block := types.NewBlockWithWithdrawals(header, txs, uncles, receipts, withdrawals, hasher)
	return block, nil

}
func (e *engine) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	if len(block.Transactions()) == 0 {
		return errors.New("sealing paused while waiting for transactions")
	}
	go func () {
		results <- block 
		close(results)
	} ()
	// TO DO: the stop channel will need to be addressed in a non test case scenerio
	return nil
}
func (e *engine) SealHash(header *types.Header) core.Hash {
	return header.Hash()
}
func (e *engine) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	return new(big.Int).SetUint64(uint64(123456789))
}
func (e *engine) APIs(chain consensus.ChainHeaderReader) []core.API {
	return []core.API{}
}
func (e *engine) Close() error {
	return nil
}

func CreateEngine() consensus.Engine {
	return &engine{}
}