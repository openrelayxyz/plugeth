package rawdb


import (
	"fmt"
	"os"
	"testing"
	"github.com/ethereum/go-ethereum/ethdb"
)




func TestAncientsInjections(t *testing.T) {

	test_dir_path := "./injection_test_dir"
	f, _ := NewFreezer(test_dir_path, "plugeth hook test", false, uint32(0), map[string]bool{"test": false})

	t.Run(fmt.Sprintf("test ModifyAncients"), func(t *testing.T) {
		called := false
		injectionCalled = &called
		_, _ = f.ModifyAncients(func (ethdb.AncientWriteOp) error {return nil})
		if *injectionCalled != true {
			t.Fatalf("pluginCommitUpdate injection in ModifyAncients not called")
		}
	})

	os.RemoveAll(test_dir_path)

	fb := newFreezerBatch(f)

	t.Run(fmt.Sprintf("test Append"), func(t *testing.T) {
		var item interface{}
		called := false
		injectionCalled = &called
		_ = fb.Append("kind", uint64(0), item)
		if *injectionCalled != true {
			t.Fatalf("PluginTrackUpdate injection in Append not called")
		}
	})

	t.Run(fmt.Sprintf("test AppendRaw"), func(t *testing.T) {
		called := false
		injectionCalled = &called
		_ = fb.AppendRaw("kind", uint64(100), []byte{})
		if *injectionCalled != true {
			t.Fatalf("PluginTrackUpdate injection in AppendRaw not called")
		}
	})
}