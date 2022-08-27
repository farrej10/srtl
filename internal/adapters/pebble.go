package adapters

import (
	"errors"

	"github.com/cockroachdb/pebble"
	"github.com/farrej10/srtl/internal/ports"
	"go.uber.org/zap"
)

type pebbleDb struct {
	db     *pebble.DB
	logger zap.SugaredLogger
}

func NewPebbleDb(location string, logger zap.SugaredLogger) (ports.IDatabaseAccessor, error) {
	db, err := pebble.Open(location, &pebble.Options{})
	if err != nil {
		return pebbleDb{}, err
	}
	return pebbleDb{db: db, logger: logger}, nil
}

func (p pebbleDb) Get(key []byte) ([]byte, error) {
	value, closer, err := p.db.Get(key)
	if err != nil && err.Error() == "pebble: not found" {
		return nil, errors.New("key not found")
	}
	if err != nil {
		return nil, err
	}
	defer closer.Close()
	returnVal := make([]byte, len(value))
	copy(returnVal, value)
	p.logger.Debugw("Get", "returnVal", string(returnVal), "value", string(value))
	return returnVal, nil
}

func (p pebbleDb) Set(key []byte, value []byte) error {
	p.logger.Debugw("Set", "key", string(key), "value", string(value))
	return p.db.Set(key, value, pebble.Sync)
}
