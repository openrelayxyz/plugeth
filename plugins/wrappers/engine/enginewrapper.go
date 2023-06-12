package engine

import (
	"fmt"
	"math/big"
	"reflect"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/plugins/wrappers"
	"github.com/openrelayxyz/plugeth-utils/core"
	ptypes "github.com/openrelayxyz/plugeth-utils/restricted/types"
	pconsensus "github.com/openrelayxyz/plugeth-utils/restricted/consensus"
	pparams "github.com/openrelayxyz/plugeth-utils/restricted/params"
)


func gethToUtilsHeader(header *types.Header) *ptypes.Header {
	if header == nil { return nil }
	return &ptypes.Header{
		ParentHash: core.Hash(header.ParentHash),
		UncleHash: core.Hash(header.UncleHash),
		Coinbase: core.Address(header.Coinbase),
		Root: core.Hash(header.Root),
		TxHash: core.Hash(header.TxHash),
		ReceiptHash: core.Hash(header.ReceiptHash),
		Bloom: ptypes.Bloom(header.Bloom),
		Difficulty: header.Difficulty,
		Number: header.Number,
		GasLimit: header.GasLimit,
		GasUsed: header.GasUsed,
		Time: header.Time,
		Extra: header.Extra,
		MixDigest: core.Hash(header.MixDigest),
		Nonce: ptypes.BlockNonce(header.Nonce),
		BaseFee: header.BaseFee,
		WithdrawalsHash: (*core.Hash)(header.WithdrawalsHash),
	}
}
func utilsToGethHeader(header *ptypes.Header) *types.Header {
	if header == nil { return nil }
	return &types.Header{
		ParentHash: common.Hash(header.ParentHash),
		UncleHash: common.Hash(header.UncleHash),
		Coinbase: common.Address(header.Coinbase),
		Root: common.Hash(header.Root),
		TxHash: common.Hash(header.TxHash),
		ReceiptHash: common.Hash(header.ReceiptHash),
		Bloom: types.Bloom(header.Bloom),
		Difficulty: header.Difficulty,
		Number: header.Number,
		GasLimit: header.GasLimit,
		GasUsed: header.GasUsed,
		Time: header.Time,
		Extra: header.Extra,
		MixDigest: common.Hash(header.MixDigest),
		Nonce: types.BlockNonce(header.Nonce),
		BaseFee: header.BaseFee,
		WithdrawalsHash: (*common.Hash)(header.WithdrawalsHash),
	}
}

func gethToUtilsTransactions(transactions []*types.Transaction) []*ptypes.Transaction {
	if transactions == nil { return nil }
	txs := make([]*ptypes.Transaction, len(transactions))
	for i, tx := range transactions {
		bin, err := tx.MarshalBinary()
		if err != nil { panic (err) }
		txs[i] = &ptypes.Transaction{}
		txs[i].UnmarshalBinary(bin)
	}
	return txs
}

func gethToUtilsHeaders(headers []*types.Header) []*ptypes.Header {
	if headers == nil { return nil }
	pheaders := make([]*ptypes.Header, len(headers))
	for i, header := range headers {
		pheaders[i] = gethToUtilsHeader(header)
	}
	return pheaders
}

func gethToUtilsReceipts(receipts []*types.Receipt) []*ptypes.Receipt {
	if receipts == nil { return nil }
	preceipts := make([]*ptypes.Receipt, len(receipts))
	for i, receipt := range receipts {
		preceipts[i] = &ptypes.Receipt{
			Type: receipt.Type,
			PostState: receipt.PostState,
			Status: receipt.Status,
			CumulativeGasUsed: receipt.CumulativeGasUsed,
			Bloom: ptypes.Bloom(receipt.Bloom),
			Logs: gethToUtilsLogs(receipt.Logs),
			TxHash: core.Hash(receipt.TxHash),
			ContractAddress: core.Address(receipt.ContractAddress),
			GasUsed: receipt.GasUsed,
			BlockHash: core.Hash(receipt.BlockHash),
			BlockNumber: receipt.BlockNumber,
			TransactionIndex: receipt.TransactionIndex,
		}
	}
	return preceipts
}

func gethToUtilsWithdrawals(withdrawals []*types.Withdrawal) []*ptypes.Withdrawal {
	if withdrawals == nil { return nil }
	pwithdrawals := make([]*ptypes.Withdrawal, len(withdrawals))
	for i, withdrawal := range withdrawals {
		pwithdrawals[i] = &ptypes.Withdrawal{
			Index: withdrawal.Index,
			Validator: withdrawal.Validator,
			Address: core.Address(withdrawal.Address),
			Amount: withdrawal.Amount,
		}
	}
	return pwithdrawals
}

func gethToUtilsBlock(block *types.Block) *ptypes.Block {
	if block == nil { return nil }
	return ptypes.NewBlockWithHeader(gethToUtilsHeader(block.Header())).WithBody(gethToUtilsTransactions(block.Transactions()), gethToUtilsHeaders(block.Uncles())).WithWithdrawals(gethToUtilsWithdrawals(block.Withdrawals()))
}

func utilsToGethBlock(block *ptypes.Block) *types.Block {
	if block == nil { return nil }
	return types.NewBlockWithHeader(utilsToGethHeader(block.Header())).WithBody(utilsToGethTransactions(block.Transactions()), utilsToGethHeaders(block.Uncles())).WithWithdrawals(utilsToGethWithdrawals(block.Withdrawals()))
}

func utilsToGethHeaders(headers []*ptypes.Header) []*types.Header {
	if headers == nil { return nil }
	pheaders := make([]*types.Header, len(headers))
	for i, header := range headers {
		pheaders[i] = utilsToGethHeader(header)
	}
	return pheaders
}

func utilsToGethTransactions(transactions []*ptypes.Transaction) []*types.Transaction {
	if transactions == nil { return nil }
	txs := make([]*types.Transaction, len(transactions))
	for i, tx := range transactions {
		bin, err := tx.MarshalBinary()
		if err != nil { panic (err) }
		txs[i] = &types.Transaction{}
		txs[i].UnmarshalBinary(bin)
	}
	return txs
}

func utilsToGethWithdrawals(withdrawals []*ptypes.Withdrawal) []*types.Withdrawal {
	if withdrawals == nil { return nil }
	pwithdrawals := make([]*types.Withdrawal, len(withdrawals))
	for i, withdrawal := range withdrawals {
		pwithdrawals[i] = &types.Withdrawal{
			Index: withdrawal.Index,
			Validator: withdrawal.Validator,
			Address: common.Address(withdrawal.Address),
			Amount: withdrawal.Amount,
		}
	}
	return pwithdrawals
}

func gethToUtilsBlockChan(ch chan<- *types.Block) chan<- *ptypes.Block {
	pchan := make(chan *ptypes.Block)
	go func() {
		for block := range pchan {
			ch <- utilsToGethBlock(block)
		}
	}()
	return pchan
}

func gethToUtilsLog(logRecord *types.Log) *ptypes.Log {
	if logRecord == nil { return nil }
	topics := make([]core.Hash, len(logRecord.Topics))
	for i, t := range logRecord.Topics {
		topics[i] = core.Hash(t)
	}
	return &ptypes.Log{
		Address: core.Address(logRecord.Address),
		Topics: topics,
		Data: logRecord.Data,
		BlockNumber: logRecord.BlockNumber,
		TxHash: core.Hash(logRecord.TxHash),
		TxIndex: logRecord.TxIndex,
		BlockHash: core.Hash(logRecord.BlockHash),
		Index: logRecord.Index,
		Removed: logRecord.Removed,
	}
}

func utilsToGethLog(logRecord *ptypes.Log) *types.Log {
	if logRecord == nil { return nil }
	topics := make([]common.Hash, len(logRecord.Topics))
	for i, t := range logRecord.Topics {
		topics[i] = common.Hash(t)
	}
	return &types.Log{
		Address: common.Address(logRecord.Address),
		Topics: topics,
		Data: logRecord.Data,
		BlockNumber: logRecord.BlockNumber,
		TxHash: common.Hash(logRecord.TxHash),
		TxIndex: logRecord.TxIndex,
		BlockHash: common.Hash(logRecord.BlockHash),
		Index: logRecord.Index,
		Removed: logRecord.Removed,
	}
}

func gethToUtilsLogs(logs []*types.Log) []*ptypes.Log {
	result := make([]*ptypes.Log, len(logs))
	for i, logRecord := range logs {
		result[i] = gethToUtilsLog(logRecord)
	}
	return result
}

func utilsToGethLogs(logs []*ptypes.Log) []*types.Log {
	result := make([]*types.Log, len(logs))
	for i, logRecord := range logs {
		result[i] = utilsToGethLog(logRecord)
	}
	return result
}


func convertAndSet(a, b reflect.Value) (err error) {
	defer func() {
		if recover() != nil {
			fmt.Errorf("error converting: %v", err.Error())
		}
	}()
	a.Set(b.Convert(a.Type()))
	return nil
}

func gethToUtilsConfig(gcfg *params.ChainConfig) *pparams.ChainConfig {
	cfg := &pparams.ChainConfig{}
	nval := reflect.ValueOf(gcfg)
	ntype := nval.Elem().Type()
	lval := reflect.ValueOf(cfg)
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
	return cfg
}

type WrappedHeaderReader struct {
	chr consensus.ChainHeaderReader
	cfg *pparams.ChainConfig
}

func (whr *WrappedHeaderReader) Config() *pparams.ChainConfig {
	if whr.cfg == nil {
		whr.cfg = gethToUtilsConfig(whr.chr.Config())
	}
	return whr.cfg
}

// CurrentHeader retrieves the current header from the local chain.
func (whr *WrappedHeaderReader) CurrentHeader() *ptypes.Header {
	return gethToUtilsHeader(whr.chr.CurrentHeader())
}

// GetHeader retrieves a block header from the database by hash and number.
func (whr *WrappedHeaderReader) GetHeader(hash core.Hash, number uint64) *ptypes.Header {
	return gethToUtilsHeader(whr.chr.GetHeader(common.Hash(hash), number))
}

// GetHeaderByNumber retrieves a block header from the database by number.
func (whr *WrappedHeaderReader) GetHeaderByNumber(number uint64) *ptypes.Header {
	return gethToUtilsHeader(whr.chr.GetHeaderByNumber(number))
}

// GetHeaderByHash retrieves a block header from the database by its hash.
func (whr *WrappedHeaderReader) GetHeaderByHash(hash core.Hash) *ptypes.Header {
	return gethToUtilsHeader(whr.chr.GetHeaderByHash(common.Hash(hash)))
}

// GetTd retrieves the total difficulty from the database by hash and number.
func (whr *WrappedHeaderReader) GetTd(hash core.Hash, number uint64) *big.Int {
	return whr.chr.GetTd(common.Hash(hash), number)
}


type WrappedChainReader struct {
	chr consensus.ChainReader
	cfg *pparams.ChainConfig
}

func (whr *WrappedChainReader) Config() *pparams.ChainConfig {
	// We're using the reflect library to copy data from params.ChainConfig to
	// pparams.ChainConfig, so this function shouldn't need to be touched for
	// simple changes to ChainConfig (though pparams.ChainConfig may need to be
	// updated). Note that this probably won't carry over consensus engine data.
	if whr.cfg == nil {
		whr.cfg = gethToUtilsConfig(whr.chr.Config())
	}
	return whr.cfg
}

// CurrentHeader retrieves the current header from the local chain.
func (whr *WrappedChainReader) CurrentHeader() *ptypes.Header {
	return gethToUtilsHeader(whr.chr.CurrentHeader())
}

// GetHeader retrieves a block header from the database by hash and number.
func (whr *WrappedChainReader) GetHeader(hash core.Hash, number uint64) *ptypes.Header {
	return gethToUtilsHeader(whr.chr.GetHeader(common.Hash(hash), number))
}

// GetHeaderByNumber retrieves a block header from the database by number.
func (whr *WrappedChainReader) GetHeaderByNumber(number uint64) *ptypes.Header {
	return gethToUtilsHeader(whr.chr.GetHeaderByNumber(number))
}

// GetHeaderByHash retrieves a block header from the database by its hash.
func (whr *WrappedChainReader) GetHeaderByHash(hash core.Hash) *ptypes.Header {
	return gethToUtilsHeader(whr.chr.GetHeaderByHash(common.Hash(hash)))
}

// GetTd retrieves the total difficulty from the database by hash and number.
func (whr *WrappedChainReader) GetTd(hash core.Hash, number uint64) *big.Int {
	return whr.chr.GetTd(common.Hash(hash), number)
}

func (whr *WrappedChainReader) GetBlock(hash core.Hash, number uint64) *ptypes.Block {
	return gethToUtilsBlock(whr.chr.GetBlock(common.Hash(hash), number))
}

type hasherWrapper struct {
	th types.TrieHasher
}
func (hw *hasherWrapper) Reset() {
	hw.th.Reset()
}
func (hw *hasherWrapper) Update(a, b []byte) {
	hw.th.Update(a, b)
}
func (hw *hasherWrapper) Hash() core.Hash {
	return core.Hash(hw.th.Hash())
}

type engineWrapper struct {
	engine pconsensus.Engine
}

func NewWrappedEngine(e pconsensus.Engine) consensus.Engine {
	return &engineWrapper {
		engine: e,
	}
}

func (ew *engineWrapper) Author(header *types.Header) (common.Address, error) {
	addr, err := ew.engine.Author(gethToUtilsHeader(header))
	return common.Address(addr), err
}
func (ew *engineWrapper) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, seal bool) error {
	return ew.engine.VerifyHeader(&WrappedHeaderReader{chain, nil}, gethToUtilsHeader(header), seal)
}
func (ew *engineWrapper) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	pheaders := make([]*ptypes.Header, len(headers))
	for i, header := range headers {
		pheaders[i] = gethToUtilsHeader(header)
	}
	return ew.engine.VerifyHeaders(&WrappedHeaderReader{chain, nil}, pheaders, seals)
}
func (ew *engineWrapper) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	return ew.engine.VerifyUncles(&WrappedChainReader{chain, nil}, gethToUtilsBlock(block))
}
func (ew *engineWrapper) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	uHeader := gethToUtilsHeader(header)
	if err := ew.engine.Prepare(&WrappedHeaderReader{chain, nil}, uHeader); err != nil {
		return err
	}
	*header = *utilsToGethHeader(uHeader)
	return nil
}
func (ew *engineWrapper) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, withdrawals []*types.Withdrawal) {
	ew.engine.Finalize(&WrappedHeaderReader{chain, nil}, gethToUtilsHeader(header), wrappers.NewWrappedStateDB(state), gethToUtilsTransactions(txs), gethToUtilsHeaders(uncles), gethToUtilsWithdrawals(withdrawals))
}
func (ew *engineWrapper) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt, withdrawals []*types.Withdrawal) (*types.Block, error) {
	block, err := ew.engine.FinalizeAndAssemble(&WrappedHeaderReader{chain, nil}, gethToUtilsHeader(header), wrappers.NewWrappedStateDB(state), gethToUtilsTransactions(txs), gethToUtilsHeaders(uncles), gethToUtilsReceipts(receipts), gethToUtilsWithdrawals(withdrawals))
	return utilsToGethBlock(block), err
}
func (ew *engineWrapper) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	return ew.engine.Seal(&WrappedHeaderReader{chain, nil}, gethToUtilsBlock(block), gethToUtilsBlockChan(results), stop)
}
func (ew *engineWrapper) SealHash(header *types.Header) common.Hash {
	return common.Hash(ew.engine.SealHash(gethToUtilsHeader(header)))
}
func (ew *engineWrapper) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	return ew.engine.CalcDifficulty(&WrappedHeaderReader{chain, nil}, time, gethToUtilsHeader(parent))
}
func (ew *engineWrapper) APIs(chain consensus.ChainHeaderReader) []rpc.API {
	papis := ew.engine.APIs(&WrappedHeaderReader{chain, nil})
	apis := make([]rpc.API, len(papis))
	for i, api := range papis {
		apis[i] = rpc.API{
			Namespace: api.Namespace,
			Version: api.Version,
			Service: api.Service,
			Public: api.Public,
		}
	}
	return apis
}
func (ew *engineWrapper) Close() error {
	return ew.engine.Close()
}
