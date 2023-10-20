package vm

import (
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/plugins"
)

func (st *Stack) Len() int {
	return len(st.data)
}

func PluginOpCodeSelect(pl *plugins.PluginLoader, jt *JumpTable) *JumpTable {
	var opCodes []int
	fnList := pl.Lookup("OpCodeSelect", func(item interface{}) bool {
		_, ok := item.(func() []int)
		return ok
	})
	for _, fni := range fnList {
		if fn, ok := fni.(func() []int); ok {
			opCodes = append(opCodes, fn()...)
		}
	}
	if len(opCodes) > 0 {
		jt = copyJumpTable(jt)
	}
	for _, idx := range opCodes {
		(*jt)[idx] = nil
	}
	return jt
}

func pluginOpCodeSelect(jt *JumpTable) *JumpTable {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting PluginOpCodeSelect, but default PluginLoader has not been initialized")
		return nil
	}
	return PluginOpCodeSelect(plugins.DefaultPluginLoader, jt)
}

