package main

import (
	"context"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/plugins"
	"github.com/ethereum/go-ethereum/plugins/interfaces"
	"github.com/ethereum/go-ethereum/rpc"
	lru "github.com/hashicorp/golang-lru"
	"gopkg.in/urfave/cli.v1"
)

type MyService struct {
	backend interfaces.Backend
	stack   *node.Node
}

var pl *plugins.PluginLoader
var cache *lru.Cache

//The initialize method needs to be modified in order to allow geth to start without appending flags

func Initialize(ctx *cli.Context, loader *plugins.PluginLoader) {
	pl = loader

	cache, _ = lru.New(128) // TODO: Make size configurable
	if !ctx.GlobalBool(utils.SnapshotFlag.Name) {
		log.Warn("Snapshots are required for StateUpdate plugins, but are currently disabled. State Updates will be unavailable")
	}
	log.Info("loaded isSynced plugin")
}

func GetAPIs(stack *node.Node, backend interfaces.Backend) []rpc.API {
	return []rpc.API{
		{
			Namespace: "mynamespace",
			Version:   "1.0",
			Service:   &MyService{backend, stack},
			Public:    true,
		},
	}
}

func (myserv *MyService) IsSynced(ctx context.Context) bool {
	dwnlder := myserv.backend.Downloader()
	return myserv.stack.Server().PeerCount() > 0 && dwnlder.Progress().CurrentBlock >= dwnlder.Progress().HighestBlock
}
