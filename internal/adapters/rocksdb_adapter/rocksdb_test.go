package rocksdbadapter

import (
	"fmt"
	"testing"

	"go.uber.org/zap"
)

var table = []struct {
	key   []byte
	value []byte
}{
	{key: []byte("abcde"), value: []byte("abcde")},
	{key: []byte("wqert"), value: []byte("wqert")},
	{key: []byte("asdfg"), value: []byte("asdfg")},
	{key: []byte("zxcvb"), value: []byte("zxcvb")},
}

func BenchmarkSet(b *testing.B) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	db, err := NewRocksDB("./db", 9999, sugar)
	defer db.Close()
	if err != nil {
		panic(err)
	}
	for _, v := range table {
		b.Run(fmt.Sprintf("Key: %s", string(v.key)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				db.Set(v.key, v.value)
			}
		})
	}
}

func BenchmarkGet(b *testing.B) {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	db, err := NewRocksDB("./db", 9999, sugar)
	defer db.Close()
	if err != nil {
		panic(err)
	}
	for _, v := range table {
		b.Run(fmt.Sprintf("Key: %s", string(v.key)), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				db.Get(v.key)
			}
		})
	}
}
