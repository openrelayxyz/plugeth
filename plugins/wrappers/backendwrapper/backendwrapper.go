package backendwrapper

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	gcore "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/trie"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted"
	"github.com/openrelayxyz/plugeth-utils/restricted/params"
)

type Backend struct {
	b               ethapi.Backend
	newTxsFeed      event.Feed
	newTxsOnce      sync.Once
	chainFeed       event.Feed
	chainOnce       sync.Once
	chainHeadFeed   event.Feed
	chainHeadOnce   sync.Once
	chainSideFeed   event.Feed
	chainSideOnce   sync.Once
	logsFeed        event.Feed
	logsOnce        sync.Once
	pendingLogsFeed event.Feed
	pendingLogsOnce sync.Once
	removedLogsFeed event.Feed
	removedLogsOnce sync.Once
	chainConfig     *params.ChainConfig
}

func NewBackend(b ethapi.Backend) *Backend {
	return &Backend{b: b}
}

func (b *Backend) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	return b.b.SuggestGasTipCap(ctx)
}
func (b *Backend) ChainDb() restricted.Database {
	return &dbWrapper{b.b.ChainDb()}
}
func (b *Backend) ExtRPCEnabled() bool {
	return b.b.ExtRPCEnabled()
}
func (b *Backend) RPCGasCap() uint64 {
	return b.b.RPCGasCap()
}
func (b *Backend) RPCTxFeeCap() float64 {
	return b.b.RPCTxFeeCap()
}
func (b *Backend) UnprotectedAllowed() bool {
	return b.b.UnprotectedAllowed()
}
func (b *Backend) SetHead(number uint64) {
	b.b.SetHead(number)
}
func (b *Backend) HeaderByNumber(ctx context.Context, number int64) ([]byte, error) {
	header, err := b.b.HeaderByNumber(ctx, rpc.BlockNumber(number))
	if err != nil {
		return nil, err
	}
	return rlp.EncodeToBytes(header)
}
func (b *Backend) HeaderByHash(ctx context.Context, hash core.Hash) ([]byte, error) {
	header, err := b.b.HeaderByHash(ctx, common.Hash(hash))
	if err != nil {
		return nil, err
	}
	return rlp.EncodeToBytes(header)
}
func (b *Backend) CurrentHeader() []byte {
	ret, _ := rlp.EncodeToBytes(b.b.CurrentHeader())
	return ret
}
func (b *Backend) CurrentBlock() []byte {
	ret, _ := rlp.EncodeToBytes(b.b.CurrentBlock())
	return ret
}
func (b *Backend) BlockByNumber(ctx context.Context, number int64) ([]byte, error) {
	block, err := b.b.BlockByNumber(ctx, rpc.BlockNumber(number))
	if err != nil {
		return nil, err
	}
	return rlp.EncodeToBytes(block)
}
func (b *Backend) BlockByHash(ctx context.Context, hash core.Hash) ([]byte, error) {
	block, err := b.b.BlockByHash(ctx, common.Hash(hash))
	if err != nil {
		return nil, err
	}
	return rlp.EncodeToBytes(block)
}
func (b *Backend) GetReceipts(ctx context.Context, hash core.Hash) ([]byte, error) {
	receipts, err := b.b.GetReceipts(ctx, common.Hash(hash))
	if err != nil {
		return nil, err
	}
	return json.Marshal(receipts)
}
func (b *Backend) GetTd(ctx context.Context, hash core.Hash) *big.Int {
	return b.b.GetTd(ctx, common.Hash(hash))
}
func (b *Backend) SendTx(ctx context.Context, signedTx []byte) error {
	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(signedTx); err != nil {
		return err
	}
	return b.b.SendTx(ctx, tx)
}
func (b *Backend) GetTransaction(ctx context.Context, txHash core.Hash) ([]byte, core.Hash, uint64, uint64, error) { // RLP Encoded transaction {
	tx, blockHash, blockNumber, index, err := b.b.GetTransaction(ctx, common.Hash(txHash))
	if err != nil {
		return nil, core.Hash(blockHash), blockNumber, index, err
	}
	enc, err := tx.MarshalBinary()
	return enc, core.Hash(blockHash), blockNumber, index, err
}
func (b *Backend) GetPoolTransactions() ([][]byte, error) {
	txs, err := b.b.GetPoolTransactions()
	if err != nil {
		return nil, err
	}
	results := make([][]byte, len(txs))
	for i, tx := range txs {
		results[i], _ = rlp.EncodeToBytes(tx)
	}
	return results, nil
}
func (b *Backend) GetPoolTransaction(txHash core.Hash) []byte {
	tx := b.b.GetPoolTransaction(common.Hash(txHash))
	if tx == nil {
		return []byte{}
	}
	enc, _ := rlp.EncodeToBytes(tx)
	return enc
}
func (b *Backend) GetPoolNonce(ctx context.Context, addr core.Address) (uint64, error) {
	return b.b.GetPoolNonce(ctx, common.Address(addr))
}
func (b *Backend) Stats() (pending int, queued int) {
	return b.b.Stats()
}
func (b *Backend) TxPoolContent() (map[core.Address][][]byte, map[core.Address][][]byte) {
	pending, queued := b.b.TxPoolContent()
	trpending, trqueued := make(map[core.Address][][]byte), make(map[core.Address][][]byte)
	for k, v := range pending {
		trpending[core.Address(k)] = make([][]byte, len(v))
		for i, tx := range v {
			trpending[core.Address(k)][i], _ = tx.MarshalBinary()
		}
	}
	for k, v := range queued {
		trqueued[core.Address(k)] = make([][]byte, len(v))
		for i, tx := range v {
			trpending[core.Address(k)][i], _ = tx.MarshalBinary()
		}
	}
	return trpending, trqueued
} // RLP encoded transactions
func (b *Backend) BloomStatus() (uint64, uint64) {
	return b.b.BloomStatus()
}
func (b *Backend) GetLogs(ctx context.Context, blockHash core.Hash) ([][]byte, error) {
	header, _ := b.b.HeaderByHash(ctx, common.Hash(blockHash))
	logs, err := b.b.GetLogs(ctx, common.Hash(blockHash), uint64(header.Number.Uint64()))
	if err != nil {
		return nil, err
	}
	encLogs := make([][]byte, len(logs))
	for i, log := range logs {
		encLogs[i], _ = rlp.EncodeToBytes(log)
	}
	return encLogs, nil
} // []RLP encoded logs

type dli interface {
	SyncProgress() ethereum.SyncProgress
}

type dl struct {
	dl dli
}

type progress struct {
	p ethereum.SyncProgress
}

func (p *progress) StartingBlock() uint64 {
	return p.p.StartingBlock
}
func (p *progress) CurrentBlock() uint64 {
	return p.p.CurrentBlock
}
func (p *progress) HighestBlock() uint64 {
	return p.p.HighestBlock
}

func (p *progress) PulledStates() uint64 {
	log.Warn("PulledStates is no longer supported by Geth")
	return 0
}
func (p *progress) KnownStates() uint64 {
	log.Warn("KnownStates is no longer supported by Geth")
	return 0
}
func (p *progress) SyncedAccounts() uint64 {
	return p.p.SyncedAccounts
}
func (p *progress) SyncedAccountBytes() uint64 {
	return p.p.SyncedAccountBytes
}
func (p *progress) SyncedBytecodes() uint64 {
	return p.p.SyncedBytecodes
}
func (p *progress) SyncedBytecodeBytes() uint64 {
	return p.p.SyncedBytecodeBytes
}
func (p *progress) SyncedStorage() uint64 {
	return p.p.SyncedStorage
}
func (p *progress) SyncedStorageBytes() uint64 {
	return p.p.SyncedStorageBytes
}
func (p *progress) HealedTrienodes() uint64 {
	return p.p.HealedTrienodes
}
func (p *progress) HealedTrienodeBytes() uint64 {
	return p.p.HealedTrienodeBytes
}
func (p *progress) HealedBytecodes() uint64 {
	return p.p.HealedBytecodes
}
func (p *progress) HealedBytecodeBytes() uint64 {
	return p.p.HealedBytecodeBytes
}
func (p *progress) HealingTrienodes() uint64 {
	return p.p.HealingTrienodes
}
func (p *progress) HealingBytecode() uint64 {
	return p.p.HealingBytecode
}

func (d *dl) Progress() core.Progress {
	return &progress{d.dl.SyncProgress()}
}

func (b *Backend) Downloader() core.Downloader {
	return &dl{b.b}
}

func (b *Backend) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) core.Subscription {
	var sub event.Subscription
	b.newTxsOnce.Do(func() {
		bch := make(chan gcore.NewTxsEvent, 100)
		sub = b.b.SubscribeNewTxsEvent(bch)
		go func() {
			for {
				select {
				case item := <-bch:
					txe := core.NewTxsEvent{
						Txs: make([][]byte, len(item.Txs)),
					}
					for i, tx := range item.Txs {
						txe.Txs[i], _ = tx.MarshalBinary()
					}
					b.newTxsFeed.Send(txe)
				case err := <-sub.Err():
					log.Warn("Subscription error for NewTxs", "err", err)
					return
				}
			}
		}()
	})
	return b.newTxsFeed.Subscribe(ch)
}
func (b *Backend) SubscribeChainEvent(ch chan<- core.ChainEvent) core.Subscription {
	var sub event.Subscription
	b.chainOnce.Do(func() {
		bch := make(chan gcore.ChainEvent, 100)
		sub = b.b.SubscribeChainEvent(bch)
		go func() {
			for {
				select {
				case item := <-bch:
					ce := core.ChainEvent{
						Hash: core.Hash(item.Hash),
					}
					ce.Block, _ = rlp.EncodeToBytes(item.Block)
					ce.Logs, _ = rlp.EncodeToBytes(item.Logs)
					b.chainFeed.Send(ce)
				case err := <-sub.Err():
					log.Warn("Subscription error for Chain", "err", err)
					return
				}
			}
		}()
	})
	return b.chainFeed.Subscribe(ch)
}
func (b *Backend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) core.Subscription {
	var sub event.Subscription
	b.chainHeadOnce.Do(func() {
		bch := make(chan gcore.ChainHeadEvent, 100)
		sub = b.b.SubscribeChainHeadEvent(bch)
		go func() {
			for {
				select {
				case item := <-bch:
					che := core.ChainHeadEvent{}
					che.Block, _ = rlp.EncodeToBytes(item.Block)
					b.chainHeadFeed.Send(che)
				case err := <-sub.Err():
					log.Warn("Subscription error for ChainHead", "err", err)
					return
				}
			}
		}()
	})
	return b.chainHeadFeed.Subscribe(ch)
}
func (b *Backend) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) core.Subscription {
	var sub event.Subscription
	b.chainSideOnce.Do(func() {
		bch := make(chan gcore.ChainSideEvent, 100)
		sub = b.b.SubscribeChainSideEvent(bch)
		go func() {
			for {
				select {
				case item := <-bch:
					cse := core.ChainSideEvent{}
					cse.Block, _ = rlp.EncodeToBytes(item.Block)
					b.chainSideFeed.Send(cse)
				case err := <-sub.Err():
					log.Warn("Subscription error for ChainSide", "err", err)
					return
				}
			}
		}()
	})
	return b.chainSideFeed.Subscribe(ch)
}
func (b *Backend) SubscribeLogsEvent(ch chan<- [][]byte) core.Subscription {
	var sub event.Subscription
	b.logsOnce.Do(func() {
		bch := make(chan []*types.Log, 100)
		sub = b.b.SubscribeLogsEvent(bch)
		go func() {
			for {
				select {
				case item := <-bch:
					logs := make([][]byte, len(item))
					for i, log := range item {
						logs[i], _ = rlp.EncodeToBytes(log)
					}
					b.logsFeed.Send(logs)
				case err := <-sub.Err():
					log.Warn("Subscription error for Logs", "err", err)
					return
				}
			}
		}()
	})
	return b.logsFeed.Subscribe(ch)
} // []RLP encoded logs
func (b *Backend) SubscribePendingLogsEvent(ch chan<- [][]byte) core.Subscription {
	var sub event.Subscription
	b.pendingLogsOnce.Do(func() {
		bch := make(chan []*types.Log, 100)
		sub = b.b.SubscribePendingLogsEvent(bch)
		go func() {
			for {
				select {
				case item := <-bch:
					logs := make([][]byte, len(item))
					for i, log := range item {
						logs[i], _ = rlp.EncodeToBytes(log)
					}
					b.pendingLogsFeed.Send(logs)
				case err := <-sub.Err():
					log.Warn("Subscription error for PendingLogs", "err", err)
					return
				}
			}
		}()
	})
	return b.pendingLogsFeed.Subscribe(ch)
} // RLP Encoded logs
func (b *Backend) SubscribeRemovedLogsEvent(ch chan<- []byte) core.Subscription {
	var sub event.Subscription
	b.removedLogsOnce.Do(func() {
		bch := make(chan gcore.RemovedLogsEvent, 100)
		sub = b.b.SubscribeRemovedLogsEvent(bch)
		go func() {
			for {
				select {
				case item := <-bch:
					logs := make([][]byte, len(item.Logs))
					for i, log := range item.Logs {
						logs[i], _ = rlp.EncodeToBytes(log)
					}
					b.removedLogsFeed.Send(item)
				case err := <-sub.Err():
					log.Warn("Subscription error for RemovedLogs", "err", err)
					return
				}
			}
		}()
	})
	return b.removedLogsFeed.Subscribe(ch)
} // RLP encoded logs

func convertAndSet(a, b reflect.Value) (err error) {
	defer func() {
		if recover() != nil {
			fmt.Errorf("error converting: %v", err.Error())
		}
	}()
	a.Set(b.Convert(a.Type()))
	return nil
}

func (b *Backend) ChainConfig() *params.ChainConfig {
	// We're using the reflect library to copy data from params.ChainConfig to
	// pparams.ChainConfig, so this function shouldn't need to be touched for
	// simple changes to ChainConfig (though pparams.ChainConfig may need to be
	// updated). Note that this probably won't carry over consensus engine data.
	if b.chainConfig != nil {
		return b.chainConfig
	}
	b.chainConfig = &params.ChainConfig{}
	nval := reflect.ValueOf(b.b.ChainConfig())
	ntype := nval.Elem().Type()
	lval := reflect.ValueOf(b.chainConfig)
	for i := 0; i < nval.Elem().NumField(); i++ {
		field := ntype.Field(i)
		v := nval.Elem().FieldByName(field.Name)
		lv := lval.Elem().FieldByName(field.Name)
		log.Info("Checking value for", "field", field.Name)
		if lv.Kind() != reflect.Invalid {
			// If core.ChainConfig doesn't have this field, skip it.
			if v.Type() == lv.Type() && lv.CanSet() {
				lv.Set(v)
			} else {
				convertAndSet(lv, v)
			}
		}
	}
	return b.chainConfig
}

func (b *Backend) GetTrie(h core.Hash) (core.Trie, error) {
	tr, err := trie.NewStateTrie(trie.TrieID(common.Hash(h)), trie.NewDatabase(b.b.ChainDb()))
	if err != nil {
		return nil, err
	}
	return NewWrappedTrie(tr), nil
}

func (b *Backend) GetAccountTrie(stateRoot core.Hash, account core.Address) (core.Trie, error) {
	tr, err := b.GetTrie(stateRoot)
	if err != nil {
		return nil, err
	}
	act, err := tr.GetAccount(account)
	if err != nil {
		return nil, err
	}
	acTr, err := trie.NewStateTrie(trie.StorageTrieID(common.Hash(stateRoot), crypto.Keccak256Hash(account[:]), common.Hash(act.Root)), trie.NewDatabase(b.b.ChainDb()))
	if err != nil {
		return nil, err
	}
	return NewWrappedTrie(acTr), nil
}





