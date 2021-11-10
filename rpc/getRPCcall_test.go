package rpc

import (
	"testing"

	"github.com/ethereum/go-ethereum/plugins"
)

func TestGetRPCCalls(t *testing.T) {
	invoked := false
	done := plugins.HookTester("GetRPCCalls", func(id, method, params string) {
		invoked = true
		if id == "" {
			t.Errorf("Expected id to be non-nil")
		}
		if method == "" {
			t.Errorf("Expected method to be non-nil")
		}
		if params == "" {
			t.Errorf("Expected params to be non-nil")
		}
	})
	defer done()
	TestClientResponseType(t)
	if !invoked {
		t.Errorf("Expected plugin invocation")
	}
}
