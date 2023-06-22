package rawdb


import (
	"fmt"
	"testing"
	"math/big"

	"github.com/ethereum/go-ethereum/ethdb"
)


func TestPlugethInjections(t *testing.T) {
	var valuesRaw [][]byte
	var valuesRLP []*big.Int
	for x := 0; x < 100; x++ {
		v := getChunk(256, x)
		valuesRaw = append(valuesRaw, v)
		iv := big.NewInt(int64(x))
		iv = iv.Exp(iv, iv, nil)
		valuesRLP = append(valuesRLP, iv)
	}
	tables := map[string]bool{"raw": true, "rlp": false}
	f, _ := newFreezerForTesting(t, tables)

	t.Run(fmt.Sprintf("test plugeth injections"), func(t *testing.T) {
		called := false
		modifyAncientsInjection = &called

		_, _ = f.ModifyAncients(func(op ethdb.AncientWriteOp) error {

			appendRawInjection = &called
			_ = op.AppendRaw("raw", uint64(0), valuesRaw[0])
			if *appendRawInjection != true {
				t.Fatalf("pluginTrackUpdate injection in AppendRaw not called")
			}

			appendInjection = &called
			_ = op.Append("rlp", uint64(0), valuesRaw[0])
			if *appendInjection != true {
				t.Fatalf("pluginTrackUpdate injection in Append not called")
			}

			return nil
	})
		if *modifyAncientsInjection != true {
			t.Fatalf("pluginCommitUpdate injection in ModifyAncients not called")
		}
	})
}