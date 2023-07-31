package main

import (
	"time"
	"math/big"
	"sync"
	
	"github.com/openrelayxyz/plugeth-utils/core"
	
)

var apis []core.API

type engineService struct {
	backend core.Backend
	stack core.Node
}

// cmd/geth/

func GetAPIs(stack core.Node, backend core.Backend) []core.API {
	// GetAPIs is covered by virtue of the plugeth_captureShutdown method functioning.
	apis = []core.API{
		{
			Namespace: "plugeth",
			Version:   "1.0",
			Service:   &engineService{backend, stack},
			Public:    true,
		},
		{
			Namespace: "plugeth",
			Version:   "1.0",
			Service:   &LiveTracerResult{},
			Public:    true,
		},
	}
	return apis
}

// func OnShutdown(){
	// this injection is covered by another test in this package. See documentation for details. 
// }

// core/


func PreProcessBlock(hash core.Hash, number uint64, encoded []byte) {
	m := map[string]struct{}{
		"PreProcessBlock":struct{}{},
	}
	hookChan <- m
}

func PreProcessTransaction(txBytes []byte, txHash, blockHash core.Hash, i int) {
	m := map[string]struct{}{
		"PreProcessTransaction":struct{}{},
	}
	hookChan <- m
}

func BlockProcessingError(tx core.Hash, block core.Hash, err error) { 
	// this injection is covered by a stand alone test: plugeth_injection_test.go in the core/ package. 
}

func PostProcessTransaction(tx core.Hash, block core.Hash, i int, receipt []byte) {
	m := map[string]struct{}{
		"PostProcessTransaction":struct{}{},
	}
	hookChan <- m
}

func PostProcessBlock(block core.Hash) {
	m := map[string]struct{}{
		"PostProcessBlock":struct{}{},
	}
	hookChan <- m
}

func NewHead(block []byte, hash core.Hash, logs [][]byte, td *big.Int) {
	m := map[string]struct{}{
		"NewHead":struct{}{},
	}
	hookChan <- m
}

func NewSideBlock(block []byte, hash core.Hash, logs [][]byte) { // beyond the scope of the test at this time
	// this injection is covered by a stand alone test: plugeth_injection_test.go in the core/ package.
}

func Reorg(commonBlock core.Hash, oldChain, newChain []core.Hash) { // beyond the scope of the test at this time
	// this injection is covered by a stand alone test: plugeth_injection_test.go in the core/ package.
}

func SetTrieFlushIntervalClone(duration time.Duration) time.Duration {
	m := map[string]struct{}{
		"SetTrieFlushIntervalClone":struct{}{},
	}
	hookChan <- m
	return duration
}

// core/rawdb/

func ModifyAncients(index uint64, freezerUpdate map[string]struct{}) {
	// this injection is covered by a stand alone test: plugeth_injection_test.go in the core/rawdb package. 
}

func AppendAncient(number uint64, hash, header, body, receipts, td []byte) {
	// this injection is covered by a stand alone test: plugeth_injection_test.go in the core/rawdb package.
}

// core/state/

func StateUpdate(blockRoot core.Hash, parentRoot core.Hash, coreDestructs map[core.Hash]struct{}, coreAccounts map[core.Hash][]byte, coreStorage map[core.Hash]map[core.Hash][]byte, coreCode map[core.Hash][]byte) {
	// log.Warn("StatueUpdate", "blockRoot", blockRoot, "parentRoot", parentRoot, "coreDestructs", coreDestructs, "coreAccounts", coreAccounts, "coreStorage", coreStorage, "coreCode", coreCode)
	m := map[string]struct{}{
		"StateUpdate":struct{}{},
	}
	hookChan <- m
}

// rpc/


func GetRPCCalls(method string, id string, params string) {
	m := map[string]struct{}{
		"GetRPCCalls":struct{}{},
	}
	hookChan <- m
}

var once sync.Once

func RPCSubscriptionTest() {
	go func() {
		once.Do(func() {
			m := map[string]struct{}{
			"RPCSubscriptionTest":struct{}{},
			}
			hookChan <- m
		})
	}()
}

// trie/

// func PreTrieCommit(node core.Hash) {
	// this injection is covered by another test in this package. See documentation for details.
// }

// func PostTrieCommit(node core.Hash) {
	// this injection is covered by another test in this package. See documentation for details.
// }

var plugins map[string]struct{} = map[string]struct{}{
	"OnShutdown": struct{}{},
	"SetTrieFlushIntervalClone":struct{}{},
	"StateUpdate": struct{}{},
	"PreProcessBlock": struct{}{},
	"PreProcessTransaction": struct{}{},
	"PostProcessTransaction": struct{}{},
	"PostProcessBlock": struct{}{},
	"NewHead": struct{}{},
	"StandardCaptureStart": struct{}{},
	"StandardCaptureState": struct{}{},
	"StandardCaptureFault": struct{}{},
	"StandardCaptureEnter": struct{}{},
	"StandardCaptureExit": struct{}{},
	"StandardCaptureEnd": struct{}{},
	"StandardTracerResult": struct{}{},
	"GetRPCCalls": struct{}{},
	"RPCSubscriptionTest": struct{}{},
	"LivePreProcessBlock": struct{}{},
	"LivePreProcessTransaction": struct{}{},
	"LivePostProcessTransaction": struct{}{},
	"LivePostProcessBlock": struct{}{},
	"LiveCaptureStart": struct{}{},
	"LiveCaptureState": struct{}{},
	// "LiveCaptureFault": struct{}{},
	// "LiveCaptureEnter": struct{}{},
	// "LiveCaptureExit": struct{}{},
	// "LiveTracerResult": struct{}{},
	"LiveCaptureEnd": struct{}{},
	"PreTrieCommit": struct{}{},
	"PostTrieCommit": struct{}{},
} 

