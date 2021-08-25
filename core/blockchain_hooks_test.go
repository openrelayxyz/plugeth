package core

import (
  "testing"
  "github.com/ethereum/go-ethereum/plugins"
  "github.com/ethereum/go-ethereum/core/types"
  "github.com/ethereum/go-ethereum/common"
)


func TestReorgLongHeadersHook(t *testing.T) {
  invoked := false
  done := plugins.HookTester("NewHead", func(b *types.Block, h common.Hash, logs []*types.Log) {
    // invoked = true
    if b == nil { t.Errorf("Expected block to be non-nil") }
    if h == (common.Hash{}) { t.Errorf("Expected hash to be non-empty") }
    if len(logs) > 0 { t.Errorf("Expected some logs") }
  })
  defer done()
  testReorgLong(t, true)
  if !invoked { t.Errorf("Expected plugin invocation")}
}
