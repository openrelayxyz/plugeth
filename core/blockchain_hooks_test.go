package core

import (
	"testing"

	"github.com/ethereum/go-ethereum/plugins"
	"github.com/openrelayxyz/plugeth-utils/core"
)

func TestReorgLongHeadersHook(t *testing.T) {
	invoked := false
	done := plugins.HookTester("NewHead", func(b []byte, h core.Hash, logs [][]byte) {
		invoked = true
		if b == nil {
			t.Errorf("Expected block to be non-nil")
		}
		if h == (core.Hash{}) {
			t.Errorf("Expected hash to be non-empty")
		}
		if len(logs) > 0 {
			t.Errorf("Expected some logs")
		}
	})
	defer done()
	testReorgLong(t, true)
	if !invoked {
		t.Errorf("Expected plugin invocation")
	}
}
