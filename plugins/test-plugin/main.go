package main

import (
	"context"
	"math/big"
	"time"
	"os"
	
	"github.com/openrelayxyz/plugeth-utils/core"
	"github.com/openrelayxyz/plugeth-utils/restricted/hexutil"
)

var hookChan chan map[string]struct{} = make(chan map[string]struct{}, 10)
var quit chan string = make(chan string)

func (service *engineService) CaptureShutdown(ctx context.Context) {
	m := map[string]struct{}{
		"OnShutdown":struct{}{},
	}
	hookChan <- m
}

func (service *engineService) CapturePreTrieCommit(ctx context.Context) {
	m := map[string]struct{}{
		"PreTrieCommit":struct{}{},
	}
	hookChan <- m
}

func (service *engineService) CapturePostTrieCommit(ctx context.Context) {
	m := map[string]struct{}{
		"PostTrieCommit":struct{}{},
	}
	hookChan <- m
}

func BlockChain() {

	go func () {
		for {
			select {
				case <- quit:
					if len(plugins) > 0 {
						log.Error("Exit with Error, Plugins map not empty", "Plugins not called", plugins)
						os.Exit(1)
					} else {
						log.Info("Exit without error")
						os.Exit(0)
					}
				case m := <- hookChan:
					var ok bool
					f := func(key string) bool {_, ok = m[key]; return ok}
					switch {
						case f("OnShutdown"):
							delete(plugins, "OnShutdown")
						case f("StateUpdate"):
							delete(plugins, "StateUpdate")
						case f("PreProcessBlock"):
							delete(plugins, "PreProcessBlock")
						case f("PreProcessTransaction"):
							delete(plugins, "PreProcessTransaction")
						case f("PostProcessTransaction"):
							delete(plugins, "PostProcessTransaction")
						case f("PostProcessBlock"):
							delete(plugins, "PostProcessBlock")
						case f("NewHead"):
							delete(plugins, "NewHead")
						case f("LivePreProcessBlock"):
							delete(plugins, "LivePreProcessBlock")
						case f("LivePreProcessTransaction"):
							delete(plugins, "LivePreProcessTransaction")
						case f("LivePostProcessTransaction"):
							delete(plugins, "LivePostProcessTransaction")
						case f("LivePostProcessBlock"):
							delete(plugins, "LivePostProcessBlock")
						case f("GetRPCCalls"):
							delete(plugins, "GetRPCCalls")
						case f("RPCSubscriptionTest"):
							delete(plugins, "RPCSubscriptionTest")
						case f("SetTrieFlushIntervalClone"):
							delete(plugins, "SetTrieFlushIntervalClone")
						case f("StandardCaptureStart"):
							delete(plugins, "StandardCaptureStart")
						case f("StandardCaptureState"):
							delete(plugins, "StandardCaptureState")
						case f("StandardCaptureFault"):
							delete(plugins, "StandardCaptureFault")
						case f("StandardCaptureEnter"):
							delete(plugins, "StandardCaptureEnter")
						case f("StandardCaptureExit"):
							delete(plugins, "StandardCaptureExit")
						case f("StandardCaptureEnd"):
							delete(plugins, "StandardCaptureEnd")
						case f("StandardTracerResult"):
							delete(plugins, "StandardTracerResult")
						case f("LivePreProcessBlock"):
							delete(plugins, "LivePreProcessBlock")
						case f("LiveCaptureStart"):
							delete(plugins, "LiveCaptureStart")
						case f("LiveCaptureState"):
							delete(plugins, "LiveCaptureState")
						// These methods are not covered by tests at this time
						// case f("LiveCaptureFault"):
						// 	delete(plugins, "LiveCaptureFault")
						// case f("LiveCaptureEnter"):
						// 	delete(plugins, "LiveCaptureEnter")
						// case f("LiveCaptureExit"):
						// 	delete(plugins, "LiveCaptureExit")
						// case f("LiveTracerResult"):
						// 	delete(plugins, "LiveTracerResult")
						case f("LiveCaptureEnd"):
							delete(plugins, "LiveCaptureEnd")
						case f("PreTrieCommit"):
							delete(plugins, "PreTrieCommit")
						case f("PostTrieCommit"):
							delete(plugins, "PostTrieCommit")
				}
			}
		}
	}()
	
	txFactory()
	txTracer()
}

var t0 core.Hash
var t1 core.Hash
var t2 core.Hash
var t3 core.Hash
var coinBase *core.Address

func txFactory() {

	cl := apis[0].Service.(*engineService).stack
	client, err := cl.Attach()
	if err != nil {
		log.Error("Error connecting with client txFactory", "err", err)
	}

	err = client.Call(&coinBase, "eth_coinbase")
	if err != nil {
		log.Error("failed to call eth_coinbase txFactory", "err", err)
	}

	var peerCount hexutil.Uint64
	for peerCount == 0 {
		err = client.Call(&peerCount, "net_peerCount")
		if err != nil {
			log.Error("failed to call net_peerCount", "err", err)
		}
		time.Sleep(100 * time.Millisecond)
	} 

	tx0_params := map[string]interface{}{
		"from": coinBase,
		"to": coinBase,
		"value": (*hexutil.Big)(big.NewInt(1)),
	}
	
	err = client.Call(&t0, "eth_sendTransaction", tx0_params)
	if err != nil {
		log.Error("transaction zero failed", "err", err)
	}

	tx1_params := map[string]interface{}{
		"input": "0x60018080600053f3",
		"from": coinBase,
	}

	time.Sleep(2 * time.Second)
	err = client.Call(&t1, "eth_sendTransaction", tx1_params)
	if err != nil {
		log.Error("transaction one failed", "err", err)
	}

	tx2_params := map[string]interface{}{
		"input": "0x61520873000000000000000000000000000000000000000060006000600060006000f1",
		"from": coinBase,
	}
	
	time.Sleep(2 * time.Second)
	err = client.Call(&t2, "eth_sendTransaction", tx2_params)
	if err != nil {
		log.Error("transaction two failed", "err", err)
	}
	
	genericArg := map[string]interface{}{
		"input": "0x608060405234801561001057600080fd5b5061011a806100206000396000f3fe608060405234801561001057600080fd5b50600436106100375760003560e01c806360fe47b11461003c5780636d4ce63c1461005d57610037565b600080fd5b61004561007e565b60405161005291906100c5565b60405180910390f35b61007c6004803603602081101561007a57600080fd5b50356100c2565b6040516020018083838082843780820191505050505b565b005b6100946100c4565b60405161005291906100bf565b6100d1565b60405180910390f35b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb60e11b815260040161010060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16146101e557600080fd5b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663e7ba30df6040518163ffffffff1660e01b8152600401600060405180830381600087803b1580156101ae57600080fd5b505af11580156101c2573d6000803e3d6000fd5b50505050505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161461029157600080fd5b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663fdacd5766040518163ffffffff1660e01b8152600401600060405180830381600087803b1580156102f957600080fd5b505af115801561030d573d6000803e3d6000fd5b50505050505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff168156fea2646970667358221220d4f2763f3a0ae2826cc9ef37a65ff0c14d7a3aafe8d1636ff99f72e2f705413d64736f6c634300060c0033",
		"from": coinBase,
	}

	for i := 0; i < 126; i ++ {
		time.Sleep(2 * time.Second)
		err = client.Call(&t3, "eth_sendTransaction", genericArg)
		if err != nil {
			log.Error("looped transaction failed on index", "i", i, "err", err)
		}
	}

}

type TraceConfig struct {
	Tracer  *string
}

func txTracer() {

	cl := apis[0].Service.(*engineService).stack
	client, err := cl.Attach()
	if err != nil {
		log.Error("Error connecting with client block factory")
	}

	time.Sleep(2 * time.Second)
	tr := "testTracer"
	t := TraceConfig{
		Tracer: &tr,
	}

	var trResult interface{}
	err = client.Call(&trResult, "debug_traceTransaction", t0, t)
	if err != nil {
		log.Error("debug_traceTransaction failed",  "err", err)
	}

	debugArg0 := map[string]interface{}{
		"input": "0x60006000fd",
		"from": coinBase,
	}

	var trResult0 interface{}
	err = client.Call(&trResult0, "debug_traceCall", debugArg0, "latest", t)
	if err != nil {
		log.Error("debug_traceCall 0 failed",  "err", err)
	}

	debugArg1 := map[string]interface{}{
		"input": "0x61520873000000000000000000000000000000000000000060006000600060006000f1",
		"from": coinBase,
	}

	var trResult1 interface{}
	err = client.Call(&trResult1, "debug_traceCall", debugArg1, "latest", t)


	final := map[string]interface{}{
		"input": "0x61520873000000000000000000000000000000000000000060006000600060006000f1",
		"from": coinBase,
	}

	time.Sleep(2 * time.Second)
	err = client.Call(&t3, "eth_sendTransaction", final)
	if err != nil {
		log.Error("contract call failed", "err", err)
	}

	quit <- "quit"

}

