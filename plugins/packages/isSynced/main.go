package main

import (
	"context"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/plugins"
	"github.com/ethereum/go-ethereum/plugins/interfaces"
	"github.com/ethereum/go-ethereum/rpc"
	"gopkg.in/urfave/cli.v1"
)

type IsSyncedService struct {
	backend interfaces.Backend
	stack   *node.Node
}

var HTTPApiFlag = cli.StringFlag{
	Name:  "http.api",
	Usage: "API's offered over the HTTP-RPC interface",
	Value: "",
}

func Initialize(ctx *cli.Context, loader *plugins.PluginLoader) {
	v := ctx.GlobalString(utils.HTTPApiFlag.Name)
	if v != "" {
		ctx.GlobalSet(HTTPApiFlag.Name, v+"plugeth")
	} else {
		ctx.GlobalSet(HTTPApiFlag.Name, "eth,net,web3,plugeth")
		log.Info("Loaded isSynced plugin")
	}
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

func (service *IsSyncedService) IsSynced(ctx context.Context) interface{} {
	progress := service.backend.Downloader().Progress()
	peercount := service.stack.Server().PeerCount()
	return map[string]interface{}{
		"startingBlock": hexutil.Uint64(progress.StartingBlock),
		"currentBlock":  hexutil.Uint64(progress.CurrentBlock),
		"highestBlock":  hexutil.Uint64(progress.HighestBlock),
		"pulledStates":  hexutil.Uint64(progress.PulledStates),
		"knownStates":   hexutil.Uint64(progress.KnownStates),
		"activePeers":   peercount > 0,
		"nodeIsSynced":  progress.CurrentBlock >= progress.HighestBlock,
	}
}
