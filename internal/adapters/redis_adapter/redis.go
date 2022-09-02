package redisadapter

import (
	"context"
	"errors"

	"github.com/farrej10/srtl/internal/ports"
	redis "github.com/go-redis/redis/v9"
)

type redisDb struct {
	db  *redis.Client
	ctx context.Context
}

func NewRedisDb(ip string, port string) (ports.IDatabaseAccessor, error) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     ip + ":" + port,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return redisDb{db: rdb, ctx: ctx}, nil
}

func (r redisDb) Get(key []byte) ([]byte, error) {
	val, err := r.db.Get(r.ctx, string(key)).Result()
	if err == redis.Nil {
		return nil, errors.New("key not found")
	}
	return []byte(val), err
}
func (r redisDb) Set(key []byte, value []byte) error {
	return r.db.Set(r.ctx, string(key), string(value), 0).Err()
}
func (r redisDb) Close() error {
	return r.db.Close()
}
