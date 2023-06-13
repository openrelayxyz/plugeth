package core

import (
	// "reflect"
	// "sync/atomic"
	"hash"
	"fmt"
	"testing"
	"math/big"
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"golang.org/x/crypto/sha3"
)

var (
	config = &params.ChainConfig{
		ChainID:             big.NewInt(1),
		HomesteadBlock:      big.NewInt(0),
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		IstanbulBlock:       big.NewInt(0),
		MuirGlacierBlock:    big.NewInt(0),
		BerlinBlock:         big.NewInt(0),
		LondonBlock:         big.NewInt(0),
		Ethash:              new(params.EthashConfig),
	}
	signer  = types.LatestSigner(config)
	key1, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	key2, _ = crypto.HexToECDSA("0202020202020202020202020202020202020202020202020202002020202020")
)

var makeTx = func(key *ecdsa.PrivateKey, nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *types.Transaction {
	tx, _ := types.SignTx(types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, data), signer, key)
	return tx
}

var (
	db    = rawdb.NewMemoryDatabase()
	gspec = &Genesis{
		Config: config,
		Alloc: GenesisAlloc{
			common.HexToAddress("0x71562b71999873DB5b286dF957af199Ec94617F7"): GenesisAccount{
				Balance: big.NewInt(1000000000000000000), // 1 ether
				Nonce:   0,
			},
			common.HexToAddress("0xfd0810DD14796680f72adf1a371963d0745BCc64"): GenesisAccount{
				Balance: big.NewInt(1000000000000000000), // 1 ether
				Nonce:   math.MaxUint64,
			},
		},
	}
)

type testHasher struct {
	hasher hash.Hash
}

func newHasher() *testHasher {
	return &testHasher{hasher: sha3.NewLegacyKeccak256()}
}

func (h *testHasher) Reset() {
	h.hasher.Reset()
}

func (h *testHasher) Update(key, val []byte) {
	h.hasher.Write(key)
	h.hasher.Write(val)
}

func (h *testHasher) Hash() common.Hash {
	return common.BytesToHash(h.hasher.Sum(nil))
}


func TestPlugethInjections(t *testing.T) {

	
	blockchain, _ := NewBlockChain(db, nil, gspec, nil, ethash.NewFaker(), vm.Config{}, nil, nil)

	engine := ethash.NewFaker()
	
	sp := NewStateProcessor(config, blockchain, engine) 

	txns := []*types.Transaction{
		makeTx(key1, 0, common.Address{}, big.NewInt(1000), params.TxGas-1000, big.NewInt(875000000), nil),
	}

	block := GenerateBadBlock(gspec.ToBlock(), engine, txns, gspec.Config)

	statedb, _ := state.New(blockchain.GetBlockByHash(block.ParentHash()).Root(), blockchain.stateCache, nil)

	t.Run(fmt.Sprintf("test BlockProcessingError"), func(t *testing.T) {
		called := false
		injectionCalled = &called

		_, _, _, _ = sp.Process(block, statedb, vm.Config{})
		
		if *injectionCalled != true {
			t.Fatalf("pluginBlockProcessingError injection in stateProcessor.Process() not called")
		}
	})

	t.Run(fmt.Sprintf("test Reorg"), func(t *testing.T) {
		called := false
		injectionCalled = &called

		// the transaction has to be initialized with a different gas price than the previous tx in order to trigger a reorg
		txns2 := []*types.Transaction{
			makeTx(key1, 0, common.Address{}, big.NewInt(1000), params.TxGas-1000, big.NewInt(875000001), nil),
		}
		block2 := GenerateBadBlock(gspec.ToBlock(), engine, txns2, gspec.Config)


		_, _ = blockchain.writeBlockAndSetHead(block, []*types.Receipt{}, []*types.Log{}, statedb, false)

		_ = blockchain.reorg(block.Header(), block2)
		
		if *injectionCalled != true {
			t.Fatalf("pluginReorg injection in blockChain.Reorg() not called")
		}
	})

	t.Run(fmt.Sprintf("test treiIntervarFlushClone"), func(t *testing.T) {
		called := false
		injectionCalled = &called

		_ = blockchain.writeBlockWithState(block, []*types.Receipt{}, statedb)

		if *injectionCalled != true {
			t.Fatalf("pluginNewSideBlock injection in blockChain.writeBlockAndSetHead() not called")
		}
	})

	t.Run(fmt.Sprintf("test NewSideBlock"), func(t *testing.T) {
		called := false
		injectionCalled = &called

		TestReorgToShorterRemovesCanonMapping(t)

		if *injectionCalled != true {
			t.Fatalf("pluginNewSideBlock injection in blockChain.writeBlockAndSetHead() not called")
		}
	})
	
}	