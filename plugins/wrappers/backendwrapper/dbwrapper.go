package backendwrapper

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/openrelayxyz/plugeth-utils/restricted"
)

type dbWrapper struct {
	db ethdb.Database
}

func (d *dbWrapper) Has(key []byte) (bool, error)             { return d.db.Has(key) }
func (d *dbWrapper) Get(key []byte) ([]byte, error)           { return d.db.Get(key) }
func (d *dbWrapper) Put(key []byte, value []byte) error       { return d.db.Put(key, value) }
func (d *dbWrapper) Delete(key []byte) error                  { return d.db.Delete(key) }
func (d *dbWrapper) Stat(property string) (string, error)     { return d.db.Stat(property) }
func (d *dbWrapper) Compact(start []byte, limit []byte) error { return d.db.Compact(start, limit) }
func (d *dbWrapper) HasAncient(kind string, number uint64) (bool, error) {
	return d.db.HasAncient(kind, number)
}
func (d *dbWrapper) Ancient(kind string, number uint64) ([]byte, error) {
	return d.db.Ancient(kind, number)
}
func (d *dbWrapper) Ancients() (uint64, error)               { return d.db.Ancients() }
func (d *dbWrapper) AncientSize(kind string) (uint64, error) { return d.db.AncientSize(kind) }
func (d *dbWrapper) AppendAncient(number uint64, hash, header, body, receipt, td []byte) error {
	return fmt.Errorf("AppendAncient is no longer supported in geth 1.10.9 and above. Use ModifyAncients instead.")
}
func (d *dbWrapper) ModifyAncients(fn func(ethdb.AncientWriteOperator) error) (int64, error) {
	return d.db.ModifyAncients(fn)
}
func (d *dbWrapper) TruncateAncients(n uint64) error {
	return fmt.Errorf("TruncateAncients is no longer supported in geth 1.10.17 and above.")
}
func (d *dbWrapper) Sync() error                     { return d.db.Sync() }
func (d *dbWrapper) Close() error                    { return d.db.Close() }
func (d *dbWrapper) NewIterator(prefix []byte, start []byte) restricted.Iterator {
	return &iterWrapper{d.db.NewIterator(prefix, start)}
}

type iterWrapper struct {
	iter ethdb.Iterator
}

func (it *iterWrapper) Next() bool    { return it.iter.Next() }
func (it *iterWrapper) Error() error  { return it.iter.Error() }
func (it *iterWrapper) Key() []byte   { return it.iter.Key() }
func (it *iterWrapper) Value() []byte { return it.iter.Value() }
func (it *iterWrapper) Release()      { it.iter.Release() }
