package rawdb


import (
	"github.com/ethereum/go-ethereum/plugins"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	freezerUpdates map[uint64]map[string]interface{}
)

func PluginTrackUpdate(num uint64, kind string, value interface{}) {
	if freezerUpdates == nil { freezerUpdates = make(map[uint64]map[string]interface{}) }
	update, ok := freezerUpdates[num]
	if !ok {
		update = make(map[string]interface{})
		freezerUpdates[num] = update
	}
	update[kind] = value
}

func PluginResetUpdate(num uint64) {
	delete(freezerUpdates, num)
}


func pluginCommitUpdate(num uint64) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting CommitUpdate, but default PluginLoader has not been initialized")
		return
	}
	PluginCommitUpdate(plugins.DefaultPluginLoader, num)
}

func PluginCommitUpdate(pl *plugins.PluginLoader, num uint64) {
	if freezerUpdates == nil { freezerUpdates = make(map[uint64]map[string]interface{}) }
	defer func() { delete(freezerUpdates, num) }()
	update, ok := freezerUpdates[num]
	if !ok {
		log.Warn("Attempting to commit untracked block", "num", num)
		return
	}
	fnList := pl.Lookup("ModifyAncients", func(item interface{}) bool {
		_, ok := item.(func(uint64, map[string]interface{}))
		return ok
	})
	for _, fni := range fnList {
		if fn, ok := fni.(func(uint64, map[string]interface{})); ok {
			fn(num, update)
		}
	}
	appendAncientFnList := pl.Lookup("AppendAncient", func(item interface{}) bool {
		_, ok := item.(func(number uint64, hash, header, body, receipts, td []byte))
		if ok { log.Warn("PlugEth's AppendAncient is deprecated. Please update to ModifyAncients.") }
		return ok
	})
	if len(appendAncientFnList) > 0 {
		var (
			hash []byte
			header []byte
			body []byte
			receipts []byte
			td []byte
		)
		if hashi, ok := update[freezerHashTable]; ok {
			switch v := hashi.(type) {
			case []byte:
				hash = v
			default:
				hash, _ = rlp.EncodeToBytes(v)
			}
		}
		if headeri, ok := update[freezerHeaderTable]; ok {
			switch v := headeri.(type) {
			case []byte:
				header = v
			default:
				header, _ = rlp.EncodeToBytes(v)
			}
		}
		if bodyi, ok := update[freezerBodiesTable]; ok {
			switch v := bodyi.(type) {
			case []byte:
				body = v
			default:
				body, _ = rlp.EncodeToBytes(v)
			}
		}
		if receiptsi, ok := update[freezerReceiptTable]; ok {
			switch v := receiptsi.(type) {
			case []byte:
				receipts = v
			default:
				receipts, _ = rlp.EncodeToBytes(v)
			}
		}
		if tdi, ok := update[freezerDifficultyTable]; ok {
			switch v := tdi.(type) {
			case []byte:
				td = v
			default:
				td, _ = rlp.EncodeToBytes(v)
			}
		}
		for _, fni := range appendAncientFnList {
			if fn, ok := fni.(func(number uint64, hash, header, body, receipts, td []byte)); ok {
				fn(num, hash, header, body, receipts, td)
			}
		}
	}
}
