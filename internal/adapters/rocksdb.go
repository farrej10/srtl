package adapters

import (
	"errors"

	"github.com/farrej10/srtl/internal/ports"
	"github.com/linxGnu/grocksdb"
)

type rocksDb struct {
	db *grocksdb.DB
}

// create rocksdb with some defaults
func NewRocksDB(location string, ttl int) (ports.IDatabaseAccessor, error) {
	bbto := grocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(grocksdb.NewLRUCache(3 << 30))

	opts := grocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)
	db, err := grocksdb.OpenDbWithTTL(opts, location, ttl)
	if err != nil {
		return rocksDb{}, err
	}
	return rocksDb{db: db}, nil
}

func (r rocksDb) Get(key []byte) ([]byte, error) {
	val, err := r.db.Get(grocksdb.NewDefaultReadOptions(), key)
	defer val.Free()
	if err != nil {
		return nil, err
	}
	if !val.Exists() {
		return nil, errors.New("key not found")
	}
	var returnVal []byte
	copy(returnVal, val.Data())
	return returnVal, nil
}

func (r rocksDb) Set(key []byte, value []byte) error {
	return r.db.Put(grocksdb.NewDefaultWriteOptions(), key, value)
}
