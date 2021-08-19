package main

import (
	"context"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/plugins"
	"github.com/ethereum/go-ethereum/plugins/interfaces"
	"github.com/ethereum/go-ethereum/rpc"
	lru "github.com/hashicorp/golang-lru"
	"gopkg.in/urfave/cli.v1"
)

type IsSyncedService struct {
	backend interfaces.Backend
	stack   *node.Node
}

var HTTPApiFlag = cli.StringFlag{
	Name:  "isSynced",
	Usage: "API's offered over the HTTP-RPC interface",
	Value: "",
}

var pl *plugins.PluginLoader
var cache *lru.Cache

//This still doesnt work I assume is has to do with some funtion not being called correctly
//as far as I can tell the name is registered but there is a warning logged saying the method is unavailable
func Initialize(ctx *cli.Context, loader *plugins.PluginLoader) {
	pl = loader
	cache, _ = lru.New(128)
	if ctx.GlobalString(name) == "" {
		ctx.GlobalSet("", HTTPApiFlag.Name)
		log.Info("name is actually set")
	}
	//I added these and other log statements to try and track down the issue
	log.Info("default message")
}
func GetAPIs(stack *node.Node, backend interfaces.Backend) []rpc.API {
	return []rpc.API{
		{
			Namespace: "plugeth",
			Version:   "1.0",
			Service:   &IsSyncedService{backend, stack},
			Public:    true,
		},
	}
}

var peers bool
var synced bool

func (service *IsSyncedService) IsSynced(ctx context.Context) interface{} {
	progress := service.backend.Downloader().Progress()
	peercount := service.stack.Server().PeerCount()
	if peercount > 0 {
		peers = true
	}
	if progress.CurrentBlock >= progress.HighestBlock {
		synced = true
	}
	return map[string]interface{}{
		"startingBlock": hexutil.Uint64(progress.StartingBlock),
		"currentBlock":  hexutil.Uint64(progress.CurrentBlock),
		"highestBlock":  hexutil.Uint64(progress.HighestBlock),
		"pulledStates":  hexutil.Uint64(progress.PulledStates),
		"knownStates":   hexutil.Uint64(progress.KnownStates),
		"activePeers":   peers,
		"nodeIsSynced":  synced,
	}
}
