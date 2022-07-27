package rpc

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/plugins"
)

func PluginGetRPCCalls(pl *plugins.PluginLoader, id, method, params string) {
	fnList := pl.Lookup("GetRPCCalls", func(item interface{}) bool {
		_, ok := item.(func(string, string, string))
		return ok
	})
	for _, fni := range fnList {
		if fn, ok := fni.(func(string, string, string)); ok {
			fn(id, method, params)
		}
	}
}

func pluginGetRPCCalls(id, method, params string) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting GerRPCCalls, but default PluginLoader has not been initialized")
		return
	}
	PluginGetRPCCalls(plugins.DefaultPluginLoader, id, method, params)
}
