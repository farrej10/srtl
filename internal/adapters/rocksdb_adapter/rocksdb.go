package rocksdbadapter

import (
	"errors"

	"github.com/farrej10/srtl/internal/ports"
	"github.com/linxGnu/grocksdb"
	"go.uber.org/zap"
)

type rocksDb struct {
	db     *grocksdb.DB
	logger *zap.SugaredLogger
}

// create rocksdb with some defaults
func NewRocksDB(location string, ttl int, logger *zap.SugaredLogger) (ports.IDatabaseAccessor, error) {
	bbto := grocksdb.NewDefaultBlockBasedTableOptions()
	bbto.SetBlockCache(grocksdb.NewLRUCache(3 << 30))

	opts := grocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)
	db, err := grocksdb.OpenDbWithTTL(opts, location, ttl)
	if err != nil {
		return rocksDb{}, err
	}
	return rocksDb{db: db, logger: logger}, nil
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
	returnVal := make([]byte, len(val.Data()))
	copy(returnVal, val.Data())
	r.logger.Debugw("Get", "returnVal", string(returnVal), "value", string(val.Data()))
	return returnVal, nil
}

func (r rocksDb) Set(key []byte, value []byte) error {
	r.logger.Debugw("Set", "key", string(key), "value", string(value))
	return r.db.Put(grocksdb.NewDefaultWriteOptions(), key, value)
}

func (r rocksDb) Close() error {
	r.db.Close()
	return nil
}
