package rawdb


import (
	"github.com/ethereum/go-ethereum/plugins"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"sync"
)

var (
	freezerUpdates map[uint64]map[string]interface{}
	lock sync.Mutex
)

func PluginTrackUpdate(num uint64, kind string, value interface{}) {
	lock.Lock()
	defer lock.Unlock()
	if freezerUpdates == nil { freezerUpdates = make(map[uint64]map[string]interface{}) }
	update, ok := freezerUpdates[num]
	if !ok {
		update = make(map[string]interface{})
		freezerUpdates[num] = update
	}
	update[kind] = value
}

func pluginCommitUpdate(num uint64) {
	if plugins.DefaultPluginLoader == nil {
		log.Warn("Attempting CommitUpdate, but default PluginLoader has not been initialized")
		return
	}
	PluginCommitUpdate(plugins.DefaultPluginLoader, num)
}

func PluginCommitUpdate(pl *plugins.PluginLoader, num uint64) {
	lock.Lock()
	defer lock.Unlock()
	if freezerUpdates == nil { freezerUpdates = make(map[uint64]map[string]interface{}) }
	min := ^uint64(0)
	for i := range freezerUpdates{
		if min > i { min = i }
	}
	for i := min ; i < num; i++ {
		update, ok := freezerUpdates[i]
		defer func(i uint64) { delete(freezerUpdates, i) }(i)
		if !ok {
			log.Warn("Attempting to commit untracked block", "num", i)
			continue
		}
		fnList := pl.Lookup("ModifyAncients", func(item interface{}) bool {
			_, ok := item.(func(uint64, map[string]interface{}))
			return ok
		})
		for _, fni := range fnList {
			if fn, ok := fni.(func(uint64, map[string]interface{})); ok {
				fn(i, update)
			}
		}
		appendAncientFnList := pl.Lookup("AppendAncient", func(item interface{}) bool {
			_, ok := item.(func(number uint64, hash, header, body, receipts, td []byte))
			if ok { log.Warn("PluGeth's AppendAncient is deprecated. Please update to ModifyAncients.") }
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
			if hashi, ok := update[ChainFreezerHashTable]; ok {
				switch v := hashi.(type) {
				case []byte:
					hash = v
				default:
					hash, _ = rlp.EncodeToBytes(v)
				}
			}
			if headeri, ok := update[ChainFreezerHeaderTable]; ok {
				switch v := headeri.(type) {
				case []byte:
					header = v
				default:
					header, _ = rlp.EncodeToBytes(v)
				}
			}
			if bodyi, ok := update[ChainFreezerBodiesTable]; ok {
				switch v := bodyi.(type) {
				case []byte:
					body = v
				default:
					body, _ = rlp.EncodeToBytes(v)
				}
			}
			if receiptsi, ok := update[ChainFreezerReceiptTable]; ok {
				switch v := receiptsi.(type) {
				case []byte:
					receipts = v
				default:
					receipts, _ = rlp.EncodeToBytes(v)
				}
			}
			if tdi, ok := update[ChainFreezerDifficultyTable]; ok {
				switch v := tdi.(type) {
				case []byte:
					td = v
				default:
					td, _ = rlp.EncodeToBytes(v)
				}
			}
			for _, fni := range appendAncientFnList {
				if fn, ok := fni.(func(number uint64, hash, header, body, receipts, td []byte)); ok {
					fn(i, hash, header, body, receipts, td)
				}
			}
		}
	}
}
